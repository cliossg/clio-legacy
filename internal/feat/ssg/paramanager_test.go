package ssg

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

// mockRepo implements only the Repo methods needed for ParamManager tests
type mockRepo struct {
	hm.Core
	params         map[string]Param
	getParamErr    error
	createParamErr error
	updateParamErr error
}

func newMockRepo() *mockRepo {
	// Create a minimal config for the mock
	cfg := hm.NewConfig()
	core := hm.NewCore("mock-repo", hm.XParams{Cfg: cfg})
	return &mockRepo{
		Core:   core,
		params: make(map[string]Param),
	}
}

func (m *mockRepo) GetParamByRefKey(ctx context.Context, refKey string) (Param, error) {
	if m.getParamErr != nil {
		return Param{}, m.getParamErr
	}
	param, exists := m.params[refKey]
	if !exists {
		return Param{}, fmt.Errorf("param not found")
	}
	return param, nil
}

func (m *mockRepo) GetParamByName(ctx context.Context, name string) (Param, error) {
	if m.getParamErr != nil {
		return Param{}, m.getParamErr
	}
	for _, param := range m.params {
		if param.Name == name {
			return param, nil
		}
	}
	return Param{}, fmt.Errorf("param not found")
}

func (m *mockRepo) CreateParam(ctx context.Context, param *Param) error {
	if m.createParamErr != nil {
		return m.createParamErr
	}
	m.params[param.RefKey] = *param
	return nil
}

func (m *mockRepo) UpdateParam(ctx context.Context, param *Param) error {
	if m.updateParamErr != nil {
		return m.updateParamErr
	}
	m.params[param.RefKey] = *param
	return nil
}

// Stub implementations for hm.Repo interface methods
func (m *mockRepo) Query() *hm.QueryManager {
	return nil
}

func (m *mockRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

// Stub implementations for unused Repo interface methods
func (m *mockRepo) CreateContent(ctx context.Context, content *Content) error { return nil }
func (m *mockRepo) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {
	return Content{}, nil
}
func (m *mockRepo) UpdateContent(ctx context.Context, content *Content) error { return nil }
func (m *mockRepo) DeleteContent(ctx context.Context, id uuid.UUID) error     { return nil }
func (m *mockRepo) GetAllContentWithMeta(ctx context.Context) ([]Content, error) {
	return nil, nil
}
func (m *mockRepo) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]Content, int, error) {
	return nil, 0, nil
}
func (m *mockRepo) CreateSection(ctx context.Context, section Section) error { return nil }
func (m *mockRepo) GetSection(ctx context.Context, id uuid.UUID) (Section, error) {
	return Section{}, nil
}
func (m *mockRepo) GetSections(ctx context.Context) ([]Section, error)   { return nil, nil }
func (m *mockRepo) UpdateSection(ctx context.Context, section Section) error { return nil }
func (m *mockRepo) DeleteSection(ctx context.Context, id uuid.UUID) error    { return nil }
func (m *mockRepo) CreateLayout(ctx context.Context, layout Layout) error    { return nil }
func (m *mockRepo) GetLayout(ctx context.Context, id uuid.UUID) (Layout, error) {
	return Layout{}, nil
}
func (m *mockRepo) GetAllLayouts(ctx context.Context) ([]Layout, error)   { return nil, nil }
func (m *mockRepo) UpdateLayout(ctx context.Context, layout Layout) error { return nil }
func (m *mockRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error  { return nil }
func (m *mockRepo) CreateTag(ctx context.Context, tag Tag) error          { return nil }
func (m *mockRepo) GetTag(ctx context.Context, id uuid.UUID) (Tag, error) {
	return Tag{}, nil
}
func (m *mockRepo) GetTagByName(ctx context.Context, name string) (Tag, error) {
	return Tag{}, nil
}
func (m *mockRepo) GetAllTags(ctx context.Context) ([]Tag, error)   { return nil, nil }
func (m *mockRepo) UpdateTag(ctx context.Context, tag Tag) error    { return nil }
func (m *mockRepo) DeleteTag(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) GetParam(ctx context.Context, id uuid.UUID) (Param, error) {
	return Param{}, nil
}
func (m *mockRepo) ListParams(ctx context.Context) ([]Param, error) { return nil, nil }
func (m *mockRepo) DeleteParam(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) CreateImage(ctx context.Context, image *Image) error { return nil }
func (m *mockRepo) GetImage(ctx context.Context, id uuid.UUID) (Image, error) {
	return Image{}, nil
}
func (m *mockRepo) GetImageByShortID(ctx context.Context, shortID string) (Image, error) {
	return Image{}, nil
}
func (m *mockRepo) GetImageByContentHash(ctx context.Context, contentHash string) (Image, error) {
	return Image{}, nil
}
func (m *mockRepo) UpdateImage(ctx context.Context, image *Image) error { return nil }
func (m *mockRepo) DeleteImage(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) ListImages(ctx context.Context) ([]Image, error)     { return nil, nil }
func (m *mockRepo) CreateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return nil
}
func (m *mockRepo) GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error) {
	return ImageVariant{}, nil
}
func (m *mockRepo) UpdateImageVariant(ctx context.Context, variant *ImageVariant) error {
	return nil
}
func (m *mockRepo) DeleteImageVariant(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error) {
	return nil, nil
}
func (m *mockRepo) CreateContentImage(ctx context.Context, contentImage *ContentImage) error {
	return nil
}
func (m *mockRepo) DeleteContentImage(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]ContentImage, error) {
	return nil, nil
}
func (m *mockRepo) CreateSectionImage(ctx context.Context, sectionImage *SectionImage) error {
	return nil
}
func (m *mockRepo) DeleteSectionImage(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockRepo) GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]SectionImage, error) {
	return nil, nil
}
func (m *mockRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	return nil
}
func (m *mockRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	return nil
}
func (m *mockRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error) {
	return nil, nil
}
func (m *mockRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error) {
	return nil, nil
}
func (m *mockRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	return auth.User{}, nil
}
func (m *mockRepo) GetSiteBySlug(ctx context.Context, slug string) (Site, error) {
	return Site{}, nil
}

