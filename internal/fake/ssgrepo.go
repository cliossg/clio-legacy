package fake

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/auth"
	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
)

type SsgRepo struct {
	hm.Core

	CreateContentFn                      func(ctx context.Context, content *ssg.Content) error
	GetContentFn                         func(ctx context.Context, id uuid.UUID) (ssg.Content, error)
	UpdateContentFn                      func(ctx context.Context, content *ssg.Content) error
	DeleteContentFn                      func(ctx context.Context, id uuid.UUID) error
	GetAllContentWithMetaFn              func(ctx context.Context) ([]ssg.Content, error)
	GetContentWithPaginationAndSearchFn  func(ctx context.Context, offset, limit int, searchQuery string) ([]ssg.Content, int, error)
	CreateSectionFn                      func(ctx context.Context, section ssg.Section) error
	GetSectionFn                         func(ctx context.Context, id uuid.UUID) (ssg.Section, error)
	GetSectionsFn                        func(ctx context.Context) ([]ssg.Section, error)
	UpdateSectionFn                      func(ctx context.Context, section ssg.Section) error
	DeleteSectionFn                      func(ctx context.Context, id uuid.UUID) error
	CreateLayoutFn                       func(ctx context.Context, layout ssg.Layout) error
	GetLayoutFn                          func(ctx context.Context, id uuid.UUID) (ssg.Layout, error)
	GetAllLayoutsFn                      func(ctx context.Context) ([]ssg.Layout, error)
	UpdateLayoutFn                       func(ctx context.Context, layout ssg.Layout) error
	DeleteLayoutFn                       func(ctx context.Context, id uuid.UUID) error
	CreateTagFn                          func(ctx context.Context, tag ssg.Tag) error
	GetTagFn                             func(ctx context.Context, id uuid.UUID) (ssg.Tag, error)
	GetTagByNameFn                       func(ctx context.Context, name string) (ssg.Tag, error)
	GetAllTagsFn                         func(ctx context.Context) ([]ssg.Tag, error)
	UpdateTagFn                          func(ctx context.Context, tag ssg.Tag) error
	DeleteTagFn                          func(ctx context.Context, id uuid.UUID) error
	CreateParamFn                        func(ctx context.Context, param *ssg.Param) error
	GetParamFn                           func(ctx context.Context, id uuid.UUID) (ssg.Param, error)
	GetParamByNameFn                     func(ctx context.Context, name string) (ssg.Param, error)
	GetParamByRefKeyFn                   func(ctx context.Context, refKey string) (ssg.Param, error)
	ListParamsFn                         func(ctx context.Context) ([]ssg.Param, error)
	UpdateParamFn                        func(ctx context.Context, param *ssg.Param) error
	DeleteParamFn                        func(ctx context.Context, id uuid.UUID) error
	CreateImageFn                        func(ctx context.Context, image *ssg.Image) error
	GetImageFn                           func(ctx context.Context, id uuid.UUID) (ssg.Image, error)
	GetImageByShortIDFn                  func(ctx context.Context, shortID string) (ssg.Image, error)
	GetImageByContentHashFn              func(ctx context.Context, contentHash string) (ssg.Image, error)
	UpdateImageFn                        func(ctx context.Context, image *ssg.Image) error
	DeleteImageFn                        func(ctx context.Context, id uuid.UUID) error
	ListImagesFn                         func(ctx context.Context) ([]ssg.Image, error)
	CreateImageVariantFn                 func(ctx context.Context, variant *ssg.ImageVariant) error
	GetImageVariantFn                    func(ctx context.Context, id uuid.UUID) (ssg.ImageVariant, error)
	UpdateImageVariantFn                 func(ctx context.Context, variant *ssg.ImageVariant) error
	DeleteImageVariantFn                 func(ctx context.Context, id uuid.UUID) error
	ListImageVariantsByImageIDFn         func(ctx context.Context, imageID uuid.UUID) ([]ssg.ImageVariant, error)
	CreateContentImageFn                 func(ctx context.Context, contentImage *ssg.ContentImage) error
	DeleteContentImageFn                 func(ctx context.Context, id uuid.UUID) error
	GetContentImagesByContentIDFn        func(ctx context.Context, contentID uuid.UUID) ([]ssg.ContentImage, error)
	CreateSectionImageFn                 func(ctx context.Context, sectionImage *ssg.SectionImage) error
	DeleteSectionImageFn                 func(ctx context.Context, id uuid.UUID) error
	GetSectionImagesBySectionIDFn        func(ctx context.Context, sectionID uuid.UUID) ([]ssg.SectionImage, error)
	AddTagToContentFn                    func(ctx context.Context, contentID, tagID uuid.UUID) error
	RemoveTagFromContentFn               func(ctx context.Context, contentID, tagID uuid.UUID) error
	GetTagsForContentFn                  func(ctx context.Context, contentID uuid.UUID) ([]ssg.Tag, error)
	GetContentForTagFn                   func(ctx context.Context, tagID uuid.UUID) ([]ssg.Content, error)
	GetUserByUsernameFn                  func(ctx context.Context, username string) (auth.User, error)
	GetSiteBySlugFn                      func(ctx context.Context, slug string) (ssg.Site, error)

	contents       map[uuid.UUID]ssg.Content
	sections       map[uuid.UUID]ssg.Section
	layouts        map[uuid.UUID]ssg.Layout
	tags           map[uuid.UUID]ssg.Tag
	tagsByName     map[string]ssg.Tag
	params         map[uuid.UUID]ssg.Param
	paramsByName   map[string]ssg.Param
	paramsByRefKey map[string]ssg.Param
	images         map[uuid.UUID]ssg.Image
	imagesByShort  map[string]ssg.Image
	imagesByPath   map[string]ssg.Image
	imageVariants  map[uuid.UUID]ssg.ImageVariant
	contentImages  map[uuid.UUID][]ssg.ContentImage
	sectionImages  map[uuid.UUID][]ssg.SectionImage
	contentTags    map[uuid.UUID][]uuid.UUID
	users          map[string]auth.User
	sites          map[string]ssg.Site
}

