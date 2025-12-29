package ssg

import (
	"context"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewTestRepo(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)

	if repo == nil {
		t.Fatal("newTestRepo() returned nil")
	}
	if repo.siteRepo != mockSite {
		t.Error("newTestRepo() did not set siteRepo correctly")
	}
}

func TestTestRepoGetDB(t *testing.T) {
	repo := newTestRepo(nil)
	if repo.GetDB() != nil {
		t.Error("GetDB() should return nil for repo without DB")
	}
}

func TestTestRepoQuery(t *testing.T) {
	repo := newTestRepo(nil)
	if repo.Query() != nil {
		t.Error("Query() should return nil")
	}
}

func TestTestRepoBeginTx(t *testing.T) {
	repo := newTestRepo(nil)
	ctx := context.Background()

	gotCtx, tx, err := repo.BeginTx(ctx)
	if err != nil {
		t.Errorf("BeginTx() error = %v, want nil", err)
	}
	if tx != nil {
		t.Error("BeginTx() should return nil transaction")
	}
	if gotCtx != ctx {
		t.Error("BeginTx() should return same context")
	}
}

func TestTestRepoCreateSite(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site := &feat.Site{}
	site.ID = uuid.New()

	err := repo.CreateSite(ctx, site)
	if err != nil {
		t.Errorf("CreateSite() error = %v, want nil", err)
	}

	created, err := mockSite.GetSite(ctx, site.ID)
	if err != nil {
		t.Errorf("Site was not created in mock repo")
	}
	if created.ID != site.ID {
		t.Errorf("Created site ID = %v, want %v", created.ID, site.ID)
	}
}

func TestTestRepoGetSite(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site := &feat.Site{}
	site.ID = uuid.New()
	mockSite.CreateSite(ctx, site)

	got, err := repo.GetSite(ctx, site.ID)
	if err != nil {
		t.Errorf("GetSite() error = %v, want nil", err)
	}
	if got.ID != site.ID {
		t.Errorf("GetSite() ID = %v, want %v", got.ID, site.ID)
	}
}

func TestTestRepoGetSiteBySlug(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site := feat.NewSite("Test", "test-slug", "structured")
	site.ID = uuid.New()
	mockSite.CreateSite(ctx, &site)

	got, err := repo.GetSiteBySlug(ctx, "test-slug")
	if err != nil {
		t.Errorf("GetSiteBySlug() error = %v, want nil", err)
	}
	if got.Slug() != "test-slug" {
		t.Errorf("GetSiteBySlug() slug = %v, want test-slug", got.Slug())
	}
}

func TestTestRepoListSites(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site1 := feat.NewSite("Test1", "test1", "structured")
	site1.ID = uuid.New()
	site2 := feat.NewSite("Test2", "test2", "blog")
	site2.ID = uuid.New()

	mockSite.CreateSite(ctx, &site1)
	mockSite.CreateSite(ctx, &site2)

	sites, err := repo.ListSites(ctx, false)
	if err != nil {
		t.Errorf("ListSites() error = %v, want nil", err)
	}
	if len(sites) != 2 {
		t.Errorf("ListSites() count = %v, want 2", len(sites))
	}
}

func TestTestRepoUpdateSite(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site := feat.NewSite("Original", "test", "structured")
	site.ID = uuid.New()
	mockSite.CreateSite(ctx, &site)

	site.Name = "Updated"
	err := repo.UpdateSite(ctx, &site)
	if err != nil {
		t.Errorf("UpdateSite() error = %v, want nil", err)
	}

	updated, _ := mockSite.GetSite(ctx, site.ID)
	if updated.Name != "Updated" {
		t.Errorf("UpdateSite() name = %v, want Updated", updated.Name)
	}
}

func TestTestRepoDeleteSite(t *testing.T) {
	mockSite := newMockSiteRepo()
	repo := newTestRepo(mockSite)
	ctx := context.Background()

	site := feat.NewSite("Test", "test", "structured")
	site.ID = uuid.New()
	mockSite.CreateSite(ctx, &site)

	err := repo.DeleteSite(ctx, site.ID)
	if err != nil {
		t.Errorf("DeleteSite() error = %v, want nil", err)
	}

	_, err = mockSite.GetSite(ctx, site.ID)
	if err == nil {
		t.Error("Site should have been deleted")
	}
}

