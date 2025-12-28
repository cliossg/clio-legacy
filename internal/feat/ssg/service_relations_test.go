package ssg

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestServiceAddTagToContent(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo)
		contentID uuid.UUID
		tagName   string
		wantErr   bool
		checkRepo func(*testing.T, *mockServiceRepo)
	}{
		{
			name: "uses existing tag",
			setup: func(m *mockServiceRepo) {
				existingTag := Tag{
					ID:   uuid.New(),
					Name: "existing-tag",
				}
				m.tags[existingTag.ID] = existingTag
				m.tagsByName["existing-tag"] = existingTag
			},
			contentID: uuid.New(),
			tagName:   "existing-tag",
			wantErr:   false,
			checkRepo: func(t *testing.T, m *mockServiceRepo) {
				if m.createTagCalled {
					t.Error("CreateTag should not be called when tag exists")
				}
				if !m.addTagToContentCalled {
					t.Error("AddTagToContent should be called")
				}
			},
		},
		{
			name: "creates new tag when not found",
			setup: func(m *mockServiceRepo) {
			},
			contentID: uuid.New(),
			tagName:   "new-tag",
			wantErr:   false,
			checkRepo: func(t *testing.T, m *mockServiceRepo) {
				if !m.createTagCalled {
					t.Error("CreateTag should be called for new tag")
				}
				if !m.addTagToContentCalled {
					t.Error("AddTagToContent should be called")
				}
			},
		},
		{
			name: "returns error when GetTagByName fails",
			setup: func(m *mockServiceRepo) {
				m.getTagErr = fmt.Errorf("db error")
			},
			contentID: uuid.New(),
			tagName:   "any-tag",
			wantErr:   true,
		},
		{
			name: "returns error when CreateTag fails",
			setup: func(m *mockServiceRepo) {
				m.createTagErr = fmt.Errorf("create error")
			},
			contentID: uuid.New(),
			tagName:   "new-tag",
			wantErr:   true,
		},
		{
			name: "returns error when AddTagToContent fails",
			setup: func(m *mockServiceRepo) {
				m.addTagToContentErr = fmt.Errorf("add error")
			},
			contentID: uuid.New(),
			tagName:   "new-tag",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.AddTagToContent(context.Background(), tt.contentID, tt.tagName)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTagToContent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkRepo != nil {
				tt.checkRepo(t, repo)
			}
		})
	}
}

func TestServiceRemoveTagFromContent(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo)
		contentID uuid.UUID
		tagID     uuid.UUID
		wantErr   bool
	}{
		{
			name:      "removes tag from content successfully",
			setup:     func(m *mockServiceRepo) {},
			contentID: uuid.New(),
			tagID:     uuid.New(),
			wantErr:   false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.removeTagFromContentErr = fmt.Errorf("db error")
			},
			contentID: uuid.New(),
			tagID:     uuid.New(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.RemoveTagFromContent(context.Background(), tt.contentID, tt.tagID)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveTagFromContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetTagsForContent(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo) uuid.UUID
		wantLen   int
		wantErr   bool
	}{
		{
			name: "gets tags for content successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				contentID := uuid.New()
				tag1 := Tag{ID: uuid.New(), Name: "tag1"}
				tag2 := Tag{ID: uuid.New(), Name: "tag2"}
				m.contentTags[contentID] = []Tag{tag1, tag2}
				return contentID
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getTagsForContentErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setup(repo)
			svc := newTestService(repo)

			tags, err := svc.GetTagsForContent(context.Background(), contentID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsForContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(tags) != tt.wantLen {
				t.Errorf("GetTagsForContent() len = %d, want %d", len(tags), tt.wantLen)
			}
		})
	}
}

func TestServiceGetContentForTag(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantLen int
		wantErr bool
	}{
		{
			name: "gets content for tag successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				tagID := uuid.New()
				content1 := Content{ID: uuid.New(), Heading: "Content 1"}
				content2 := Content{ID: uuid.New(), Heading: "Content 2"}
				m.tagContent[tagID] = []Content{content1, content2}
				return tagID
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getContentForTagErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tagID := tt.setup(repo)
			svc := newTestService(repo)

			content, err := svc.GetContentForTag(context.Background(), tagID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentForTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(content) != tt.wantLen {
				t.Errorf("GetContentForTag() len = %d, want %d", len(content), tt.wantLen)
			}
		})
	}
}

func TestServiceGetContentImages(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantLen int
		wantErr bool
	}{
		{
			name: "aggregates multiple images with metadata",
			setup: func(m *mockServiceRepo) uuid.UUID {
				contentID := uuid.New()
				img1 := Image{ID: uuid.New(), FileName: "img1.jpg"}
				img2 := Image{ID: uuid.New(), FileName: "img2.jpg"}
				m.images[img1.ID] = img1
				m.images[img2.ID] = img2
				m.contentImages[contentID] = []ContentImage{
					{ImageID: img1.ID, IsHeader: true},
					{ImageID: img2.ID, IsHeader: false},
				}
				return contentID
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "handles missing images gracefully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				contentID := uuid.New()
				img1 := Image{ID: uuid.New(), FileName: "img1.jpg"}
				m.images[img1.ID] = img1
				missingImageID := uuid.New()
				m.contentImages[contentID] = []ContentImage{
					{ImageID: img1.ID, IsHeader: true},
					{ImageID: missingImageID, IsHeader: false},
				}
				return contentID
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name: "returns error when GetContentImagesByContentID fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getContentImagesByContentIDErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setup(repo)
			svc := newTestService(repo)

			images, err := svc.GetContentImages(context.Background(), contentID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(images) != tt.wantLen {
				t.Errorf("GetContentImages() len = %d, want %d", len(images), tt.wantLen)
			}
		})
	}
}