func NewSsgRepo() *SsgRepo {
	cfg := hm.NewConfig()
	return &SsgRepo{
		Core:           hm.NewCore("fake-ssg-repo", hm.XParams{Cfg: cfg}),
		contents:       make(map[uuid.UUID]ssg.Content),
		sections:       make(map[uuid.UUID]ssg.Section),
		layouts:        make(map[uuid.UUID]ssg.Layout),
		tags:           make(map[uuid.UUID]ssg.Tag),
		tagsByName:     make(map[string]ssg.Tag),
		params:         make(map[uuid.UUID]ssg.Param),
		paramsByName:   make(map[string]ssg.Param),
		paramsByRefKey: make(map[string]ssg.Param),
		images:         make(map[uuid.UUID]ssg.Image),
		imagesByShort:  make(map[string]ssg.Image),
		imagesByPath:   make(map[string]ssg.Image),
		imageVariants:  make(map[uuid.UUID]ssg.ImageVariant),
		contentImages:  make(map[uuid.UUID][]ssg.ContentImage),
		sectionImages:  make(map[uuid.UUID][]ssg.SectionImage),
		contentTags:    make(map[uuid.UUID][]uuid.UUID),
		users:          make(map[string]auth.User),
		sites:          make(map[string]ssg.Site),
	}
}

func (f *SsgRepo) Query() *hm.QueryManager {
	return nil
}

func (f *SsgRepo) BeginTx(ctx context.Context) (context.Context, hm.Tx, error) {
	return ctx, nil, nil
}

func (f *SsgRepo) CreateContent(ctx context.Context, content *ssg.Content) error {
	if f.CreateContentFn != nil {
		return f.CreateContentFn(ctx, content)
	}
	f.contents[content.ID] = *content
	return nil
}

func (f *SsgRepo) GetContent(ctx context.Context, id uuid.UUID) (ssg.Content, error) {
	if f.GetContentFn != nil {
		return f.GetContentFn(ctx, id)
	}
	if c, ok := f.contents[id]; ok {
		return c, nil
	}
	return ssg.Content{}, fmt.Errorf("content not found")
}

func (f *SsgRepo) UpdateContent(ctx context.Context, content *ssg.Content) error {
	if f.UpdateContentFn != nil {
		return f.UpdateContentFn(ctx, content)
	}
	f.contents[content.ID] = *content
	return nil
}

