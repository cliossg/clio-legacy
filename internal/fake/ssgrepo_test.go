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
			name:        "returns nil when tag not found",
			setupFake:   func(f *fake.SsgRepo) {},
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

func TestSsgRepoQuery(t *testing.T) {
	f := fake.NewSsgRepo()
	qm := f.Query()
	if qm != nil {
		t.Errorf("expected Query() to return nil, got %v", qm)
	}
}

func TestSsgRepoBeginTx(t *testing.T) {
	f := fake.NewSsgRepo()
	ctx := context.Background()
	returnedCtx, tx, err := f.BeginTx(ctx)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if tx != nil {
		t.Errorf("expected tx to be nil, got %v", tx)
	}
	if returnedCtx != ctx {
		t.Errorf("expected context to be unchanged")
	}
}

func TestSsgRepoGetContentWithPaginationAndSearch(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		offset      int
		limit       int
		search      string
		wantCount   int
		wantTotal   int
		expectedErr error
	}{
		{
			name: "returns paginated content",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContent(context.Background(), &ssg.Content{ID: uuid.New(), Heading: "Content 1"})
				f.CreateContent(context.Background(), &ssg.Content{ID: uuid.New(), Heading: "Content 2"})
			},
			offset:      0,
			limit:       10,
			search:      "",
			wantCount:   2,
			wantTotal:   2,
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.GetContentWithPaginationAndSearchFn = func(ctx context.Context, offset, limit int, search string) ([]ssg.Content, int, error) {
					return nil, 0, errors.New("search failed")
				}
			},
			offset:      0,
			limit:       10,
			search:      "",
			wantCount:   0,
			wantTotal:   0,
			expectedErr: errors.New("search failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			contents, total, err := f.GetContentWithPaginationAndSearch(context.Background(), tt.offset, tt.limit, tt.search)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(contents) != tt.wantCount {
				t.Errorf("expected %d contents, got %d", tt.wantCount, len(contents))
			}

			if total != tt.wantTotal {
				t.Errorf("expected total %d, got %d", tt.wantTotal, total)
			}
		})
	}
}

func TestSsgRepoCreateSection(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		section     ssg.Section
		expectedErr error
	}{
		{
			name:        "creates section successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			section:     ssg.Section{ID: uuid.New(), Name: "Blog"},
			expectedErr: nil,
		},
		{
			name: "returns error from custom function",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSectionFn = func(ctx context.Context, section ssg.Section) error {
					return errors.New("create failed")
				}
			},
			section:     ssg.Section{ID: uuid.New(), Name: "Blog"},
			expectedErr: errors.New("create failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateSection(context.Background(), tt.section)

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

func TestSsgRepoGetSection(t *testing.T) {
	sectionID := uuid.New()

	tests := []struct {
		name            string
		setupFake       func(f *fake.SsgRepo)
		id              uuid.UUID
		expectedSection ssg.Section
		expectedErr     error
	}{
		{
			name: "gets existing section",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSection(context.Background(), ssg.Section{ID: sectionID, Name: "Blog"})
			},
			id:              sectionID,
			expectedSection: ssg.Section{ID: sectionID, Name: "Blog"},
			expectedErr:     nil,
		},
		{
			name:            "returns error when section not found",
			setupFake:       func(f *fake.SsgRepo) {},
			id:              uuid.New(),
			expectedSection: ssg.Section{},
			expectedErr:     errors.New("section not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			section, err := f.GetSection(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if section.ID != tt.expectedSection.ID || section.Name != tt.expectedSection.Name {
				t.Errorf("expected section %+v, got %+v", tt.expectedSection, section)
			}
		})
	}
}

func TestSsgRepoGetSections(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns all sections",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSection(context.Background(), ssg.Section{ID: uuid.New(), Name: "Blog"})
				f.CreateSection(context.Background(), ssg.Section{ID: uuid.New(), Name: "News"})
			},
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no sections",
			setupFake:   func(f *fake.SsgRepo) {},
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			sections, err := f.GetSections(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(sections) != tt.wantCount {
				t.Errorf("expected %d sections, got %d", tt.wantCount, len(sections))
			}
		})
	}
}

