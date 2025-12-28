package ssg

import (
	"context"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
)

type testRepo struct {
	hm.Core
	siteRepo feat.SiteRepo
	db       *sqlx.DB
}

func newTestRepo(siteRepo feat.SiteRepo) *testRepo {
	cfg := hm.NewConfig()
	return &testRepo{
		Core:     hm.NewCore("test-repo", hm.XParams{Cfg: cfg}),
		siteRepo: siteRepo,
	}
}

func newTestRepoWithDB(siteRepo feat.SiteRepo, db *sqlx.DB) *testRepo {
	cfg := hm.NewConfig()
	return &testRepo{
		Core:     hm.NewCore("test-repo", hm.XParams{Cfg: cfg}),
		siteRepo: siteRepo,
		db:       db,
	}
}

// GetDB returns the test DB
func (r *testRepo) GetDB() *sqlx.DB { return r.db }

// Query returns nil (not used in site tests)
func (r *testRepo) Query() *hm.QueryManager { return nil }

// BeginTx returns a stub transaction
func (r *testRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

// Site methods - delegate to siteRepo
func (r *testRepo) CreateSite(ctx context.Context, site *feat.Site) error {
	return r.siteRepo.CreateSite(ctx, site)
}
func (r *testRepo) GetSite(ctx context.Context, id uuid.UUID) (feat.Site, error) {
	return r.siteRepo.GetSite(ctx, id)
}
func (r *testRepo) GetSiteBySlug(ctx context.Context, slug string) (feat.Site, error) {
	return r.siteRepo.GetSiteBySlug(ctx, slug)
}
func (r *testRepo) ListSites(ctx context.Context, activeOnly bool) ([]feat.Site, error) {
	return r.siteRepo.ListSites(ctx, activeOnly)
}
func (r *testRepo) UpdateSite(ctx context.Context, site *feat.Site) error {
	return r.siteRepo.UpdateSite(ctx, site)
}
func (r *testRepo) DeleteSite(ctx context.Context, id uuid.UUID) error {
	return r.siteRepo.DeleteSite(ctx, id)
}

// All other methods - stubs (not used in site handler tests)
func (r *testRepo) CreateContent(ctx context.Context, content *feat.Content) error { return nil }
func (r *testRepo) GetContent(ctx context.Context, id uuid.UUID) (feat.Content, error) {
	return feat.Content{}, nil
}
func (r *testRepo) UpdateContent(ctx context.Context, content *feat.Content) error                       { return nil }
func (r *testRepo) DeleteContent(ctx context.Context, id uuid.UUID) error                                 { return nil }
func (r *testRepo) GetAllContentWithMeta(ctx context.Context) ([]feat.Content, error)                     { return nil, nil }
func (r *testRepo) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]feat.Content, int, error) {
	return nil, 0, nil
}
func (r *testRepo) CreateSection(ctx context.Context, section feat.Section) error { return nil }
func (r *testRepo) GetSection(ctx context.Context, id uuid.UUID) (feat.Section, error) {
	return feat.Section{}, nil
}
func (r *testRepo) GetSections(ctx context.Context) ([]feat.Section, error) { return nil, nil }
func (r *testRepo) UpdateSection(ctx context.Context, section feat.Section) error { return nil }
func (r *testRepo) DeleteSection(ctx context.Context, id uuid.UUID) error          { return nil }
func (r *testRepo) CreateLayout(ctx context.Context, layout feat.Layout) error     { return nil }
func (r *testRepo) GetLayout(ctx context.Context, id uuid.UUID) (feat.Layout, error) {
	return feat.Layout{}, nil
}
func (r *testRepo) GetAllLayouts(ctx context.Context) ([]feat.Layout, error)            { return nil, nil }
func (r *testRepo) UpdateLayout(ctx context.Context, layout feat.Layout) error          { return nil }
func (r *testRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error                { return nil }
func (r *testRepo) CreateTag(ctx context.Context, tag feat.Tag) error                   { return nil }
func (r *testRepo) GetTag(ctx context.Context, id uuid.UUID) (feat.Tag, error)          { return feat.Tag{}, nil }
func (r *testRepo) GetTagByName(ctx context.Context, name string) (feat.Tag, error)     { return feat.Tag{}, nil }
func (r *testRepo) GetAllTags(ctx context.Context) ([]feat.Tag, error)                  { return nil, nil }
func (r *testRepo) UpdateTag(ctx context.Context, tag feat.Tag) error                   { return nil }
func (r *testRepo) DeleteTag(ctx context.Context, id uuid.UUID) error                   { return nil }
func (r *testRepo) CreateParam(ctx context.Context, param *feat.Param) error            { return nil }
func (r *testRepo) GetParam(ctx context.Context, id uuid.UUID) (feat.Param, error)      { return feat.Param{}, nil }
func (r *testRepo) GetParamByName(ctx context.Context, name string) (feat.Param, error) {
	return feat.Param{}, nil
}
func (r *testRepo) GetParamByRefKey(ctx context.Context, refKey string) (feat.Param, error) {
	return feat.Param{}, nil
}
func (r *testRepo) ListParams(ctx context.Context) ([]feat.Param, error)            { return nil, nil }
func (r *testRepo) UpdateParam(ctx context.Context, param *feat.Param) error        { return nil }
func (r *testRepo) DeleteParam(ctx context.Context, id uuid.UUID) error             { return nil }
func (r *testRepo) CreateImage(ctx context.Context, image *feat.Image) error        { return nil }
func (r *testRepo) GetImage(ctx context.Context, id uuid.UUID) (feat.Image, error)  { return feat.Image{}, nil }
func (r *testRepo) GetImageByShortID(ctx context.Context, shortID string) (feat.Image, error) {
	return feat.Image{}, nil
}
func (r *testRepo) GetImageByContentHash(ctx context.Context, contentHash string) (feat.Image, error) {
	return feat.Image{}, nil
}
func (r *testRepo) UpdateImage(ctx context.Context, image *feat.Image) error            { return nil }
func (r *testRepo) DeleteImage(ctx context.Context, id uuid.UUID) error                 { return nil }
func (r *testRepo) ListImages(ctx context.Context) ([]feat.Image, error)                { return nil, nil }
func (r *testRepo) CreateImageVariant(ctx context.Context, variant *feat.ImageVariant) error { return nil }
func (r *testRepo) GetImageVariant(ctx context.Context, id uuid.UUID) (feat.ImageVariant, error) {
	return feat.ImageVariant{}, nil
}
func (r *testRepo) UpdateImageVariant(ctx context.Context, variant *feat.ImageVariant) error { return nil }
func (r *testRepo) DeleteImageVariant(ctx context.Context, id uuid.UUID) error { return nil }
func (r *testRepo) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]feat.ImageVariant, error) {
	return nil, nil
}
func (r *testRepo) CreateContentImage(ctx context.Context, contentImage *feat.ContentImage) error {
	return nil
}
func (r *testRepo) DeleteContentImage(ctx context.Context, id uuid.UUID) error { return nil }
func (r *testRepo) GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]feat.ContentImage, error) {
	return nil, nil
}
func (r *testRepo) CreateSectionImage(ctx context.Context, sectionImage *feat.SectionImage) error {
	return nil
}
func (r *testRepo) DeleteSectionImage(ctx context.Context, id uuid.UUID) error { return nil }
func (r *testRepo) GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]feat.SectionImage, error) {
	return nil, nil
}
func (r *testRepo) CreateMeta(ctx context.Context, meta *feat.Meta) error { return nil }
func (r *testRepo) GetMeta(ctx context.Context, id uuid.UUID) (feat.Meta, error) {
	return feat.Meta{}, nil
}
func (r *testRepo) GetAllMeta(ctx context.Context) ([]feat.Meta, error)            { return nil, nil }
func (r *testRepo) GetMetaByContentID(ctx context.Context, contentID uuid.UUID) ([]feat.Meta, error) {
	return nil, nil
}
func (r *testRepo) UpdateMeta(ctx context.Context, meta *feat.Meta) error { return nil }
func (r *testRepo) DeleteMeta(ctx context.Context, id uuid.UUID) error    { return nil }
func (r *testRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error { return nil }
func (r *testRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	return nil
}
func (r *testRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]feat.Tag, error) {
	return nil, nil
}
func (r *testRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]feat.Content, error) {
	return nil, nil
}
func (r *testRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	return auth.User{}, nil
}
