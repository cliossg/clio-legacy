package fake_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/fake"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

func TestSsgRepoCreateContent(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		content     *ssg.Content
		expectedErr error
	}{
		{
			name:      "creates content successfully",
			setupFake: func(f *fake.SsgRepo) {},
			content: &ssg.Content{
				ID:      uuid.New(),
				Heading: "Test Content",
			},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContentFn = func(ctx context.Context, content *ssg.Content) error {
					return errors.New("db error")
				}
			},
			content:     &ssg.Content{ID: uuid.New()},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateContent(context.Background(), tt.content)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetContent(t *testing.T) {
	contentID := uuid.New()

	tests := []struct {
		name            string
		setupFake       func(f *fake.SsgRepo)
		id              uuid.UUID
		expectedContent ssg.Content
		expectedErr     error
	}{
		{
			name: "gets existing content",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: contentID, Heading: "Test"})
			},
			id:              contentID,
			expectedContent: ssg.Content{ID: contentID, Heading: "Test"},
			expectedErr:     nil,
		},
		{
			name:            "returns error when not found",
			setupFake:       func(f *fake.SsgRepo) {},
			id:              uuid.New(),
			expectedContent: ssg.Content{},
			expectedErr:     errors.New("content not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetContentFn = func(ctx context.Context, id uuid.UUID) (ssg.Content, error) {
					return ssg.Content{}, errors.New("db error")
				}
			},
			id:              contentID,
			expectedContent: ssg.Content{},
			expectedErr:     errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			content, err := f.GetContent(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if content.ID != tt.expectedContent.ID || content.Heading != tt.expectedContent.Heading {
				t.Errorf("expected content %+v, got %+v", tt.expectedContent, content)
			}
		})
	}
}

func TestSsgRepoUpdateContent(t *testing.T) {
	contentID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		content     *ssg.Content
		expectedErr error
	}{
		{
			name: "updates content successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: contentID, Heading: "Old"})
			},
			content:     &ssg.Content{ID: contentID, Heading: "New"},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.UpdateContentFn = func(ctx context.Context, content *ssg.Content) error {
					return errors.New("db error")
				}
			},
			content:     &ssg.Content{ID: contentID, Heading: "New"},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateContent(context.Background(), tt.content)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoDeleteContent(t *testing.T) {
	contentID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes content successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: contentID})
			},
			id:          contentID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.DeleteContentFn = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("db error")
				}
			},
			id:          contentID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteContent(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetAllContentWithMeta(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		expectedLen int
		expectedErr error
	}{
		{
			name: "returns all content",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: uuid.New(), Heading: "Content 1"})
				f.CreateContent(context.Background(), &ssg.Content{ID: uuid.New(), Heading: "Content 2"})
			},
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no content",
			setupFake:   func(f *fake.SsgRepo) {},
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetAllContentWithMetaFn = func(ctx context.Context) ([]ssg.Content, error) {
					return nil, errors.New("db error")
				}
			},
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			contents, err := f.GetAllContentWithMeta(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(contents) != tt.expectedLen {
				t.Errorf("expected %d contents, got %d", tt.expectedLen, len(contents))
			}
		})
	}
}

func TestSsgRepoCreateTag(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		tag         ssg.Tag
		expectedErr error
	}{
		{
			name:        "creates tag successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			tag:         ssg.Tag{ID: uuid.New(), Name: "golang"},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTagFn = func(ctx context.Context, tag ssg.Tag) error {
					return errors.New("db error")
				}
			},
			tag:         ssg.Tag{ID: uuid.New(), Name: "golang"},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateTag(context.Background(), tt.tag)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetTagByName(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		tagName     string
		expectedTag ssg.Tag
		expectedErr error
	}{
		{
			name: "gets tag by name",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: uuid.New(), Name: "golang"})
			},
			tagName:     "golang",
			expectedTag: ssg.Tag{Name: "golang"},
			expectedErr: nil,
		},
		{
			name:        "returns error when not found",
			setupFake:   func(f *fake.SsgRepo) {},
			tagName:     "nonexistent",
			expectedTag: ssg.Tag{},
			expectedErr: errors.New("tag not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetTagByNameFn = func(ctx context.Context, name string) (ssg.Tag, error) {
					return ssg.Tag{}, errors.New("db error")
				}
			},
			tagName:     "golang",
			expectedTag: ssg.Tag{},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			tag, err := f.GetTagByName(context.Background(), tt.tagName)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tag.Name != tt.expectedTag.Name {
				t.Errorf("expected tag name %v, got %v", tt.expectedTag.Name, tag.Name)
			}
		})
	}
}

