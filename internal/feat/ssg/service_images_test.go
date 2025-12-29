package ssg

import (
	"context"
	"fmt"
	"mime/multipart"
	"testing"
	"unsafe"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

type mockImageManager struct {
	hm.Core
	processResult *ImageProcessResult
	processErr    error
	deleteErr     error
	deleteCalled  bool
	deletedPaths  []string
}

func newMockImageManager() *mockImageManager {
	cfg := hm.NewConfig()
	core := hm.NewCore("mock-image-manager", hm.XParams{Cfg: cfg})
	return &mockImageManager{
		Core:         core,
		deletedPaths: []string{},
	}
}

func (m *mockImageManager) ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, content *Content, section *Section, imageType ImageType, altText, caption string) (*ImageProcessResult, error) {
	if m.processErr != nil {
		return nil, m.processErr
	}
	if m.processResult == nil {
		m.processResult = &ImageProcessResult{
			RelativePath: "images/test.jpg",
		}
	}
	return m.processResult, nil
}

func (m *mockImageManager) DeleteImage(ctx context.Context, path string) error {
	m.deleteCalled = true
	m.deletedPaths = append(m.deletedPaths, path)
	if m.deleteErr != nil {
		return m.deleteErr
	}
	return nil
}

func newTestServiceWithImageManager(repo Repo, im *mockImageManager) *BaseService {
	cfg := hm.NewConfig()
	log := hm.NewLogger("debug")
	params := hm.XParams{Cfg: cfg, Log: log}

	svc := &BaseService{
		Service: hm.NewService("test-service", params),
		repo:    repo,
		im:      (*ImageManager)(unsafe.Pointer(im)),
	}

	return svc
}

func TestServiceUploadContentImage(t *testing.T) {
	tests := []struct {
		name          string
		setupRepo     func(*mockServiceRepo)
		setupIM       func(*mockImageManager)
		contentID     uuid.UUID
		wantErr       bool
		checkRollback func(*testing.T, *mockImageManager)
	}{
		{
			name: "uploads content image successfully",
			setupRepo: func(m *mockServiceRepo) {
				contentID := uuid.New()
				m.contents[contentID] = Content{
					ID:      contentID,
					Heading: "Test Content",
				}
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/uploaded.jpg",
				}
			},
			wantErr: false,
		},
		{
			name: "uploads content image with section",
			setupRepo: func(m *mockServiceRepo) {
				sectionID := uuid.New()
				contentID := uuid.New()
				m.sections[sectionID] = Section{
					ID:   sectionID,
					Name: "Blog",
					Path: "blog",
				}
				m.contents[contentID] = Content{
					ID:        contentID,
					SectionID: sectionID,
					Heading:   "Test Content",
				}
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/blog/test.jpg",
				}
			},
			wantErr: false,
		},
		{
			name: "fails when content not found",
			setupRepo: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("content not found")
			},
			setupIM:   func(m *mockImageManager) {},
			contentID: uuid.New(),
			wantErr:   true,
		},
		{
			name: "fails when section not found",
			setupRepo: func(m *mockServiceRepo) {
				sectionID := uuid.New()
				contentID := uuid.New()
				m.contents[contentID] = Content{
					ID:        contentID,
					SectionID: sectionID,
				}
				m.getSectionErr = fmt.Errorf("section not found")
			},
			setupIM: func(m *mockImageManager) {},
			wantErr: true,
		},
		{
			name: "fails and does not rollback when image processing fails",
			setupRepo: func(m *mockServiceRepo) {
				contentID := uuid.New()
				m.contents[contentID] = Content{
					ID: contentID,
				}
			},
			setupIM: func(m *mockImageManager) {
				m.processErr = fmt.Errorf("upload failed")
			},
			wantErr: true,
			checkRollback: func(t *testing.T, im *mockImageManager) {
				if im.deleteCalled {
					t.Error("DeleteImage should not be called when ProcessUpload fails")
				}
			},
		},
		{
			name: "rolls back file when image record creation fails",
			setupRepo: func(m *mockServiceRepo) {
				contentID := uuid.New()
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.createImageErr = fmt.Errorf("db error")
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/rollback.jpg",
				}
			},
			wantErr: true,
			checkRollback: func(t *testing.T, im *mockImageManager) {
				if !im.deleteCalled {
					t.Error("DeleteImage should be called for rollback")
				}
				if len(im.deletedPaths) != 1 || im.deletedPaths[0] != "images/rollback.jpg" {
					t.Errorf("Expected rollback deletion of images/rollback.jpg, got %v", im.deletedPaths)
				}
			},
		},
		{
			name: "rolls back file and image when content image relationship creation fails",
			setupRepo: func(m *mockServiceRepo) {
				contentID := uuid.New()
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.createContentImageErr = fmt.Errorf("relationship error")
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/rollback2.jpg",
				}
			},
			wantErr: true,
			checkRollback: func(t *testing.T, im *mockImageManager) {
				if !im.deleteCalled {
					t.Error("DeleteImage should be called for rollback")
				}
				if len(im.deletedPaths) != 1 || im.deletedPaths[0] != "images/rollback2.jpg" {
					t.Errorf("Expected rollback deletion of images/rollback2.jpg, got %v", im.deletedPaths)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			im := newMockImageManager()
			tt.setupRepo(repo)
			tt.setupIM(im)
			svc := newTestServiceWithImageManager(repo, im)

			if tt.contentID == uuid.Nil {
				for id := range repo.contents {
					tt.contentID = id
					break
				}
			}

			ctx := context.Background()
			_, err := svc.UploadContentImage(ctx, tt.contentID, nil, nil, ImageTypeContent, "alt text", "caption")

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadContentImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkRollback != nil {
				tt.checkRollback(t, im)
			}
		})
	}
}