func TestSsgRepoUpdateSection(t *testing.T) {
	sectionID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		section     ssg.Section
		expectedErr error
	}{
		{
			name: "updates section successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSection(context.Background(), ssg.Section{ID: sectionID, Name: "Old Name"})
			},
			section:     ssg.Section{ID: sectionID, Name: "New Name"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateSection(context.Background(), tt.section)

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

func TestSsgRepoDeleteSection(t *testing.T) {
	sectionID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes section successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSection(context.Background(), ssg.Section{ID: sectionID, Name: "Blog"})
			},
			id:          sectionID,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteSection(context.Background(), tt.id)

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

func TestSsgRepoCreateLayout(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		layout      ssg.Layout
		expectedErr error
	}{
		{
			name:        "creates layout successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			layout:      ssg.Layout{ID: uuid.New(), Name: "Default"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateLayout(context.Background(), tt.layout)

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

func TestSsgRepoGetLayout(t *testing.T) {
	layoutID := uuid.New()

	tests := []struct {
		name           string
		setupFake      func(f *fake.SsgRepo)
		id             uuid.UUID
		expectedLayout ssg.Layout
		expectedErr    error
	}{
		{
			name: "gets existing layout",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateLayout(context.Background(), ssg.Layout{ID: layoutID, Name: "Default"})
			},
			id:             layoutID,
			expectedLayout: ssg.Layout{ID: layoutID, Name: "Default"},
			expectedErr:    nil,
		},
		{
			name:           "returns error when layout not found",
			setupFake:      func(f *fake.SsgRepo) {},
			id:             uuid.New(),
			expectedLayout: ssg.Layout{},
			expectedErr:    errors.New("layout not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			layout, err := f.GetLayout(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if layout.ID != tt.expectedLayout.ID || layout.Name != tt.expectedLayout.Name {
				t.Errorf("expected layout %+v, got %+v", tt.expectedLayout, layout)
			}
		})
	}
}

func TestSsgRepoGetAllLayouts(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns all layouts",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateLayout(context.Background(), ssg.Layout{ID: uuid.New(), Name: "Default"})
				f.CreateLayout(context.Background(), ssg.Layout{ID: uuid.New(), Name: "Alternative"})
			},
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no layouts",
			setupFake:   func(f *fake.SsgRepo) {},
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			layouts, err := f.GetAllLayouts(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(layouts) != tt.wantCount {
				t.Errorf("expected %d layouts, got %d", tt.wantCount, len(layouts))
			}
		})
	}
}

func TestSsgRepoUpdateLayout(t *testing.T) {
	layoutID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		layout      ssg.Layout
		expectedErr error
	}{
		{
			name: "updates layout successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateLayout(context.Background(), ssg.Layout{ID: layoutID, Name: "Old"})
			},
			layout:      ssg.Layout{ID: layoutID, Name: "New"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateLayout(context.Background(), tt.layout)

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

func TestSsgRepoDeleteLayout(t *testing.T) {
	layoutID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes layout successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateLayout(context.Background(), ssg.Layout{ID: layoutID, Name: "Default"})
			},
			id:          layoutID,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteLayout(context.Background(), tt.id)

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

func TestSsgRepoGetTag(t *testing.T) {
	tagID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedTag ssg.Tag
		expectedErr error
	}{
		{
			name: "gets existing tag",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: tagID, Name: "tech"})
			},
			id:          tagID,
			expectedTag: ssg.Tag{ID: tagID, Name: "tech"},
			expectedErr: nil,
		},
		{
			name:        "returns error when tag not found",
			setupFake:   func(f *fake.SsgRepo) {},
			id:          uuid.New(),
			expectedTag: ssg.Tag{},
			expectedErr: errors.New("tag not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			tag, err := f.GetTag(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tag.ID != tt.expectedTag.ID || tag.Name != tt.expectedTag.Name {
				t.Errorf("expected tag %+v, got %+v", tt.expectedTag, tag)
			}
		})
	}
}

func TestSsgRepoGetAllTags(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns all tags",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: uuid.New(), Name: "tech"})
				f.CreateTag(context.Background(), ssg.Tag{ID: uuid.New(), Name: "news"})
			},
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no tags",
			setupFake:   func(f *fake.SsgRepo) {},
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			tags, err := f.GetAllTags(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(tags) != tt.wantCount {
				t.Errorf("expected %d tags, got %d", tt.wantCount, len(tags))
			}
		})
	}
}