func TestSsgRepoDeleteTag(t *testing.T) {
	tagID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes tag and removes from name index",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: tagID, Name: "golang"})
			},
			id:          tagID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.DeleteTagFn = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("db error")
				}
			},
			id:          tagID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteTag(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoCreateParam(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		param       *ssg.Param
		expectedErr error
	}{
		{
			name:      "creates param successfully",
			setupFake: func(f *fake.SsgRepo) {},
			param: &ssg.Param{
				ID:     uuid.New(),
				Name:   "site.title",
				RefKey: "ssg.site.title",
				Value:  "My Site",
			},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParamFn = func(ctx context.Context, param *ssg.Param) error {
					return errors.New("db error")
				}
			},
			param:       &ssg.Param{ID: uuid.New()},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateParam(context.Background(), tt.param)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetParamByName(t *testing.T) {
	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		paramName     string
		expectedParam ssg.Param
		expectedErr   error
	}{
		{
			name: "gets param by name",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{
					ID:     uuid.New(),
					Name:   "site.title",
					RefKey: "ssg.site.title",
					Value:  "My Site",
				})
			},
			paramName:     "site.title",
			expectedParam: ssg.Param{Name: "site.title", Value: "My Site"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when not found",
			setupFake:     func(f *fake.SsgRepo) {},
			paramName:     "nonexistent",
			expectedParam: ssg.Param{},
			expectedErr:   errors.New("param not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetParamByNameFn = func(ctx context.Context, name string) (ssg.Param, error) {
					return ssg.Param{}, errors.New("db error")
				}
			},
			paramName:     "site.title",
			expectedParam: ssg.Param{},
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			param, err := f.GetParamByName(context.Background(), tt.paramName)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if param.Name != tt.expectedParam.Name || param.Value != tt.expectedParam.Value {
				t.Errorf("expected param %+v, got %+v", tt.expectedParam, param)
			}
		})
	}
}

func TestSsgRepoGetParamByRefKey(t *testing.T) {
	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		refKey        string
		expectedParam ssg.Param
		expectedErr   error
	}{
		{
			name: "gets param by refkey",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{
					ID:     uuid.New(),
					Name:   "site.title",
					RefKey: "ssg.site.title",
					Value:  "My Site",
				})
			},
			refKey:        "ssg.site.title",
			expectedParam: ssg.Param{RefKey: "ssg.site.title", Value: "My Site"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when not found",
			setupFake:     func(f *fake.SsgRepo) {},
			refKey:        "nonexistent",
			expectedParam: ssg.Param{},
			expectedErr:   errors.New("param not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetParamByRefKeyFn = func(ctx context.Context, refKey string) (ssg.Param, error) {
					return ssg.Param{}, errors.New("db error")
				}
			},
			refKey:        "ssg.site.title",
			expectedParam: ssg.Param{},
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			param, err := f.GetParamByRefKey(context.Background(), tt.refKey)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if param.RefKey != tt.expectedParam.RefKey || param.Value != tt.expectedParam.Value {
				t.Errorf("expected param %+v, got %+v", tt.expectedParam, param)
			}
		})
	}
}