func TestServiceDeleteContentImage(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockServiceRepo, uuid.UUID)
		setupIM   func(*mockImageManager)
		imagePath string
		wantErr   bool
		checkIM   func(*testing.T, *mockImageManager)
	}{
		{
			name: "deletes content image successfully",
			setupRepo: func(m *mockServiceRepo, contentID uuid.UUID) {
				imageID := uuid.New()
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.images[imageID] = Image{
					ID:       imageID,
					FilePath: "images/delete-me.jpg",
				}
				m.contentImages[contentID] = []ContentImage{
					{
						ID:        uuid.New(),
						ContentID: contentID,
						ImageID:   imageID,
					},
				}
			},
			setupIM:   func(m *mockImageManager) {},
			imagePath: "images/delete-me.jpg",
			wantErr:   false,
			checkIM: func(t *testing.T, im *mockImageManager) {
				if !im.deleteCalled {
					t.Error("DeleteImage should be called")
				}
			},
		},
		{
			name: "deletes orphaned file when not in database",
			setupRepo: func(m *mockServiceRepo, contentID uuid.UUID) {
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.contentImages[contentID] = []ContentImage{}
			},
			setupIM:   func(m *mockImageManager) {},
			imagePath: "images/orphan.jpg",
			wantErr:   false,
			checkIM: func(t *testing.T, im *mockImageManager) {
				if !im.deleteCalled {
					t.Error("DeleteImage should be called for orphaned file")
				}
			},
		},
		{
			name: "fails when content not found",
			setupRepo: func(m *mockServiceRepo, contentID uuid.UUID) {
				m.getContentErr = fmt.Errorf("content not found")
			},
			setupIM:   func(m *mockImageManager) {},
			imagePath: "images/test.jpg",
			wantErr:   true,
		},
		{
			name: "fails when cannot get content images",
			setupRepo: func(m *mockServiceRepo, contentID uuid.UUID) {
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.getContentImagesByContentIDErr = fmt.Errorf("db error")
			},
			setupIM:   func(m *mockImageManager) {},
			imagePath: "images/test.jpg",
			wantErr:   true,
		},
		{
			name: "fails when cannot delete relationship",
			setupRepo: func(m *mockServiceRepo, contentID uuid.UUID) {
				imageID := uuid.New()
				m.contents[contentID] = Content{
					ID: contentID,
				}
				m.images[imageID] = Image{
					ID:       imageID,
					FilePath: "images/test.jpg",
				}
				m.contentImages[contentID] = []ContentImage{
					{
						ID:        uuid.New(),
						ContentID: contentID,
						ImageID:   imageID,
					},
				}
				m.deleteContentImageErr = fmt.Errorf("cannot delete relationship")
			},
			setupIM:   func(m *mockImageManager) {},
			imagePath: "images/test.jpg",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			im := newMockImageManager()
			contentID := uuid.New()
			tt.setupRepo(repo, contentID)
			tt.setupIM(im)
			svc := newTestServiceWithImageManager(repo, im)

			ctx := context.Background()
			err := svc.DeleteContentImage(ctx, contentID, tt.imagePath)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContentImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkIM != nil {
				tt.checkIM(t, im)
			}
		})
	}
}

