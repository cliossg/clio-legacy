package ssg

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/hm"
)

type mockServiceRepo struct {
	hm.Core
	contents        map[uuid.UUID]Content
	sections        map[uuid.UUID]Section
	layouts         map[uuid.UUID]Layout
	tags            map[uuid.UUID]Tag
	tagsByName      map[string]Tag
	params          map[uuid.UUID]Param
	paramsByName    map[string]Param
	paramsByRef     map[string]Param
	images          map[uuid.UUID]Image
	imagesByShortID map[string]Image
	imageVariants   map[uuid.UUID]ImageVariant
	contentImages   map[uuid.UUID][]ContentImage
	sectionImages   map[uuid.UUID][]SectionImage
	contentTags     map[uuid.UUID][]Tag
	tagContent      map[uuid.UUID][]Content

	createContentErr error
	getContentErr    error
	updateContentErr error
	deleteContentErr error

	createSectionErr error
	getSectionErr    error
	updateSectionErr error
	deleteSectionErr error

	createLayoutErr error
	getLayoutErr    error
	updateLayoutErr error
	deleteLayoutErr error

	createTagErr error
	getTagErr    error
	updateTagErr error
	deleteTagErr error

	createParamErr error
	getParamErr    error
	updateParamErr error
	deleteParamErr error

	createImageErr error
	getImageErr    error
	updateImageErr error
	deleteImageErr error

	createImageVariantErr error
	getImageVariantErr    error
	updateImageVariantErr error
	deleteImageVariantErr error

	addTagToContentErr               error
	removeTagFromContentErr          error
	getTagsForContentErr             error
	getContentForTagErr              error
	createContentImageErr            error
	getContentImagesByContentIDErr   error
	createSectionImageErr            error
	getSectionImagesBySectionIDErr   error
	deleteContentImageErr            error
	deleteSectionImageErr            error

	createTagCalled      bool
	addTagToContentCalled bool
}

func newMockServiceRepo() *mockServiceRepo {
	cfg := hm.NewConfig()
	core := hm.NewCore("mock-service-repo", hm.XParams{Cfg: cfg})
	return &mockServiceRepo{
		Core:            core,
		contents:        make(map[uuid.UUID]Content),
		sections:        make(map[uuid.UUID]Section),
		layouts:         make(map[uuid.UUID]Layout),
		tags:            make(map[uuid.UUID]Tag),
		tagsByName:      make(map[string]Tag),
		params:          make(map[uuid.UUID]Param),
		paramsByName:    make(map[string]Param),
		paramsByRef:     make(map[string]Param),
		images:          make(map[uuid.UUID]Image),
		imagesByShortID: make(map[string]Image),
		imageVariants:   make(map[uuid.UUID]ImageVariant),
		contentImages:   make(map[uuid.UUID][]ContentImage),
		sectionImages:   make(map[uuid.UUID][]SectionImage),
		contentTags:     make(map[uuid.UUID][]Tag),
		tagContent:      make(map[uuid.UUID][]Content),
	}
}

func (m *mockServiceRepo) Query() *hm.QueryManager {
	return nil
}