func TestSsgRepoDeleteParam(t *testing.T) {
	paramID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes param and removes from all indexes",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{
					ID:     paramID,
					Name:   "site.title",
					RefKey: "ssg.site.title",
				})
			},
			id:          paramID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.DeleteParamFn = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("db error")
				}
			},
			id:          paramID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteParam(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoCreateImage(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		image       *ssg.Image
		expectedErr error
	}{
		{
			name:      "creates image successfully",
			setupFake: func(f *fake.SsgRepo) {},
			image: &ssg.Image{
				ID:       uuid.New(),
				ShortID:  "abc123",
				FilePath: "images/test.jpg",
				FileName: "test.jpg",
			},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImageFn = func(ctx context.Context, image *ssg.Image) error {
					return errors.New("db error")
				}
			},
			image:       &ssg.Image{ID: uuid.New()},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateImage(context.Background(), tt.image)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetImageByShortID(t *testing.T) {
	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		shortID       string
		expectedImage ssg.Image
		expectedErr   error
	}{
		{
			name: "gets image by short id",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{
					ID:       uuid.New(),
					ShortID:  "abc123",
					FileName: "test.jpg",
				})
			},
			shortID:       "abc123",
			expectedImage: ssg.Image{ShortID: "abc123", FileName: "test.jpg"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when not found",
			setupFake:     func(f *fake.SsgRepo) {},
			shortID:       "nonexistent",
			expectedImage: ssg.Image{},
			expectedErr:   errors.New("image not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetImageByShortIDFn = func(ctx context.Context, shortID string) (ssg.Image, error) {
					return ssg.Image{}, errors.New("db error")
				}
			},
			shortID:       "abc123",
			expectedImage: ssg.Image{},
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			image, err := f.GetImageByShortID(context.Background(), tt.shortID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if image.ShortID != tt.expectedImage.ShortID || image.FileName != tt.expectedImage.FileName {
				t.Errorf("expected image %+v, got %+v", tt.expectedImage, image)
			}
		})
	}
}

func TestSsgRepoGetImageByContentHash(t *testing.T) {
	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		contentHash   string
		expectedImage ssg.Image
		expectedErr   error
	}{
		{
			name: "gets image by content hash",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{
					ID:       uuid.New(),
					FilePath: "images/test.jpg",
					FileName: "test.jpg",
				})
			},
			contentHash:   "images/test.jpg",
			expectedImage: ssg.Image{FilePath: "images/test.jpg", FileName: "test.jpg"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when not found",
			setupFake:     func(f *fake.SsgRepo) {},
			contentHash:   "nonexistent",
			expectedImage: ssg.Image{},
			expectedErr:   errors.New("image not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetImageByContentHashFn = func(ctx context.Context, contentHash string) (ssg.Image, error) {
					return ssg.Image{}, errors.New("db error")
				}
			},
			contentHash:   "images/test.jpg",
			expectedImage: ssg.Image{},
			expectedErr:   errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			image, err := f.GetImageByContentHash(context.Background(), tt.contentHash)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if image.FilePath != tt.expectedImage.FilePath || image.FileName != tt.expectedImage.FileName {
				t.Errorf("expected image %+v, got %+v", tt.expectedImage, image)
			}
		})
	}
}

func TestSsgRepoDeleteImage(t *testing.T) {
	imageID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes image and removes from all indexes",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{
					ID:       imageID,
					ShortID:  "abc123",
					FilePath: "images/test.jpg",
				})
			},
			id:          imageID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.DeleteImageFn = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("db error")
				}
			},
			id:          imageID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteImage(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoAddTagToContent(t *testing.T) {
	contentID := uuid.New()
	tagID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		contentID   uuid.UUID
		tagID       uuid.UUID
		expectedErr error
	}{
		{
			name:        "adds tag to content successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			contentID:   contentID,
			tagID:       tagID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.AddTagToContentFn = func(ctx context.Context, contentID, tagID uuid.UUID) error {
					return errors.New("db error")
				}
			},
			contentID:   contentID,
			tagID:       tagID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.AddTagToContent(context.Background(), tt.contentID, tt.tagID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoRemoveTagFromContent(t *testing.T) {
	contentID := uuid.New()
	tagID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		contentID   uuid.UUID
		tagID       uuid.UUID
		expectedErr error
	}{
		{
			name: "removes tag from content successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.AddTagToContent(context.Background(), contentID, tagID)
			},
			contentID:   contentID,
			tagID:       tagID,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.RemoveTagFromContentFn = func(ctx context.Context, contentID, tagID uuid.UUID) error {
					return errors.New("db error")
				}
			},
			contentID:   contentID,
			tagID:       tagID,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.RemoveTagFromContent(context.Background(), tt.contentID, tt.tagID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestSsgRepoGetTagsForContent(t *testing.T) {
	contentID := uuid.New()
	tag1ID := uuid.New()
	tag2ID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		contentID   uuid.UUID
		expectedLen int
		expectedErr error
	}{
		{
			name: "returns tags for content",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: tag1ID, Name: "golang"})
				f.CreateTag(context.Background(), ssg.Tag{ID: tag2ID, Name: "testing"})
				f.AddTagToContent(context.Background(), contentID, tag1ID)
				f.AddTagToContent(context.Background(), contentID, tag2ID)
			},
			contentID:   contentID,
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no tags",
			setupFake:   func(f *fake.SsgRepo) {},
			contentID:   contentID,
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetTagsForContentFn = func(ctx context.Context, contentID uuid.UUID) ([]ssg.Tag, error) {
					return nil, errors.New("db error")
				}
			},
			contentID:   contentID,
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			tags, err := f.GetTagsForContent(context.Background(), tt.contentID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(tags) != tt.expectedLen {
				t.Errorf("expected %d tags, got %d", tt.expectedLen, len(tags))
			}
		})
	}
}