func (f *SsgRepo) DeleteContent(ctx context.Context, id uuid.UUID) error {
	if f.DeleteContentFn != nil {
		return f.DeleteContentFn(ctx, id)
	}
	delete(f.contents, id)
	return nil
}

func (f *SsgRepo) GetAllContentWithMeta(ctx context.Context) ([]ssg.Content, error) {
	if f.GetAllContentWithMetaFn != nil {
		return f.GetAllContentWithMetaFn(ctx)
	}
	var contents []ssg.Content
	for _, c := range f.contents {
		contents = append(contents, c)
	}
	return contents, nil
}

func (f *SsgRepo) GetContentWithPaginationAndSearch(ctx context.Context, offset, limit int, searchQuery string) ([]ssg.Content, int, error) {
	if f.GetContentWithPaginationAndSearchFn != nil {
		return f.GetContentWithPaginationAndSearchFn(ctx, offset, limit, searchQuery)
	}
	var contents []ssg.Content
	for _, c := range f.contents {
		contents = append(contents, c)
	}
	return contents, len(contents), nil
}

func (f *SsgRepo) CreateSection(ctx context.Context, section ssg.Section) error {
	if f.CreateSectionFn != nil {
		return f.CreateSectionFn(ctx, section)
	}
	f.sections[section.ID] = section
	return nil
}

func (f *SsgRepo) GetSection(ctx context.Context, id uuid.UUID) (ssg.Section, error) {
	if f.GetSectionFn != nil {
		return f.GetSectionFn(ctx, id)
	}
	if s, ok := f.sections[id]; ok {
		return s, nil
	}
	return ssg.Section{}, fmt.Errorf("section not found")
}

func (f *SsgRepo) GetSections(ctx context.Context) ([]ssg.Section, error) {
	if f.GetSectionsFn != nil {
		return f.GetSectionsFn(ctx)
	}
	var sections []ssg.Section
	for _, s := range f.sections {
		sections = append(sections, s)
	}
	return sections, nil
}

func (f *SsgRepo) UpdateSection(ctx context.Context, section ssg.Section) error {
	if f.UpdateSectionFn != nil {
		return f.UpdateSectionFn(ctx, section)
	}
	f.sections[section.ID] = section
	return nil
}

func (f *SsgRepo) DeleteSection(ctx context.Context, id uuid.UUID) error {
	if f.DeleteSectionFn != nil {
		return f.DeleteSectionFn(ctx, id)
	}
	delete(f.sections, id)
	return nil
}

func (f *SsgRepo) CreateLayout(ctx context.Context, layout ssg.Layout) error {
	if f.CreateLayoutFn != nil {
		return f.CreateLayoutFn(ctx, layout)
	}
	f.layouts[layout.ID] = layout
	return nil
}

func (f *SsgRepo) GetLayout(ctx context.Context, id uuid.UUID) (ssg.Layout, error) {
	if f.GetLayoutFn != nil {
		return f.GetLayoutFn(ctx, id)
	}
	if l, ok := f.layouts[id]; ok {
		return l, nil
	}
	return ssg.Layout{}, fmt.Errorf("layout not found")
}

func (f *SsgRepo) GetAllLayouts(ctx context.Context) ([]ssg.Layout, error) {
	if f.GetAllLayoutsFn != nil {
		return f.GetAllLayoutsFn(ctx)
	}
	var layouts []ssg.Layout
	for _, l := range f.layouts {
		layouts = append(layouts, l)
	}
	return layouts, nil
}

func (f *SsgRepo) UpdateLayout(ctx context.Context, layout ssg.Layout) error {
	if f.UpdateLayoutFn != nil {
		return f.UpdateLayoutFn(ctx, layout)
	}
	f.layouts[layout.ID] = layout
	return nil
}

func (f *SsgRepo) DeleteLayout(ctx context.Context, id uuid.UUID) error {
	if f.DeleteLayoutFn != nil {
		return f.DeleteLayoutFn(ctx, id)
	}
	delete(f.layouts, id)
	return nil
}

func (f *SsgRepo) CreateTag(ctx context.Context, tag ssg.Tag) error {
	if f.CreateTagFn != nil {
		return f.CreateTagFn(ctx, tag)
	}
	f.tags[tag.ID] = tag
	f.tagsByName[tag.Name] = tag
	return nil
}