func (m *mockServiceRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

func (m *mockServiceRepo) CreateContent(ctx context.Context, content *Content) error {
	if m.createContentErr != nil {
		return m.createContentErr
	}
	m.contents[content.ID] = *content
	return nil
}

func (m *mockServiceRepo) GetContent(ctx context.Context, id uuid.UUID) (Content, error) {
	if m.getContentErr != nil {
		return Content{}, m.getContentErr
	}
	content, exists := m.contents[id]
	if !exists {
		return Content{}, sql.ErrNoRows
	}
	return content, nil
}

func (m *mockServiceRepo) UpdateContent(ctx context.Context, content *Content) error {
	if m.updateContentErr != nil {
		return m.updateContentErr
	}
	m.contents[content.ID] = *content
	return nil
}

func (m *mockServiceRepo) DeleteContent(ctx context.Context, id uuid.UUID) error {
	if m.deleteContentErr != nil {
		return m.deleteContentErr
	}
	delete(m.contents, id)
	return nil
}

func (m *mockServiceRepo) GetAllContentWithMeta(ctx context.Context) ([]Content, error) {
	if m.getContentErr != nil {
		return nil, m.getContentErr
	}
	result := make([]Content, 0, len(m.contents))
	for _, c := range m.contents {
		result = append(result, c)
	}
	return result, nil
}

func (m *mockServiceRepo) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]Content, int, error) {
	if m.getContentErr != nil {
		return nil, 0, m.getContentErr
	}
	all := make([]Content, 0, len(m.contents))
	for _, c := range m.contents {
		all = append(all, c)
	}
	total := len(all)

	if offset >= total {
		return []Content{}, total, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return all[offset:end], total, nil
}

func (m *mockServiceRepo) CreateSection(ctx context.Context, section Section) error {
	if m.createSectionErr != nil {
		return m.createSectionErr
	}
	m.sections[section.ID] = section
	return nil
}

func (m *mockServiceRepo) GetSection(ctx context.Context, id uuid.UUID) (Section, error) {
	if m.getSectionErr != nil {
		return Section{}, m.getSectionErr
	}
	section, exists := m.sections[id]
	if !exists {
		return Section{}, sql.ErrNoRows
	}
	return section, nil
}

func (m *mockServiceRepo) GetSections(ctx context.Context) ([]Section, error) {
	if m.getSectionErr != nil {
		return nil, m.getSectionErr
	}
	result := make([]Section, 0, len(m.sections))
	for _, s := range m.sections {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateSection(ctx context.Context, section Section) error {
	if m.updateSectionErr != nil {
		return m.updateSectionErr
	}
	m.sections[section.ID] = section
	return nil
}

func (m *mockServiceRepo) DeleteSection(ctx context.Context, id uuid.UUID) error {
	if m.deleteSectionErr != nil {
		return m.deleteSectionErr
	}
	delete(m.sections, id)
	return nil
}

func (m *mockServiceRepo) CreateLayout(ctx context.Context, layout Layout) error {
	if m.createLayoutErr != nil {
		return m.createLayoutErr
	}
	m.layouts[layout.ID] = layout
	return nil
}

func (m *mockServiceRepo) GetLayout(ctx context.Context, id uuid.UUID) (Layout, error) {
	if m.getLayoutErr != nil {
		return Layout{}, m.getLayoutErr
	}
	layout, exists := m.layouts[id]
	if !exists {
		return Layout{}, sql.ErrNoRows
	}
	return layout, nil
}

func (m *mockServiceRepo) GetAllLayouts(ctx context.Context) ([]Layout, error) {
	if m.getLayoutErr != nil {
		return nil, m.getLayoutErr
	}
	result := make([]Layout, 0, len(m.layouts))
	for _, l := range m.layouts {
		result = append(result, l)
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateLayout(ctx context.Context, layout Layout) error {
	if m.updateLayoutErr != nil {
		return m.updateLayoutErr
	}
	m.layouts[layout.ID] = layout
	return nil
}

func (m *mockServiceRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	if m.deleteLayoutErr != nil {
		return m.deleteLayoutErr
	}
	delete(m.layouts, id)
	return nil
}

func (m *mockServiceRepo) CreateTag(ctx context.Context, tag Tag) error {
	m.createTagCalled = true
	if m.createTagErr != nil {
		return m.createTagErr
	}
	m.tags[tag.ID] = tag
	m.tagsByName[tag.Name] = tag
	return nil
}

func (m *mockServiceRepo) GetTag(ctx context.Context, id uuid.UUID) (Tag, error) {
	if m.getTagErr != nil {
		return Tag{}, m.getTagErr
	}
	tag, exists := m.tags[id]
	if !exists {
		return Tag{}, sql.ErrNoRows
	}
	return tag, nil
}

func (m *mockServiceRepo) GetTagByName(ctx context.Context, name string) (Tag, error) {
	if m.getTagErr != nil {
		return Tag{}, m.getTagErr
	}
	tag, exists := m.tagsByName[name]
	if !exists {
		return Tag{}, sql.ErrNoRows
	}
	return tag, nil
}

func (m *mockServiceRepo) GetAllTags(ctx context.Context) ([]Tag, error) {
	if m.getTagErr != nil {
		return nil, m.getTagErr
	}
	result := make([]Tag, 0, len(m.tags))
	for _, t := range m.tags {
		result = append(result, t)
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateTag(ctx context.Context, tag Tag) error {
	if m.updateTagErr != nil {
		return m.updateTagErr
	}
	m.tags[tag.ID] = tag
	m.tagsByName[tag.Name] = tag
	return nil
}

func (m *mockServiceRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	if m.deleteTagErr != nil {
		return m.deleteTagErr
	}
	tag := m.tags[id]
	delete(m.tags, id)
	delete(m.tagsByName, tag.Name)
	return nil
}

func (m *mockServiceRepo) CreateParam(ctx context.Context, param *Param) error {
	if m.createParamErr != nil {
		return m.createParamErr
	}
	m.params[param.ID] = *param
	m.paramsByName[param.Name] = *param
	m.paramsByRef[param.RefKey] = *param
	return nil
}

func (m *mockServiceRepo) GetParam(ctx context.Context, id uuid.UUID) (Param, error) {
	if m.getParamErr != nil {
		return Param{}, m.getParamErr
	}
	param, exists := m.params[id]
	if !exists {
		return Param{}, sql.ErrNoRows
	}
	return param, nil
}

func (m *mockServiceRepo) GetParamByName(ctx context.Context, name string) (Param, error) {
	if m.getParamErr != nil {
		return Param{}, m.getParamErr
	}
	param, exists := m.paramsByName[name]
	if !exists {
		return Param{}, sql.ErrNoRows
	}
	return param, nil
}

func (m *mockServiceRepo) GetParamByRefKey(ctx context.Context, refKey string) (Param, error) {
	if m.getParamErr != nil {
		return Param{}, m.getParamErr
	}
	param, exists := m.paramsByRef[refKey]
	if !exists {
		return Param{}, sql.ErrNoRows
	}
	return param, nil
}

func (m *mockServiceRepo) ListParams(ctx context.Context) ([]Param, error) {
	if m.getParamErr != nil {
		return nil, m.getParamErr
	}
	result := make([]Param, 0, len(m.params))
	for _, p := range m.params {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateParam(ctx context.Context, param *Param) error {
	if m.updateParamErr != nil {
		return m.updateParamErr
	}
	m.params[param.ID] = *param
	m.paramsByName[param.Name] = *param
	m.paramsByRef[param.RefKey] = *param
	return nil
}

func (m *mockServiceRepo) DeleteParam(ctx context.Context, id uuid.UUID) error {
	if m.deleteParamErr != nil {
		return m.deleteParamErr
	}
	param := m.params[id]
	delete(m.params, id)
	delete(m.paramsByName, param.Name)
	delete(m.paramsByRef, param.RefKey)
	return nil
}

func (m *mockServiceRepo) CreateImage(ctx context.Context, image *Image) error {
	if m.createImageErr != nil {
		return m.createImageErr
	}
	m.images[image.ID] = *image
	m.imagesByShortID[image.ShortID] = *image
	return nil
}

func (m *mockServiceRepo) GetImage(ctx context.Context, id uuid.UUID) (Image, error) {
	if m.getImageErr != nil {
		return Image{}, m.getImageErr
	}
	image, exists := m.images[id]
	if !exists {
		return Image{}, sql.ErrNoRows
	}
	return image, nil
}

func (m *mockServiceRepo) GetImageByShortID(ctx context.Context, shortID string) (Image, error) {
	if m.getImageErr != nil {
		return Image{}, m.getImageErr
	}
	image, exists := m.imagesByShortID[shortID]
	if !exists {
		return Image{}, sql.ErrNoRows
	}
	return image, nil
}

func (m *mockServiceRepo) GetImageByContentHash(ctx context.Context, contentHash string) (Image, error) {
	return Image{}, nil
}

func (m *mockServiceRepo) ListImages(ctx context.Context) ([]Image, error) {
	if m.getImageErr != nil {
		return nil, m.getImageErr
	}
	result := make([]Image, 0, len(m.images))
	for _, img := range m.images {
		result = append(result, img)
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateImage(ctx context.Context, image *Image) error {
	if m.updateImageErr != nil {
		return m.updateImageErr
	}
	m.images[image.ID] = *image
	m.imagesByShortID[image.ShortID] = *image
	return nil
}

func (m *mockServiceRepo) DeleteImage(ctx context.Context, id uuid.UUID) error {
	if m.deleteImageErr != nil {
		return m.deleteImageErr
	}
	image := m.images[id]
	delete(m.images, id)
	delete(m.imagesByShortID, image.ShortID)
	return nil
}

func (m *mockServiceRepo) CreateImageVariant(ctx context.Context, variant *ImageVariant) error {
	if m.createImageVariantErr != nil {
		return m.createImageVariantErr
	}
	m.imageVariants[variant.ID] = *variant
	return nil
}

func (m *mockServiceRepo) GetImageVariant(ctx context.Context, id uuid.UUID) (ImageVariant, error) {
	if m.getImageVariantErr != nil {
		return ImageVariant{}, m.getImageVariantErr
	}
	variant, exists := m.imageVariants[id]
	if !exists {
		return ImageVariant{}, sql.ErrNoRows
	}
	return variant, nil
}

func (m *mockServiceRepo) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ImageVariant, error) {
	if m.getImageVariantErr != nil {
		return nil, m.getImageVariantErr
	}
	result := make([]ImageVariant, 0)
	for _, v := range m.imageVariants {
		if v.ImageID == imageID {
			result = append(result, v)
		}
	}
	return result, nil
}

func (m *mockServiceRepo) UpdateImageVariant(ctx context.Context, variant *ImageVariant) error {
	if m.updateImageVariantErr != nil {
		return m.updateImageVariantErr
	}
	m.imageVariants[variant.ID] = *variant
	return nil
}

func (m *mockServiceRepo) DeleteImageVariant(ctx context.Context, id uuid.UUID) error {
	if m.deleteImageVariantErr != nil {
		return m.deleteImageVariantErr
	}
	delete(m.imageVariants, id)
	return nil
}

func (m *mockServiceRepo) CreateContentImage(ctx context.Context, contentImage *ContentImage) error {
	if m.createContentImageErr != nil {
		return m.createContentImageErr
	}
	return nil
}

func (m *mockServiceRepo) DeleteContentImage(ctx context.Context, id uuid.UUID) error {
	if m.deleteContentImageErr != nil {
		return m.deleteContentImageErr
	}
	return nil
}

func (m *mockServiceRepo) GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]ContentImage, error) {
	if m.getContentImagesByContentIDErr != nil {
		return nil, m.getContentImagesByContentIDErr
	}
	images, exists := m.contentImages[contentID]
	if !exists {
		return []ContentImage{}, nil
	}
	return images, nil
}

func (m *mockServiceRepo) CreateSectionImage(ctx context.Context, sectionImage *SectionImage) error {
	if m.createSectionImageErr != nil {
		return m.createSectionImageErr
	}
	return nil
}

func (m *mockServiceRepo) DeleteSectionImage(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockServiceRepo) GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]SectionImage, error) {
	if m.getSectionImagesBySectionIDErr != nil {
		return nil, m.getSectionImagesBySectionIDErr
	}
	images, exists := m.sectionImages[sectionID]
	if !exists {
		return []SectionImage{}, nil
	}
	return images, nil
}

func (m *mockServiceRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	m.addTagToContentCalled = true
	if m.addTagToContentErr != nil {
		return m.addTagToContentErr
	}
	return nil
}

func (m *mockServiceRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	if m.removeTagFromContentErr != nil {
		return m.removeTagFromContentErr
	}
	return nil
}

func (m *mockServiceRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]Tag, error) {
	if m.getTagsForContentErr != nil {
		return nil, m.getTagsForContentErr
	}
	tags, exists := m.contentTags[contentID]
	if !exists {
		return []Tag{}, nil
	}
	return tags, nil
}

func (m *mockServiceRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]Content, error) {
	if m.getContentForTagErr != nil {
		return nil, m.getContentForTagErr
	}
	content, exists := m.tagContent[tagID]
	if !exists {
		return []Content{}, nil
	}
	return content, nil
}

func (m *mockServiceRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	return auth.User{}, nil
}

func (m *mockServiceRepo) GetSiteBySlug(ctx context.Context, slug string) (Site, error) {
	return Site{}, nil
}

type mockPublisher struct{}

func (m *mockPublisher) Validate(cfg PublisherConfig) error {
	return nil
}

func (m *mockPublisher) Publish(ctx context.Context, cfg PublisherConfig, sourceDir string) (string, error) {
	return "", nil
}

func (m *mockPublisher) Plan(ctx context.Context, cfg PublisherConfig, sourceDir string) (PlanReport, error) {
	return PlanReport{}, nil
}

func newTestService(repo Repo) *BaseService {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	pub := &mockPublisher{}
	pm := NewParamManager(repo, params)

	return NewService(embed.FS{}, repo, nil, pub, pm, nil, params)
}

func TestServiceCreateContent(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		content *Content
		wantErr bool
	}{
		{
			name:  "creates content successfully",
			setup: func(m *mockServiceRepo) {},
			content: &Content{
				ID:      uuid.New(),
				Heading: "Test Content",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createContentErr = fmt.Errorf("db error")
			},
			content: &Content{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateContent(context.Background(), tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateContentWithNilRepo(t *testing.T) {
	svc := newTestService(nil)
	err := svc.CreateContent(context.Background(), &Content{ID: uuid.New()})

	if err == nil {
		t.Error("CreateContent() with nil repo should return error")
	}
}

func TestServiceGetContent(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets content successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when content not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getContentErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetContent(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceUpdateContent(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		content *Content
		wantErr bool
	}{
		{
			name:  "updates content successfully",
			setup: func(m *mockServiceRepo) {},
			content: &Content{
				ID:      uuid.New(),
				Heading: "Updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateContentErr = fmt.Errorf("db error")
			},
			content: &Content{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateContent(context.Background(), tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteContent(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes content successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteContentErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteContent(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetAllContentWithMeta(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "gets all content successfully",
			setup: func(m *mockServiceRepo) {
				m.contents[uuid.New()] = Content{Heading: "First"}
				m.contents[uuid.New()] = Content{Heading: "Second"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.GetAllContentWithMeta(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllContentWithMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("GetAllContentWithMeta() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceGetContentWithPaginationAndSearch(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo)
		offset    int
		limit     int
		wantCount int
		wantTotal int
		wantErr   bool
	}{
		{
			name: "gets paginated content successfully",
			setup: func(m *mockServiceRepo) {
				for i := 0; i < 10; i++ {
					m.contents[uuid.New()] = Content{}
				}
			},
			offset:    0,
			limit:     5,
			wantCount: 5,
			wantTotal: 10,
			wantErr:   false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
			},
			offset:  0,
			limit:   5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, total, err := svc.GetContentWithPaginationAndSearch(context.Background(), tt.offset, tt.limit, "")

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentWithPaginationAndSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(result) != tt.wantCount {
					t.Errorf("GetContentWithPaginationAndSearch() count = %d, want %d", len(result), tt.wantCount)
				}
				if total != tt.wantTotal {
					t.Errorf("GetContentWithPaginationAndSearch() total = %d, want %d", total, tt.wantTotal)
				}
			}
		})
	}
}

func TestServiceCreateSection(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		section Section
		wantErr bool
	}{
		{
			name:  "creates section successfully",
			setup: func(m *mockServiceRepo) {},
			section: Section{
				ID:   uuid.New(),
				Path: "test-section",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createSectionErr = fmt.Errorf("db error")
			},
			section: Section{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateSection(context.Background(), tt.section)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetSection(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets section successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.sections[id] = Section{ID: id, Path: "test"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when section not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getSectionErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetSection(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetSections(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "gets all sections successfully",
			setup: func(m *mockServiceRepo) {
				m.sections[uuid.New()] = Section{Path: "first"}
				m.sections[uuid.New()] = Section{Path: "second"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getSectionErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.GetSections(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("GetSections() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateSection(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		section Section
		wantErr bool
	}{
		{
			name:  "updates section successfully",
			setup: func(m *mockServiceRepo) {},
			section: Section{
				ID:   uuid.New(),
				Path: "updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateSectionErr = fmt.Errorf("db error")
			},
			section: Section{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateSection(context.Background(), tt.section)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteSection(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes section successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteSectionErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteSection(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateLayout(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		layout  Layout
		wantErr bool
	}{
		{
			name:  "creates layout successfully",
			setup: func(m *mockServiceRepo) {},
			layout: Layout{
				ID:   uuid.New(),
				Name: "test-layout",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createLayoutErr = fmt.Errorf("db error")
			},
			layout:  Layout{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateLayout(context.Background(), tt.layout)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetLayout(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets layout successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.layouts[id] = Layout{ID: id, Name: "test"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when layout not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getLayoutErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetLayout(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetAllLayouts(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "gets all layouts successfully",
			setup: func(m *mockServiceRepo) {
				m.layouts[uuid.New()] = Layout{Name: "first"}
				m.layouts[uuid.New()] = Layout{Name: "second"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getLayoutErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.GetAllLayouts(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllLayouts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("GetAllLayouts() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateLayout(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		layout  Layout
		wantErr bool
	}{
		{
			name:  "updates layout successfully",
			setup: func(m *mockServiceRepo) {},
			layout: Layout{
				ID:   uuid.New(),
				Name: "updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateLayoutErr = fmt.Errorf("db error")
			},
			layout:  Layout{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateLayout(context.Background(), tt.layout)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteLayout(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes layout successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteLayoutErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteLayout(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateTag(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		tag     Tag
		wantErr bool
	}{
		{
			name:  "creates tag successfully",
			setup: func(m *mockServiceRepo) {},
			tag: Tag{
				ID:   uuid.New(),
				Name: "test-tag",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createTagErr = fmt.Errorf("db error")
			},
			tag:     Tag{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateTag(context.Background(), tt.tag)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetTag(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets tag successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.tags[id] = Tag{ID: id, Name: "test"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when tag not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getTagErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetTag(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetTagByName(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		tagName string
		wantErr bool
	}{
		{
			name: "gets tag by name successfully",
			setup: func(m *mockServiceRepo) {
				tag := Tag{ID: uuid.New(), Name: "test-tag"}
				m.tags[tag.ID] = tag
				m.tagsByName["test-tag"] = tag
			},
			tagName: "test-tag",
			wantErr: false,
		},
		{
			name:    "returns error when tag not found",
			setup:   func(m *mockServiceRepo) {},
			tagName: "missing-tag",
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getTagErr = fmt.Errorf("db error")
			},
			tagName: "any-tag",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetTagByName(context.Background(), tt.tagName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetAllTags(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "gets all tags successfully",
			setup: func(m *mockServiceRepo) {
				m.tags[uuid.New()] = Tag{Name: "first"}
				m.tags[uuid.New()] = Tag{Name: "second"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getTagErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.GetAllTags(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("GetAllTags() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateTag(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		tag     Tag
		wantErr bool
	}{
		{
			name:  "updates tag successfully",
			setup: func(m *mockServiceRepo) {},
			tag: Tag{
				ID:   uuid.New(),
				Name: "updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateTagErr = fmt.Errorf("db error")
			},
			tag:     Tag{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateTag(context.Background(), tt.tag)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteTag(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes tag successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteTagErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteTag(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateParam(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		param   *Param
		wantErr bool
	}{
		{
			name:  "creates param successfully",
			setup: func(m *mockServiceRepo) {},
			param: &Param{
				ID:     uuid.New(),
				Name:   "test-param",
				RefKey: "test.param",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createParamErr = fmt.Errorf("db error")
			},
			param:   &Param{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateParam(context.Background(), tt.param)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetParam(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets param successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.params[id] = Param{ID: id, Name: "test"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when param not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getParamErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetParam(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetParamByName(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*mockServiceRepo)
		paramName string
		wantErr   bool
	}{
		{
			name: "gets param by name successfully",
			setup: func(m *mockServiceRepo) {
				param := Param{ID: uuid.New(), Name: "test-param"}
				m.params[param.ID] = param
				m.paramsByName["test-param"] = param
			},
			paramName: "test-param",
			wantErr:   false,
		},
		{
			name:      "returns error when param not found",
			setup:     func(m *mockServiceRepo) {},
			paramName: "missing-param",
			wantErr:   true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("db error")
			},
			paramName: "any-param",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetParamByName(context.Background(), tt.paramName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetParamByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetParamByRefKey(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		refKey  string
		wantErr bool
	}{
		{
			name: "gets param by ref key successfully",
			setup: func(m *mockServiceRepo) {
				param := Param{ID: uuid.New(), RefKey: "test.param"}
				m.params[param.ID] = param
				m.paramsByRef["test.param"] = param
			},
			refKey:  "test.param",
			wantErr: false,
		},
		{
			name:    "returns error when param not found",
			setup:   func(m *mockServiceRepo) {},
			refKey:  "missing.param",
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("db error")
			},
			refKey:  "any.param",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetParamByRefKey(context.Background(), tt.refKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetParamByRefKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceListParams(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "lists params successfully",
			setup: func(m *mockServiceRepo) {
				m.params[uuid.New()] = Param{Name: "first"}
				m.params[uuid.New()] = Param{Name: "second"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getParamErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.ListParams(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("ListParams() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateParam(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		param   *Param
		wantErr bool
	}{
		{
			name:  "updates param successfully",
			setup: func(m *mockServiceRepo) {},
			param: &Param{
				ID:   uuid.New(),
				Name: "updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateParamErr = fmt.Errorf("db error")
			},
			param:   &Param{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateParam(context.Background(), tt.param)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteParam(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes param successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteParamErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteParam(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateImage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		image   *Image
		wantErr bool
	}{
		{
			name:  "creates image successfully",
			setup: func(m *mockServiceRepo) {},
			image: &Image{
				ID:       uuid.New(),
				FileName: "test.jpg",
				ShortID:  "abc123",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createImageErr = fmt.Errorf("db error")
			},
			image:   &Image{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateImage(context.Background(), tt.image)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetImage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets image successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.images[id] = Image{ID: id, FileName: "test.jpg"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when image not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getImageErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetImage(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetImageByShortID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		shortID string
		wantErr bool
	}{
		{
			name: "gets image by short ID successfully",
			setup: func(m *mockServiceRepo) {
				img := Image{ID: uuid.New(), ShortID: "abc123"}
				m.images[img.ID] = img
				m.imagesByShortID["abc123"] = img
			},
			shortID: "abc123",
			wantErr: false,
		},
		{
			name:    "returns error when image not found",
			setup:   func(m *mockServiceRepo) {},
			shortID: "missing",
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getImageErr = fmt.Errorf("db error")
			},
			shortID: "any",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetImageByShortID(context.Background(), tt.shortID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageByShortID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceListImages(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantLen int
		wantErr bool
	}{
		{
			name: "lists images successfully",
			setup: func(m *mockServiceRepo) {
				m.images[uuid.New()] = Image{FileName: "first.jpg"}
				m.images[uuid.New()] = Image{FileName: "second.jpg"}
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.getImageErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.ListImages(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("ListImages() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateImage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		image   *Image
		wantErr bool
	}{
		{
			name:  "updates image successfully",
			setup: func(m *mockServiceRepo) {},
			image: &Image{
				ID:       uuid.New(),
				FileName: "updated.jpg",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateImageErr = fmt.Errorf("db error")
			},
			image:   &Image{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateImage(context.Background(), tt.image)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteImage(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes image successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteImageErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteImage(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceCreateImageVariant(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		variant *ImageVariant
		wantErr bool
	}{
		{
			name:  "creates image variant successfully",
			setup: func(m *mockServiceRepo) {},
			variant: &ImageVariant{
				ID:      uuid.New(),
				ImageID: uuid.New(),
				Kind:    "thumbnail",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.createImageVariantErr = fmt.Errorf("db error")
			},
			variant: &ImageVariant{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.CreateImageVariant(context.Background(), tt.variant)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetImageVariant(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantErr bool
	}{
		{
			name: "gets image variant successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.imageVariants[id] = ImageVariant{ID: id, Kind: "thumbnail"}
				return id
			},
			wantErr: false,
		},
		{
			name: "returns error when variant not found",
			setup: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getImageVariantErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			id := tt.setup(repo)
			svc := newTestService(repo)

			_, err := svc.GetImageVariant(context.Background(), id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceListImageVariantsByImageID(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo) uuid.UUID
		wantLen int
		wantErr bool
	}{
		{
			name: "lists image variants successfully",
			setup: func(m *mockServiceRepo) uuid.UUID {
				imageID := uuid.New()
				m.imageVariants[uuid.New()] = ImageVariant{ImageID: imageID, Kind: "thumbnail"}
				m.imageVariants[uuid.New()] = ImageVariant{ImageID: imageID, Kind: "small"}
				m.imageVariants[uuid.New()] = ImageVariant{ImageID: uuid.New(), Kind: "other"}
				return imageID
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) uuid.UUID {
				m.getImageVariantErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			imageID := tt.setup(repo)
			svc := newTestService(repo)

			result, err := svc.ListImageVariantsByImageID(context.Background(), imageID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListImageVariantsByImageID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(result) != tt.wantLen {
				t.Errorf("ListImageVariantsByImageID() len = %d, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestServiceUpdateImageVariant(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		variant *ImageVariant
		wantErr bool
	}{
		{
			name:  "updates image variant successfully",
			setup: func(m *mockServiceRepo) {},
			variant: &ImageVariant{
				ID:   uuid.New(),
				Kind: "updated",
			},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.updateImageVariantErr = fmt.Errorf("db error")
			},
			variant: &ImageVariant{ID: uuid.New()},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.UpdateImageVariant(context.Background(), tt.variant)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDeleteImageVariant(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*mockServiceRepo)
		wantErr bool
	}{
		{
			name:    "deletes image variant successfully",
			setup:   func(m *mockServiceRepo) {},
			wantErr: false,
		},
		{
			name: "returns error when repo fails",
			setup: func(m *mockServiceRepo) {
				m.deleteImageVariantErr = fmt.Errorf("db error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setup(repo)
			svc := newTestService(repo)

			err := svc.DeleteImageVariant(context.Background(), uuid.New())

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