func TestSsgRepoUpdateTag(t *testing.T) {
	tagID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		tag         ssg.Tag
		expectedErr error
	}{
		{
			name: "updates tag successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateTag(context.Background(), ssg.Tag{ID: tagID, Name: "old"})
			},
			tag:         ssg.Tag{ID: tagID, Name: "new"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateTag(context.Background(), tt.tag)

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

func TestSsgRepoGetParam(t *testing.T) {
	paramID := uuid.New()

	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		id            uuid.UUID
		expectedParam ssg.Param
		expectedErr   error
	}{
		{
			name: "gets existing param",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{ID: paramID, Name: "title"})
			},
			id:            paramID,
			expectedParam: ssg.Param{ID: paramID, Name: "title"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when param not found",
			setupFake:     func(f *fake.SsgRepo) {},
			id:            uuid.New(),
			expectedParam: ssg.Param{},
			expectedErr:   errors.New("param not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			param, err := f.GetParam(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if param.ID != tt.expectedParam.ID || param.Name != tt.expectedParam.Name {
				t.Errorf("expected param %+v, got %+v", tt.expectedParam, param)
			}
		})
	}
}

func TestSsgRepoListParams(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns all params",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{ID: uuid.New(), Name: "title"})
				f.CreateParam(context.Background(), &ssg.Param{ID: uuid.New(), Name: "desc"})
			},
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no params",
			setupFake:   func(f *fake.SsgRepo) {},
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			params, err := f.ListParams(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(params) != tt.wantCount {
				t.Errorf("expected %d params, got %d", tt.wantCount, len(params))
			}
		})
	}
}

func TestSsgRepoUpdateParam(t *testing.T) {
	paramID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		param       *ssg.Param
		expectedErr error
	}{
		{
			name: "updates param successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateParam(context.Background(), &ssg.Param{ID: paramID, Name: "old"})
			},
			param:       &ssg.Param{ID: paramID, Name: "new"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateParam(context.Background(), tt.param)

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

func TestSsgRepoGetImage(t *testing.T) {
	imageID := uuid.New()

	tests := []struct {
		name          string
		setupFake     func(f *fake.SsgRepo)
		id            uuid.UUID
		expectedImage ssg.Image
		expectedErr   error
	}{
		{
			name: "gets existing image",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{ID: imageID, FilePath: "/path/to/image.jpg"})
			},
			id:            imageID,
			expectedImage: ssg.Image{ID: imageID, FilePath: "/path/to/image.jpg"},
			expectedErr:   nil,
		},
		{
			name:          "returns error when image not found",
			setupFake:     func(f *fake.SsgRepo) {},
			id:            uuid.New(),
			expectedImage: ssg.Image{},
			expectedErr:   errors.New("image not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			image, err := f.GetImage(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if image.ID != tt.expectedImage.ID || image.FilePath != tt.expectedImage.FilePath {
				t.Errorf("expected image %+v, got %+v", tt.expectedImage, image)
			}
		})
	}
}

func TestSsgRepoUpdateImage(t *testing.T) {
	imageID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		image       *ssg.Image
		expectedErr error
	}{
		{
			name: "updates image successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{ID: imageID, FilePath: "/old.jpg"})
			},
			image:       &ssg.Image{ID: imageID, FilePath: "/new.jpg"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateImage(context.Background(), tt.image)

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

func TestSsgRepoListImages(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns all images",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImage(context.Background(), &ssg.Image{ID: uuid.New(), FilePath: "/img1.jpg"})
				f.CreateImage(context.Background(), &ssg.Image{ID: uuid.New(), FilePath: "/img2.jpg"})
			},
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no images",
			setupFake:   func(f *fake.SsgRepo) {},
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			images, err := f.ListImages(context.Background())

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(images) != tt.wantCount {
				t.Errorf("expected %d images, got %d", tt.wantCount, len(images))
			}
		})
	}
}