func (f *SsgRepo) GetTag(ctx context.Context, id uuid.UUID) (ssg.Tag, error) {
	if f.GetTagFn != nil {
		return f.GetTagFn(ctx, id)
	}
	if t, ok := f.tags[id]; ok {
		return t, nil
	}
	return ssg.Tag{}, fmt.Errorf("tag not found")
}

func (f *SsgRepo) GetTagByName(ctx context.Context, name string) (ssg.Tag, error) {
	if f.GetTagByNameFn != nil {
		return f.GetTagByNameFn(ctx, name)
	}
	if t, ok := f.tagsByName[name]; ok {
		return t, nil
	}
	return ssg.Tag{}, fmt.Errorf("tag not found")
}

func (f *SsgRepo) GetAllTags(ctx context.Context) ([]ssg.Tag, error) {
	if f.GetAllTagsFn != nil {
		return f.GetAllTagsFn(ctx)
	}
	var tags []ssg.Tag
	for _, t := range f.tags {
		tags = append(tags, t)
	}
	return tags, nil
}

func (f *SsgRepo) UpdateTag(ctx context.Context, tag ssg.Tag) error {
	if f.UpdateTagFn != nil {
		return f.UpdateTagFn(ctx, tag)
	}
	f.tags[tag.ID] = tag
	f.tagsByName[tag.Name] = tag
	return nil
}

func (f *SsgRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	if f.DeleteTagFn != nil {
		return f.DeleteTagFn(ctx, id)
	}
	if tag, ok := f.tags[id]; ok {
		delete(f.tagsByName, tag.Name)
	}
	delete(f.tags, id)
	return nil
}

func (f *SsgRepo) CreateParam(ctx context.Context, param *ssg.Param) error {
	if f.CreateParamFn != nil {
		return f.CreateParamFn(ctx, param)
	}
	f.params[param.ID] = *param
	f.paramsByName[param.Name] = *param
	f.paramsByRefKey[param.RefKey] = *param
	return nil
}

func (f *SsgRepo) GetParam(ctx context.Context, id uuid.UUID) (ssg.Param, error) {
	if f.GetParamFn != nil {
		return f.GetParamFn(ctx, id)
	}
	if p, ok := f.params[id]; ok {
		return p, nil
	}
	return ssg.Param{}, fmt.Errorf("param not found")
}

func (f *SsgRepo) GetParamByName(ctx context.Context, name string) (ssg.Param, error) {
	if f.GetParamByNameFn != nil {
		return f.GetParamByNameFn(ctx, name)
	}
	if p, ok := f.paramsByName[name]; ok {
		return p, nil
	}
	return ssg.Param{}, fmt.Errorf("param not found")
}

func (f *SsgRepo) GetParamByRefKey(ctx context.Context, refKey string) (ssg.Param, error) {
	if f.GetParamByRefKeyFn != nil {
		return f.GetParamByRefKeyFn(ctx, refKey)
	}
	if p, ok := f.paramsByRefKey[refKey]; ok {
		return p, nil
	}
	return ssg.Param{}, fmt.Errorf("param not found")
}

func (f *SsgRepo) ListParams(ctx context.Context) ([]ssg.Param, error) {
	if f.ListParamsFn != nil {
		return f.ListParamsFn(ctx)
	}
	var params []ssg.Param
	for _, p := range f.params {
		params = append(params, p)
	}
	return params, nil
}

func (f *SsgRepo) UpdateParam(ctx context.Context, param *ssg.Param) error {
	if f.UpdateParamFn != nil {
		return f.UpdateParamFn(ctx, param)
	}
	f.params[param.ID] = *param
	f.paramsByName[param.Name] = *param
	f.paramsByRefKey[param.RefKey] = *param
	return nil
}

func (f *SsgRepo) DeleteParam(ctx context.Context, id uuid.UUID) error {
	if f.DeleteParamFn != nil {
		return f.DeleteParamFn(ctx, id)
	}
	if param, ok := f.params[id]; ok {
		delete(f.paramsByName, param.Name)
		delete(f.paramsByRefKey, param.RefKey)
	}
	delete(f.params, id)
	return nil
}

func (f *SsgRepo) CreateImage(ctx context.Context, image *ssg.Image) error {
	if f.CreateImageFn != nil {
		return f.CreateImageFn(ctx, image)
	}
	f.images[image.ID] = *image
	f.imagesByShort[image.ShortID] = *image
	f.imagesByPath[image.FilePath] = *image
	return nil
}