func TestServiceGetContentHeaderImage(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo) uuid.UUID
		wantPath  string
		wantEmpty bool
		wantErr   bool
	}{
		{
			name: "returns header image path",
			setup: func(m *mockServiceRepo) uuid.UUID {
				contentID := uuid.New()
				img := Image{ID: uuid.New(), FilePath: "/images/header.jpg"}
				m.images[img.ID] = img
				m.contentImages[contentID] = []ContentImage{
					{ImageID: img.ID, IsHeader: true},
				}
				return contentID
			},
			wantPath: "/images/header.jpg",
			wantErr:  false,
		},
		{
			name: "returns empty when no header image",
			setup: func(m *mockServiceRepo) uuid.UUID {
				contentID := uuid.New()
				img := Image{ID: uuid.New(), FilePath: "/images/regular.jpg"}
				m.images[img.ID] = img
				m.contentImages[contentID] = []ContentImage{
					{ImageID: img.ID, IsHeader: false},
				}
				return contentID
			},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name: "returns error when GetContentImagesByContentID fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getContentImagesByContentIDErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setup(repo)
			svc := newTestService(repo)

			path, err := svc.GetContentHeaderImage(context.Background(), contentID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentHeaderImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantEmpty && path != "" {
					t.Errorf("GetContentHeaderImage() = %q, want empty string", path)
				} else if !tt.wantEmpty && path != tt.wantPath {
					t.Errorf("GetContentHeaderImage() = %q, want %q", path, tt.wantPath)
				}
			}
		})
	}
}

func TestServiceGetSectionHeaderImage(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo) uuid.UUID
		wantPath  string
		wantEmpty bool
		wantErr   bool
	}{
		{
			name: "returns header image path",
			setup: func(m *mockServiceRepo) uuid.UUID {
				sectionID := uuid.New()
				img := Image{ID: uuid.New(), FilePath: "/images/section-header.jpg"}
				m.images[img.ID] = img
				m.sectionImages[sectionID] = []SectionImage{
					{ImageID: img.ID, IsHeader: true},
				}
				return sectionID
			},
			wantPath: "/images/section-header.jpg",
			wantErr:  false,
		},
		{
			name: "returns empty when no header image",
			setup: func(m *mockServiceRepo) uuid.UUID {
				sectionID := uuid.New()
				m.sectionImages[sectionID] = []SectionImage{}
				return sectionID
			},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name: "returns error when GetSectionImagesBySectionID fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getSectionImagesBySectionIDErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			sectionID := tt.setup(repo)
			svc := newTestService(repo)

			path, err := svc.GetSectionHeaderImage(context.Background(), sectionID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSectionHeaderImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantEmpty && path != "" {
					t.Errorf("GetSectionHeaderImage() = %q, want empty string", path)
				} else if !tt.wantEmpty && path != tt.wantPath {
					t.Errorf("GetSectionHeaderImage() = %q, want %q", path, tt.wantPath)
				}
			}
		})
	}
}

func TestServiceGetSectionBlogHeaderImage(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo) uuid.UUID
		wantPath  string
		wantEmpty bool
		wantErr   bool
	}{
		{
			name: "returns blog header image path",
			setup: func(m *mockServiceRepo) uuid.UUID {
				sectionID := uuid.New()
				img := Image{ID: uuid.New(), FilePath: "/images/blog-header.jpg"}
				m.images[img.ID] = img
				m.sectionImages[sectionID] = []SectionImage{
					{ImageID: img.ID, IsHeader: true},
				}
				return sectionID
			},
			wantPath: "/images/blog-header.jpg",
			wantErr:  false,
		},
		{
			name: "returns empty when no header image",
			setup: func(m *mockServiceRepo) uuid.UUID {
				sectionID := uuid.New()
				m.sectionImages[sectionID] = []SectionImage{}
				return sectionID
			},
			wantEmpty: true,
			wantErr:   false,
		},
		{
			name: "returns error when GetSectionImagesBySectionID fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getSectionImagesBySectionIDErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			sectionID := tt.setup(repo)
			svc := newTestService(repo)

			path, err := svc.GetSectionBlogHeaderImage(context.Background(), sectionID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSectionBlogHeaderImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantEmpty && path != "" {
					t.Errorf("GetSectionBlogHeaderImage() = %q, want empty string", path)
				} else if !tt.wantEmpty && path != tt.wantPath {
					t.Errorf("GetSectionBlogHeaderImage() = %q, want %q", path, tt.wantPath)
				}
			}
		})
	}
}