func TestNewParamManager(t *testing.T) {
	repo := newMockRepo()
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	pm := NewParamManager(repo, params)

	if pm == nil {
		t.Fatal("NewParamManager() returned nil")
	}

	if pm.repo == nil {
		t.Error("repo was not set")
	}
}

func TestFindParam(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		setup     func(*mockRepo)
		wantErr   bool
	}{
		{
			name:      "finds existing param by name",
			paramName: "Test Param",
			setup: func(m *mockRepo) {
				param := Param{
					ID:     uuid.New(),
					Name:   "Test Param",
					RefKey: "test.param",
					Value:  "test-value",
				}
				m.params["test.param"] = param
			},
			wantErr: false,
		},
		{
			name:      "returns error for non-existent param",
			paramName: "Missing Param",
			setup:     func(m *mockRepo) {},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepo()
			tt.setup(repo)
			cfg := hm.NewConfig()
			pm := NewParamManager(repo, hm.XParams{Cfg: cfg})

			param, err := pm.FindParam(context.Background(), tt.paramName)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && param.Name != tt.paramName {
				t.Errorf("FindParam() param.Name = %v, want %v", param.Name, tt.paramName)
			}
		})
	}
}

func TestFindParamByRef(t *testing.T) {
	tests := []struct {
		name    string
		refKey  string
		setup   func(*mockRepo)
		wantErr bool
	}{
		{
			name:   "finds existing param by ref key",
			refKey: "site.mode",
			setup: func(m *mockRepo) {
				m.params["site.mode"] = Param{
					ID:     uuid.New(),
					RefKey: "site.mode",
					Value:  "structured",
				}
			},
			wantErr: false,
		},
		{
			name:    "returns error when repo is nil",
			refKey:  "any.key",
			setup:   func(m *mockRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pm *ParamManager
			cfg := hm.NewConfig()
			if tt.name == "returns error when repo is nil" {
				pm = NewParamManager(nil, hm.XParams{Cfg: cfg})
			} else {
				repo := newMockRepo()
				tt.setup(repo)
				pm = NewParamManager(repo, hm.XParams{Cfg: cfg})
			}

			param, err := pm.FindParamByRef(context.Background(), tt.refKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindParamByRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && param.RefKey != tt.refKey {
				t.Errorf("FindParamByRef() param.RefKey = %v, want %v", param.RefKey, tt.refKey)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		refKey  string
		defVal  string
		setup   func(*mockRepo)
		want    string
		useNilRepo bool
	}{
		{
			name:   "returns param value when found",
			refKey: "test.key",
			defVal: "default",
			setup: func(m *mockRepo) {
				m.params["test.key"] = Param{
					ID:     uuid.New(),
					RefKey: "test.key",
					Value:  "stored-value",
				}
			},
			want: "stored-value",
		},
		{
			name:   "returns default when param not found",
			refKey: "missing.key",
			defVal: "default-value",
			setup:  func(m *mockRepo) {},
			want:   "default-value",
		},
		{
			name:       "returns default when repo is nil",
			refKey:     "any.key",
			defVal:     "fallback",
			setup:      func(m *mockRepo) {},
			want:       "fallback",
			useNilRepo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pm *ParamManager
			cfg := hm.NewConfig()
			if tt.useNilRepo {
				pm = NewParamManager(nil, hm.XParams{Cfg: cfg})
			} else {
				repo := newMockRepo()
				tt.setup(repo)
				pm = NewParamManager(repo, hm.XParams{Cfg: cfg})
			}

			got := pm.Get(context.Background(), tt.refKey, tt.defVal)

			if got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteMode(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*mockRepo)
		want   string
	}{
		{
			name: "returns structured mode",
			setup: func(m *mockRepo) {
				m.params["site.mode"] = Param{
					ID:     uuid.New(),
					RefKey: "site.mode",
					Value:  "structured",
				}
			},
			want: "structured",
		},
		{
			name: "returns blog mode",
			setup: func(m *mockRepo) {
				m.params["site.mode"] = Param{
					ID:     uuid.New(),
					RefKey: "site.mode",
					Value:  "blog",
				}
			},
			want: "blog",
		},
		{
			name: "returns structured as default when not found",
			setup: func(m *mockRepo) {},
			want: "structured",
		},
		{
			name: "returns structured when mode is invalid",
			setup: func(m *mockRepo) {
				m.params["site.mode"] = Param{
					ID:     uuid.New(),
					RefKey: "site.mode",
					Value:  "invalid-mode",
				}
			},
			want: "structured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockRepo()
			tt.setup(repo)
			cfg := hm.NewConfig()
			pm := NewParamManager(repo, hm.XParams{Cfg: cfg})

			got := pm.GetSiteMode(context.Background())

			if got != tt.want {
				t.Errorf("GetSiteMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSiteMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		setup   func(*mockRepo)
		wantErr bool
		wantValue string
	}{
		{
			name:    "creates new param with structured mode",
			mode:    "structured",
			setup:   func(m *mockRepo) {},
			wantErr: false,
			wantValue: "structured",
		},
		{
			name:    "creates new param with blog mode",
			mode:    "blog",
			setup:   func(m *mockRepo) {},
			wantErr: false,
			wantValue: "blog",
		},
		{
			name: "updates existing param",
			mode: "blog",
			setup: func(m *mockRepo) {
				m.params["site.mode"] = Param{
					ID:     uuid.New(),
					RefKey: "site.mode",
					Value:  "structured",
				}
			},
			wantErr: false,
			wantValue: "blog",
		},
		{
			name:    "returns error for invalid mode",
			mode:    "invalid",
			setup:   func(m *mockRepo) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pm *ParamManager
			cfg := hm.NewConfig()
			if tt.name == "returns error when repo is nil" {
				pm = NewParamManager(nil, hm.XParams{Cfg: cfg})
			} else {
				repo := newMockRepo()
				tt.setup(repo)
				pm = NewParamManager(repo, hm.XParams{Cfg: cfg})
			}

			err := pm.SetSiteMode(context.Background(), tt.mode)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetSiteMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the value was set correctly
				repo := pm.repo.(*mockRepo)
				param, exists := repo.params["site.mode"]
				if !exists {
					t.Error("site.mode param was not created")
					return
				}
				if param.Value != tt.wantValue {
					t.Errorf("SetSiteMode() set value = %v, want %v", param.Value, tt.wantValue)
				}
			}
		})
	}
}

func TestSetSiteModeWithNilRepo(t *testing.T) {
	cfg := hm.NewConfig()
	pm := NewParamManager(nil, hm.XParams{Cfg: cfg})
	err := pm.SetSiteMode(context.Background(), "structured")

	if err == nil {
		t.Error("SetSiteMode() with nil repo should return error")
	}
}