func TestSsgRepoCreateImageVariant(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		variant     *ssg.ImageVariant
		expectedErr error
	}{
		{
			name:        "creates image variant successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			variant:     &ssg.ImageVariant{ID: uuid.New(), Kind: "thumbnail"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateImageVariant(context.Background(), tt.variant)

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

func TestSsgRepoGetImageVariant(t *testing.T) {
	variantID := uuid.New()

	tests := []struct {
		name            string
		setupFake       func(f *fake.SsgRepo)
		id              uuid.UUID
		expectedVariant ssg.ImageVariant
		expectedErr     error
	}{
		{
			name: "gets existing image variant",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImageVariant(context.Background(), &ssg.ImageVariant{ID: variantID, Kind: "thumbnail"})
			},
			id:              variantID,
			expectedVariant: ssg.ImageVariant{ID: variantID, Kind: "thumbnail"},
			expectedErr:     nil,
		},
		{
			name:            "returns error when variant not found",
			setupFake:       func(f *fake.SsgRepo) {},
			id:              uuid.New(),
			expectedVariant: ssg.ImageVariant{},
			expectedErr:     errors.New("image variant not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			variant, err := f.GetImageVariant(context.Background(), tt.id)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if variant.ID != tt.expectedVariant.ID || variant.Kind != tt.expectedVariant.Kind {
				t.Errorf("expected variant %+v, got %+v", tt.expectedVariant, variant)
			}
		})
	}
}

func TestSsgRepoUpdateImageVariant(t *testing.T) {
	variantID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		variant     *ssg.ImageVariant
		expectedErr error
	}{
		{
			name: "updates image variant successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImageVariant(context.Background(), &ssg.ImageVariant{ID: variantID, Kind: "old"})
			},
			variant:     &ssg.ImageVariant{ID: variantID, Kind: "new"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.UpdateImageVariant(context.Background(), tt.variant)

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

func TestSsgRepoDeleteImageVariant(t *testing.T) {
	variantID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes image variant successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImageVariant(context.Background(), &ssg.ImageVariant{ID: variantID, Kind: "thumbnail"})
			},
			id:          variantID,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteImageVariant(context.Background(), tt.id)

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

func TestSsgRepoListImageVariantsByImageID(t *testing.T) {
	imageID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		imageID     uuid.UUID
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns variants for image",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateImageVariant(context.Background(), &ssg.ImageVariant{ID: uuid.New(), ImageID: imageID, Kind: "thumb"})
				f.CreateImageVariant(context.Background(), &ssg.ImageVariant{ID: uuid.New(), ImageID: imageID, Kind: "medium"})
			},
			imageID:     imageID,
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no variants",
			setupFake:   func(f *fake.SsgRepo) {},
			imageID:     uuid.New(),
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			variants, err := f.ListImageVariantsByImageID(context.Background(), tt.imageID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(variants) != tt.wantCount {
				t.Errorf("expected %d variants, got %d", tt.wantCount, len(variants))
			}
		})
	}
}

func TestSsgRepoCreateContentImage(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		ci          *ssg.ContentImage
		expectedErr error
	}{
		{
			name:        "creates content image successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			ci:          &ssg.ContentImage{ID: uuid.New(), ContentID: uuid.New(), ImageID: uuid.New()},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateContentImage(context.Background(), tt.ci)

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

func TestSsgRepoDeleteContentImage(t *testing.T) {
	ciID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes content image successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContentImage(context.Background(), &ssg.ContentImage{ID: ciID})
			},
			id:          ciID,
			expectedErr: nil,
		},
		{
			name:        "returns nil when content image not found",
			setupFake:   func(f *fake.SsgRepo) {},
			id:          uuid.New(),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteContentImage(context.Background(), tt.id)

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

func TestSsgRepoGetContentImagesByContentID(t *testing.T) {
	contentID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		contentID   uuid.UUID
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns content images",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateContentImage(context.Background(), &ssg.ContentImage{ID: uuid.New(), ContentID: contentID})
				f.CreateContentImage(context.Background(), &ssg.ContentImage{ID: uuid.New(), ContentID: contentID})
			},
			contentID:   contentID,
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no content images",
			setupFake:   func(f *fake.SsgRepo) {},
			contentID:   uuid.New(),
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			cis, err := f.GetContentImagesByContentID(context.Background(), tt.contentID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(cis) != tt.wantCount {
				t.Errorf("expected %d content images, got %d", tt.wantCount, len(cis))
			}
		})
	}
}