func TestServiceUploadSectionImage(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockServiceRepo)
		setupIM   func(*mockImageManager)
		sectionID uuid.UUID
		wantErr   bool
	}{
		{
			name: "uploads section image successfully",
			setupRepo: func(m *mockServiceRepo) {
				sectionID := uuid.New()
				m.sections[sectionID] = Section{
					ID:   sectionID,
					Name: "Test Section",
					Path: "test",
				}
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/section/header.jpg",
				}
			},
			wantErr: false,
		},
		{
			name: "fails when section not found",
			setupRepo: func(m *mockServiceRepo) {
				m.getSectionErr = fmt.Errorf("section not found")
			},
			setupIM:   func(m *mockImageManager) {},
			sectionID: uuid.New(),
			wantErr:   true,
		},
		{
			name: "rolls back on image creation failure",
			setupRepo: func(m *mockServiceRepo) {
				sectionID := uuid.New()
				m.sections[sectionID] = Section{
					ID: sectionID,
				}
				m.createImageErr = fmt.Errorf("db error")
			},
			setupIM: func(m *mockImageManager) {
				m.processResult = &ImageProcessResult{
					RelativePath: "images/rollback-section.jpg",
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			im := newMockImageManager()
			tt.setupRepo(repo)
			tt.setupIM(im)
			svc := newTestServiceWithImageManager(repo, im)

			if tt.sectionID == uuid.Nil {
				for id := range repo.sections {
					tt.sectionID = id
					break
				}
			}

			ctx := context.Background()
			_, err := svc.UploadSectionImage(ctx, tt.sectionID, nil, nil, ImageTypeHeader, "alt", "caption")

			if (err != nil) != tt.wantErr {
				t.Errorf("UploadSectionImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteSectionImage(t *testing.T) {
	tests := []struct {
		name      string
		setupRepo func(*mockServiceRepo, uuid.UUID)
		setupIM   func(*mockImageManager)
		imageType ImageType
		wantErr   bool
	}{
		{
			name: "deletes section image successfully",
			setupRepo: func(m *mockServiceRepo, sectionID uuid.UUID) {
				imageID := uuid.New()
				m.sections[sectionID] = Section{
					ID: sectionID,
				}
				m.images[imageID] = Image{
					ID:       imageID,
					FilePath: "images/section.jpg",
				}
				m.sectionImages[sectionID] = []SectionImage{
					{
						ID:        uuid.New(),
						SectionID: sectionID,
						ImageID:   imageID,
						IsHeader:  true,
					},
				}
			},
			setupIM:   func(m *mockImageManager) {},
			imageType: ImageTypeHeader,
			wantErr:   false,
		},
		{
			name: "fails when section not found",
			setupRepo: func(m *mockServiceRepo, sectionID uuid.UUID) {
				m.getSectionErr = fmt.Errorf("section not found")
			},
			setupIM:   func(m *mockImageManager) {},
			imageType: ImageTypeHeader,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			im := newMockImageManager()
			sectionID := uuid.New()
			tt.setupRepo(repo, sectionID)
			tt.setupIM(im)
			svc := newTestServiceWithImageManager(repo, im)

			ctx := context.Background()
			err := svc.DeleteSectionImage(ctx, sectionID, tt.imageType)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSectionImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