func TestSsgRepoGetContentForTag(t *testing.T) {
	tagID := uuid.New()
	content1ID := uuid.New()
	content2ID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		tagID       uuid.UUID
		expectedLen int
		expectedErr error
	}{
		{
			name: "returns content for tag",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: content1ID, Heading: "Content 1"})
				f.CreateContent(context.Background(), &ssg.Content{ID: content2ID, Heading: "Content 2"})
				f.AddTagToContent(context.Background(), content1ID, tagID)
				f.AddTagToContent(context.Background(), content2ID, tagID)
			},
			tagID:       tagID,
			expectedLen: 2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no content",
			setupFake:   func(f *fake.SsgRepo) {},
			tagID:       tagID,
			expectedLen: 0,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetContentForTagFn = func(ctx context.Context, tagID uuid.UUID) ([]ssg.Content, error) {
					return nil, errors.New("db error")
				}
			},
			tagID:       tagID,
			expectedLen: 0,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			contents, err := f.GetContentForTag(context.Background(), tt.tagID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(contents) != tt.expectedLen {
				t.Errorf("expected %d contents, got %d", tt.expectedLen, len(contents))
			}
		})
	}
}

func TestSsgRepoGetUserByUsername(t *testing.T) {
	tests := []struct {
		name         string
		setupFake    func(f *fake.SsgRepo)
		username     string
		expectedUser auth.User
		expectedErr  error
	}{
		{
			name:         "returns error when user not found",
			setupFake:    func(f *fake.SsgRepo) {},
			username:     "testuser",
			expectedUser: auth.User{},
			expectedErr:  errors.New("user not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetUserByUsernameFn = func(ctx context.Context, username string) (auth.User, error) {
					return auth.User{}, errors.New("db error")
				}
			},
			username:     "testuser",
			expectedUser: auth.User{},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			user, err := f.GetUserByUsername(context.Background(), tt.username)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if user.Username != tt.expectedUser.Username {
				t.Errorf("expected user %+v, got %+v", tt.expectedUser, user)
			}
		})
	}
}

func TestSsgRepoGetSiteBySlug(t *testing.T) {
	tests := []struct {
		name         string
		setupFake    func(f *fake.SsgRepo)
		slug         string
		expectedSite ssg.Site
		expectedErr  error
	}{
		{
			name:         "returns error when site not found",
			setupFake:    func(f *fake.SsgRepo) {},
			slug:         "mysite",
			expectedSite: ssg.Site{},
			expectedErr:  errors.New("site not found"),
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetSiteBySlugFn = func(ctx context.Context, slug string) (ssg.Site, error) {
					return ssg.Site{}, errors.New("db error")
				}
			},
			slug:         "mysite",
			expectedSite: ssg.Site{},
			expectedErr:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			site, err := f.GetSiteBySlug(context.Background(), tt.slug)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if site.Slug() != tt.expectedSite.Slug() {
				t.Errorf("expected site %+v, got %+v", tt.expectedSite, site)
			}
		})
	}
}