func TestSsgRepoCreateSectionImage(t *testing.T) {
	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		si          *ssg.SectionImage
		expectedErr error
	}{
		{
			name:        "creates section image successfully",
			setupFake:   func(f *fake.SsgRepo) {},
			si:          &ssg.SectionImage{ID: uuid.New(), SectionID: uuid.New(), ImageID: uuid.New()},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.CreateSectionImage(context.Background(), tt.si)

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

func TestSsgRepoDeleteSectionImage(t *testing.T) {
	siID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		id          uuid.UUID
		expectedErr error
	}{
		{
			name: "deletes section image successfully",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSectionImage(context.Background(), &ssg.SectionImage{ID: siID})
			},
			id:          siID,
			expectedErr: nil,
		},
		{
			name:        "returns nil when section image not found",
			setupFake:   func(f *fake.SsgRepo) {},
			id:          uuid.New(),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			err := f.DeleteSectionImage(context.Background(), tt.id)

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

func TestSsgRepoGetSectionImagesBySectionID(t *testing.T) {
	sectionID := uuid.New()

	tests := []struct {
		name        string
		setupFake   func(f *fake.SsgRepo)
		sectionID   uuid.UUID
		wantCount   int
		expectedErr error
	}{
		{
			name: "returns section images",
			setupFake: func(f *fake.SsgRepo) {
				f.CreateSectionImage(context.Background(), &ssg.SectionImage{ID: uuid.New(), SectionID: sectionID})
				f.CreateSectionImage(context.Background(), &ssg.SectionImage{ID: uuid.New(), SectionID: sectionID})
			},
			sectionID:   sectionID,
			wantCount:   2,
			expectedErr: nil,
		},
		{
			name:        "returns empty slice when no section images",
			setupFake:   func(f *fake.SsgRepo) {},
			sectionID:   uuid.New(),
			wantCount:   0,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fake.NewSsgRepo()
			tt.setupFake(f)

			sis, err := f.GetSectionImagesBySectionID(context.Background(), tt.sectionID)

			if tt.expectedErr != nil {
				if err == nil || err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if len(sis) != tt.wantCount {
				t.Errorf("expected %d section images, got %d", tt.wantCount, len(sis))
			}
		})
	}
}

func TestSsgRepoGetSectionWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetSectionFn = func(ctx context.Context, id uuid.UUID) (ssg.Section, error) {
		return ssg.Section{}, errors.New("custom error")
	}

	_, err := f.GetSection(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetLayoutWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetLayoutFn = func(ctx context.Context, id uuid.UUID) (ssg.Layout, error) {
		return ssg.Layout{}, errors.New("custom error")
	}

	_, err := f.GetLayout(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetTagWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetTagFn = func(ctx context.Context, id uuid.UUID) (ssg.Tag, error) {
		return ssg.Tag{}, errors.New("custom error")
	}

	_, err := f.GetTag(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateTagWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateTagFn = func(ctx context.Context, tag ssg.Tag) error {
		return errors.New("custom error")
	}

	err := f.UpdateTag(context.Background(), ssg.Tag{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetParamWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetParamFn = func(ctx context.Context, id uuid.UUID) (ssg.Param, error) {
		return ssg.Param{}, errors.New("custom error")
	}

	_, err := f.GetParam(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetImageFn = func(ctx context.Context, id uuid.UUID) (ssg.Image, error) {
		return ssg.Image{}, errors.New("custom error")
	}

	_, err := f.GetImage(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetImageVariantWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetImageVariantFn = func(ctx context.Context, id uuid.UUID) (ssg.ImageVariant, error) {
		return ssg.ImageVariant{}, errors.New("custom error")
	}

	_, err := f.GetImageVariant(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetUserByUsernameWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetUserByUsernameFn = func(ctx context.Context, username string) (auth.User, error) {
		return auth.User{}, errors.New("custom error")
	}

	_, err := f.GetUserByUsername(context.Background(), "testuser")
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetSiteBySlugWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetSiteBySlugFn = func(ctx context.Context, slug string) (ssg.Site, error) {
		return ssg.Site{}, errors.New("custom error")
	}

	_, err := f.GetSiteBySlug(context.Background(), "test-site")
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoListParamsWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.ListParamsFn = func(ctx context.Context) ([]ssg.Param, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.ListParams(context.Background())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateParamWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateParamFn = func(ctx context.Context, param *ssg.Param) error {
		return errors.New("custom error")
	}

	err := f.UpdateParam(context.Background(), &ssg.Param{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateImageFn = func(ctx context.Context, image *ssg.Image) error {
		return errors.New("custom error")
	}

	err := f.UpdateImage(context.Background(), &ssg.Image{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoListImagesWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.ListImagesFn = func(ctx context.Context) ([]ssg.Image, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.ListImages(context.Background())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoCreateImageVariantWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.CreateImageVariantFn = func(ctx context.Context, variant *ssg.ImageVariant) error {
		return errors.New("custom error")
	}

	err := f.CreateImageVariant(context.Background(), &ssg.ImageVariant{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateImageVariantWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateImageVariantFn = func(ctx context.Context, variant *ssg.ImageVariant) error {
		return errors.New("custom error")
	}

	err := f.UpdateImageVariant(context.Background(), &ssg.ImageVariant{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoDeleteImageVariantWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.DeleteImageVariantFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.DeleteImageVariant(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoListImageVariantsByImageIDWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.ListImageVariantsByImageIDFn = func(ctx context.Context, imageID uuid.UUID) ([]ssg.ImageVariant, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.ListImageVariantsByImageID(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoCreateContentImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.CreateContentImageFn = func(ctx context.Context, ci *ssg.ContentImage) error {
		return errors.New("custom error")
	}

	err := f.CreateContentImage(context.Background(), &ssg.ContentImage{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoDeleteContentImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.DeleteContentImageFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.DeleteContentImage(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetContentImagesByContentIDWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetContentImagesByContentIDFn = func(ctx context.Context, contentID uuid.UUID) ([]ssg.ContentImage, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.GetContentImagesByContentID(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoCreateSectionImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.CreateSectionImageFn = func(ctx context.Context, si *ssg.SectionImage) error {
		return errors.New("custom error")
	}

	err := f.CreateSectionImage(context.Background(), &ssg.SectionImage{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoDeleteSectionImageWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.DeleteSectionImageFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.DeleteSectionImage(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetSectionImagesBySectionIDWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetSectionImagesBySectionIDFn = func(ctx context.Context, sectionID uuid.UUID) ([]ssg.SectionImage, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.GetSectionImagesBySectionID(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoRemoveTagFromContentWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.RemoveTagFromContentFn = func(ctx context.Context, contentID, tagID uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.RemoveTagFromContent(context.Background(), uuid.New(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetSectionsWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetSectionsFn = func(ctx context.Context) ([]ssg.Section, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.GetSections(context.Background())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateSectionWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateSectionFn = func(ctx context.Context, section ssg.Section) error {
		return errors.New("custom error")
	}

	err := f.UpdateSection(context.Background(), ssg.Section{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoDeleteSectionWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.DeleteSectionFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.DeleteSection(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoCreateLayoutWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.CreateLayoutFn = func(ctx context.Context, layout ssg.Layout) error {
		return errors.New("custom error")
	}

	err := f.CreateLayout(context.Background(), ssg.Layout{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetAllLayoutsWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetAllLayoutsFn = func(ctx context.Context) ([]ssg.Layout, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.GetAllLayouts(context.Background())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoUpdateLayoutWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.UpdateLayoutFn = func(ctx context.Context, layout ssg.Layout) error {
		return errors.New("custom error")
	}

	err := f.UpdateLayout(context.Background(), ssg.Layout{})
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoDeleteLayoutWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.DeleteLayoutFn = func(ctx context.Context, id uuid.UUID) error {
		return errors.New("custom error")
	}

	err := f.DeleteLayout(context.Background(), uuid.New())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}

func TestSsgRepoGetAllTagsWithCustomFn(t *testing.T) {
	f := fake.NewSsgRepo()
	f.GetAllTagsFn = func(ctx context.Context) ([]ssg.Tag, error) {
		return nil, errors.New("custom error")
	}

	_, err := f.GetAllTags(context.Background())
	if err == nil || err.Error() != "custom error" {
		t.Errorf("expected custom error, got %v", err)
	}
}