func (f *SsgRepo) GetImage(ctx context.Context, id uuid.UUID) (ssg.Image, error) {
	if f.GetImageFn != nil {
		return f.GetImageFn(ctx, id)
	}
	if i, ok := f.images[id]; ok {
		return i, nil
	}
	return ssg.Image{}, fmt.Errorf("image not found")
}

func (f *SsgRepo) GetImageByShortID(ctx context.Context, shortID string) (ssg.Image, error) {
	if f.GetImageByShortIDFn != nil {
		return f.GetImageByShortIDFn(ctx, shortID)
	}
	if i, ok := f.imagesByShort[shortID]; ok {
		return i, nil
	}
	return ssg.Image{}, fmt.Errorf("image not found")
}

func (f *SsgRepo) GetImageByContentHash(ctx context.Context, contentHash string) (ssg.Image, error) {
	if f.GetImageByContentHashFn != nil {
		return f.GetImageByContentHashFn(ctx, contentHash)
	}
	if i, ok := f.imagesByPath[contentHash]; ok {
		return i, nil
	}
	return ssg.Image{}, fmt.Errorf("image not found")
}

func (f *SsgRepo) UpdateImage(ctx context.Context, image *ssg.Image) error {
	if f.UpdateImageFn != nil {
		return f.UpdateImageFn(ctx, image)
	}
	f.images[image.ID] = *image
	f.imagesByShort[image.ShortID] = *image
	f.imagesByPath[image.FilePath] = *image
	return nil
}

func (f *SsgRepo) DeleteImage(ctx context.Context, id uuid.UUID) error {
	if f.DeleteImageFn != nil {
		return f.DeleteImageFn(ctx, id)
	}
	if image, ok := f.images[id]; ok {
		delete(f.imagesByShort, image.ShortID)
		delete(f.imagesByPath, image.FilePath)
	}
	delete(f.images, id)
	return nil
}

func (f *SsgRepo) ListImages(ctx context.Context) ([]ssg.Image, error) {
	if f.ListImagesFn != nil {
		return f.ListImagesFn(ctx)
	}
	var images []ssg.Image
	for _, i := range f.images {
		images = append(images, i)
	}
	return images, nil
}

func (f *SsgRepo) CreateImageVariant(ctx context.Context, variant *ssg.ImageVariant) error {
	if f.CreateImageVariantFn != nil {
		return f.CreateImageVariantFn(ctx, variant)
	}
	f.imageVariants[variant.ID] = *variant
	return nil
}

func (f *SsgRepo) GetImageVariant(ctx context.Context, id uuid.UUID) (ssg.ImageVariant, error) {
	if f.GetImageVariantFn != nil {
		return f.GetImageVariantFn(ctx, id)
	}
	if v, ok := f.imageVariants[id]; ok {
		return v, nil
	}
	return ssg.ImageVariant{}, fmt.Errorf("image variant not found")
}

func (f *SsgRepo) UpdateImageVariant(ctx context.Context, variant *ssg.ImageVariant) error {
	if f.UpdateImageVariantFn != nil {
		return f.UpdateImageVariantFn(ctx, variant)
	}
	f.imageVariants[variant.ID] = *variant
	return nil
}

func (f *SsgRepo) DeleteImageVariant(ctx context.Context, id uuid.UUID) error {
	if f.DeleteImageVariantFn != nil {
		return f.DeleteImageVariantFn(ctx, id)
	}
	delete(f.imageVariants, id)
	return nil
}

func (f *SsgRepo) ListImageVariantsByImageID(ctx context.Context, imageID uuid.UUID) ([]ssg.ImageVariant, error) {
	if f.ListImageVariantsByImageIDFn != nil {
		return f.ListImageVariantsByImageIDFn(ctx, imageID)
	}
	var variants []ssg.ImageVariant
	for _, v := range f.imageVariants {
		if v.ImageID == imageID {
			variants = append(variants, v)
		}
	}
	return variants, nil
}

func (f *SsgRepo) CreateContentImage(ctx context.Context, contentImage *ssg.ContentImage) error {
	if f.CreateContentImageFn != nil {
		return f.CreateContentImageFn(ctx, contentImage)
	}
	f.contentImages[contentImage.ContentID] = append(f.contentImages[contentImage.ContentID], *contentImage)
	return nil
}