func TestTestRepoStubMethods(t *testing.T) {
	repo := newTestRepo(nil)
	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "CreateContent",
			fn: func() error {
				return repo.CreateContent(ctx, &feat.Content{})
			},
		},
		{
			name: "UpdateContent",
			fn: func() error {
				return repo.UpdateContent(ctx, &feat.Content{})
			},
		},
		{
			name: "DeleteContent",
			fn: func() error {
				return repo.DeleteContent(ctx, uuid.New())
			},
		},
		{
			name: "CreateSection",
			fn: func() error {
				return repo.CreateSection(ctx, feat.Section{})
			},
		},
		{
			name: "UpdateSection",
			fn: func() error {
				return repo.UpdateSection(ctx, feat.Section{})
			},
		},
		{
			name: "DeleteSection",
			fn: func() error {
				return repo.DeleteSection(ctx, uuid.New())
			},
		},
		{
			name: "CreateLayout",
			fn: func() error {
				return repo.CreateLayout(ctx, feat.Layout{})
			},
		},
		{
			name: "UpdateLayout",
			fn: func() error {
				return repo.UpdateLayout(ctx, feat.Layout{})
			},
		},
		{
			name: "DeleteLayout",
			fn: func() error {
				return repo.DeleteLayout(ctx, uuid.New())
			},
		},
		{
			name: "CreateTag",
			fn: func() error {
				return repo.CreateTag(ctx, feat.Tag{})
			},
		},
		{
			name: "UpdateTag",
			fn: func() error {
				return repo.UpdateTag(ctx, feat.Tag{})
			},
		},
		{
			name: "DeleteTag",
			fn: func() error {
				return repo.DeleteTag(ctx, uuid.New())
			},
		},
		{
			name: "CreateParam",
			fn: func() error {
				return repo.CreateParam(ctx, &feat.Param{})
			},
		},
		{
			name: "UpdateParam",
			fn: func() error {
				return repo.UpdateParam(ctx, &feat.Param{})
			},
		},
		{
			name: "DeleteParam",
			fn: func() error {
				return repo.DeleteParam(ctx, uuid.New())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s() error = %v, want nil", tt.name, err)
			}
		})
	}
}

func TestTestRepoGetMethods(t *testing.T) {
	repo := newTestRepo(nil)
	ctx := context.Background()

	t.Run("GetContent", func(t *testing.T) {
		content, err := repo.GetContent(ctx, uuid.New())
		if err != nil {
			t.Errorf("GetContent() error = %v, want nil", err)
		}
		if !content.IsZero() {
			t.Error("GetContent() should return zero value")
		}
	})

	t.Run("GetAllContentWithMeta", func(t *testing.T) {
		contents, err := repo.GetAllContentWithMeta(ctx)
		if err != nil {
			t.Errorf("GetAllContentWithMeta() error = %v, want nil", err)
		}
		if contents != nil {
			t.Error("GetAllContentWithMeta() should return nil")
		}
	})

	t.Run("GetContentWithPaginationAndSearch", func(t *testing.T) {
		contents, count, err := repo.GetContentWithPaginationAndSearch(ctx, 0, 10, "")
		if err != nil {
			t.Errorf("GetContentWithPaginationAndSearch() error = %v, want nil", err)
		}
		if contents != nil {
			t.Error("GetContentWithPaginationAndSearch() contents should be nil")
		}
		if count != 0 {
			t.Errorf("GetContentWithPaginationAndSearch() count = %v, want 0", count)
		}
	})

	t.Run("GetSection", func(t *testing.T) {
		section, err := repo.GetSection(ctx, uuid.New())
		if err != nil {
			t.Errorf("GetSection() error = %v, want nil", err)
		}
		if !section.IsZero() {
			t.Error("GetSection() should return zero value")
		}
	})

	t.Run("GetSections", func(t *testing.T) {
		sections, err := repo.GetSections(ctx)
		if err != nil {
			t.Errorf("GetSections() error = %v, want nil", err)
		}
		if sections != nil {
			t.Error("GetSections() should return nil")
		}
	})

	t.Run("GetLayout", func(t *testing.T) {
		layout, err := repo.GetLayout(ctx, uuid.New())
		if err != nil {
			t.Errorf("GetLayout() error = %v, want nil", err)
		}
		if !layout.IsZero() {
			t.Error("GetLayout() should return zero value")
		}
	})

	t.Run("GetAllLayouts", func(t *testing.T) {
		layouts, err := repo.GetAllLayouts(ctx)
		if err != nil {
			t.Errorf("GetAllLayouts() error = %v, want nil", err)
		}
		if layouts != nil {
			t.Error("GetAllLayouts() should return nil")
		}
	})

	t.Run("GetTag", func(t *testing.T) {
		tag, err := repo.GetTag(ctx, uuid.New())
		if err != nil {
			t.Errorf("GetTag() error = %v, want nil", err)
		}
		if !tag.IsZero() {
			t.Error("GetTag() should return zero value")
		}
	})

	t.Run("GetTagByName", func(t *testing.T) {
		tag, err := repo.GetTagByName(ctx, "test")
		if err != nil {
			t.Errorf("GetTagByName() error = %v, want nil", err)
		}
		if !tag.IsZero() {
			t.Error("GetTagByName() should return zero value")
		}
	})

	t.Run("GetAllTags", func(t *testing.T) {
		tags, err := repo.GetAllTags(ctx)
		if err != nil {
			t.Errorf("GetAllTags() error = %v, want nil", err)
		}
		if tags != nil {
			t.Error("GetAllTags() should return nil")
		}
	})

	t.Run("GetParam", func(t *testing.T) {
		param, err := repo.GetParam(ctx, uuid.New())
		if err != nil {
			t.Errorf("GetParam() error = %v, want nil", err)
		}
		if !param.IsZero() {
			t.Error("GetParam() should return zero value")
		}
	})
}