func (f *SsgRepo) DeleteContentImage(ctx context.Context, id uuid.UUID) error {
	if f.DeleteContentImageFn != nil {
		return f.DeleteContentImageFn(ctx, id)
	}
	for contentID, images := range f.contentImages {
		for i, img := range images {
			if img.ID == id {
				f.contentImages[contentID] = append(images[:i], images[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (f *SsgRepo) GetContentImagesByContentID(ctx context.Context, contentID uuid.UUID) ([]ssg.ContentImage, error) {
	if f.GetContentImagesByContentIDFn != nil {
		return f.GetContentImagesByContentIDFn(ctx, contentID)
	}
	return f.contentImages[contentID], nil
}

func (f *SsgRepo) CreateSectionImage(ctx context.Context, sectionImage *ssg.SectionImage) error {
	if f.CreateSectionImageFn != nil {
		return f.CreateSectionImageFn(ctx, sectionImage)
	}
	f.sectionImages[sectionImage.SectionID] = append(f.sectionImages[sectionImage.SectionID], *sectionImage)
	return nil
}

func (f *SsgRepo) DeleteSectionImage(ctx context.Context, id uuid.UUID) error {
	if f.DeleteSectionImageFn != nil {
		return f.DeleteSectionImageFn(ctx, id)
	}
	for sectionID, images := range f.sectionImages {
		for i, img := range images {
			if img.ID == id {
				f.sectionImages[sectionID] = append(images[:i], images[i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (f *SsgRepo) GetSectionImagesBySectionID(ctx context.Context, sectionID uuid.UUID) ([]ssg.SectionImage, error) {
	if f.GetSectionImagesBySectionIDFn != nil {
		return f.GetSectionImagesBySectionIDFn(ctx, sectionID)
	}
	return f.sectionImages[sectionID], nil
}

func (f *SsgRepo) AddTagToContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	if f.AddTagToContentFn != nil {
		return f.AddTagToContentFn(ctx, contentID, tagID)
	}
	f.contentTags[contentID] = append(f.contentTags[contentID], tagID)
	return nil
}

func (f *SsgRepo) RemoveTagFromContent(ctx context.Context, contentID, tagID uuid.UUID) error {
	if f.RemoveTagFromContentFn != nil {
		return f.RemoveTagFromContentFn(ctx, contentID, tagID)
	}
	tags := f.contentTags[contentID]
	for i, tid := range tags {
		if tid == tagID {
			f.contentTags[contentID] = append(tags[:i], tags[i+1:]...)
			return nil
		}
	}
	return nil
}

func (f *SsgRepo) GetTagsForContent(ctx context.Context, contentID uuid.UUID) ([]ssg.Tag, error) {
	if f.GetTagsForContentFn != nil {
		return f.GetTagsForContentFn(ctx, contentID)
	}
	var tags []ssg.Tag
	for _, tagID := range f.contentTags[contentID] {
		if tag, ok := f.tags[tagID]; ok {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func (f *SsgRepo) GetContentForTag(ctx context.Context, tagID uuid.UUID) ([]ssg.Content, error) {
	if f.GetContentForTagFn != nil {
		return f.GetContentForTagFn(ctx, tagID)
	}
	var contents []ssg.Content
	for contentID, tags := range f.contentTags {
		for _, tid := range tags {
			if tid == tagID {
				if content, ok := f.contents[contentID]; ok {
					contents = append(contents, content)
				}
			}
		}
	}
	return contents, nil
}

func (f *SsgRepo) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	if f.GetUserByUsernameFn != nil {
		return f.GetUserByUsernameFn(ctx, username)
	}
	if u, ok := f.users[username]; ok {
		return u, nil
	}
	return auth.User{}, fmt.Errorf("user not found")
}

func (f *SsgRepo) GetSiteBySlug(ctx context.Context, slug string) (ssg.Site, error) {
	if f.GetSiteBySlugFn != nil {
		return f.GetSiteBySlugFn(ctx, slug)
	}
	if s, ok := f.sites[slug]; ok {
		return s, nil
	}
	return ssg.Site{}, fmt.Errorf("site not found")
}
