package sqlite

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestSsgRepo(t *testing.T) (*ClioRepo, uuid.UUID) {
	tmpDB := t.TempDir() + "/test.db"
	db, err := sqlx.Open("sqlite3", tmpDB)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}

	schema := `
		PRAGMA foreign_keys = OFF;

		CREATE TABLE site (
			id TEXT PRIMARY KEY,
			short_id TEXT NOT NULL DEFAULT '',
			name TEXT NOT NULL,
			slug TEXT NOT NULL UNIQUE,
			mode TEXT NOT NULL DEFAULT 'structured',
			active INTEGER NOT NULL DEFAULT 1,
			created_by TEXT NOT NULL,
			updated_by TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS content (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			user_id TEXT,
			short_id TEXT,
			section_id TEXT,
			kind TEXT,
			heading TEXT NOT NULL,
			summary TEXT,
			body TEXT,
			draft INTEGER DEFAULT 0,
			featured INTEGER DEFAULT 0,
			series TEXT,
			series_order INTEGER,
			published_at TIMESTAMP,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS section (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			path TEXT,
			layout_id TEXT,
			layout_name TEXT,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tag (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			name TEXT NOT NULL,
			slug TEXT NOT NULL,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS layout (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			code TEXT,
			header_image_id TEXT,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS meta (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			content_id TEXT NOT NULL,
			summary TEXT,
			excerpt TEXT,
			description TEXT,
			keywords TEXT,
			robots TEXT,
			canonical_url TEXT,
			sitemap TEXT,
			table_of_contents INTEGER DEFAULT 0,
			share INTEGER DEFAULT 0,
			comments INTEGER DEFAULT 0,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS param (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			value TEXT,
			ref_key TEXT,
			system INTEGER DEFAULT 0,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS content_tag (
			id TEXT PRIMARY KEY,
			content_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS image (
			id TEXT PRIMARY KEY,
			site_id TEXT NOT NULL,
			short_id TEXT,
			file_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			alt_text TEXT,
			title TEXT,
			width INTEGER,
			height INTEGER,
			created_by TEXT,
			updated_by TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS image_variant (
			id TEXT PRIMARY KEY,
			short_id TEXT NOT NULL DEFAULT '',
			image_id TEXT NOT NULL,
			kind TEXT NOT NULL,
			blob_ref TEXT NOT NULL,
			width INTEGER,
			height INTEGER,
			filesize_bytes INTEGER,
			mime TEXT,
			created_by TEXT NOT NULL DEFAULT '',
			updated_by TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS content_images (
			id TEXT PRIMARY KEY,
			content_id TEXT NOT NULL,
			image_id TEXT NOT NULL,
			is_header INTEGER DEFAULT 0,
			is_featured INTEGER DEFAULT 0,
			order_num INTEGER DEFAULT 0,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS section_images (
			id TEXT PRIMARY KEY,
			section_id TEXT NOT NULL,
			image_id TEXT NOT NULL,
			is_header INTEGER DEFAULT 0,
			is_featured INTEGER DEFAULT 0,
			order_num INTEGER DEFAULT 0,
			created_at TIMESTAMP
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		t.Fatalf("create tables: %v", err)
	}

	siteID := uuid.New()
	_, err = db.Exec(`INSERT INTO site (id, short_id, name, slug, created_by, updated_by, created_at, updated_at)
		VALUES (?, 'test', 'Test Site', 'test-site', '00000000-0000-0000-0000-000000000000', '00000000-0000-0000-0000-000000000000', datetime('now'), datetime('now'))`,
		siteID.String())
	if err != nil {
		t.Fatalf("insert test site: %v", err)
	}

	cfg := hm.NewConfig()
	log := hm.NewLogger("debug")
	params := hm.XParams{Cfg: cfg, Log: log}

	qm := hm.NewQueryManager(testAssetsFS, "sqlite", params)
	ctx := context.Background()
	err = qm.Setup(ctx)
	if err != nil {
		t.Fatalf("setup query manager: %v", err)
	}

	repo := NewClioRepo(qm, params)
	repo.SetDB(db)

	return repo, siteID
}


func TestClioRepoCreateContent(t *testing.T) {
	tests := []struct {
		name    string
		content *ssg.Content
		wantErr bool
	}{
		{
			name: "creates content successfully",
			content: &ssg.Content{
				ID:      uuid.New(),
				Heading: "Test Content",
			},
			wantErr: false,
		},
	}

	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.content.SiteID = siteID
			err := repo.CreateContent(ctx, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetContent(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Test Content",
	}
	repo.CreateContent(ctx, content)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing content",
			id:      content.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent content",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetContent(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetContent() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoUpdateContent(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Original Heading",
	}
	repo.CreateContent(ctx, content)

	tests := []struct {
		name    string
		content *ssg.Content
		wantErr bool
	}{
		{
			name: "updates content successfully",
			content: &ssg.Content{
				ID:      content.ID,
				SiteID:  content.SiteID,
				Heading: "Updated Heading",
			},
			wantErr: false,
		},
		{
			name: "updates content with meta successfully",
			content: &ssg.Content{
				ID:      content.ID,
				SiteID:  content.SiteID,
				Heading: "Updated with Meta",
				Meta: ssg.Meta{
					ID:        uuid.New(),
					ContentID: content.ID,
					Summary:   "Test summary",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateContent(ctx, tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetContent(ctx, tt.content.ID)
				if updated.Heading != tt.content.Heading {
					t.Errorf("UpdateContent() heading = %v, want %v", updated.Heading, tt.content.Heading)
				}
			}
		})
	}
}

func TestClioRepoDeleteContent(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Delete Test",
	}
	repo.CreateContent(ctx, content)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes content successfully",
			id:      content.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteContent(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteContent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetContent(ctx, tt.id)
				if err == nil {
					t.Error("DeleteContent() content still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoCreateSection(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tests := []struct {
		name    string
		section ssg.Section
		wantErr bool
	}{
		{
			name: "creates section successfully",
			section: ssg.Section{
				ID:     uuid.New(),
				SiteID: siteID,
				Name:   "Test Section",
				Path:   "test-section",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateSection(ctx, tt.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetSection(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Test Section",
		Path:   "test-section",
	}
	repo.CreateSection(ctx, section)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing section",
			id:      section.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent section",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetSection(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetSection() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoGetSections(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple sections successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				section1 := ssg.Section{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "Section 1",
					Path:   "section-1",
				}
				section2 := ssg.Section{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "Section 2",
					Path:   "section-2",
				}
				r.CreateSection(ctx, section1)
				r.CreateSection(ctx, section2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no sections",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			sections, err := repo.GetSections(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(sections) != tt.wantCount {
				t.Errorf("GetSections() got %d sections, want %d", len(sections), tt.wantCount)
			}
		})
	}
}

func TestClioRepoCreateTag(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tests := []struct {
		name    string
		tag     ssg.Tag
		wantErr bool
	}{
		{
			name: "creates tag successfully",
			tag: ssg.Tag{
				ID:        uuid.New(),
				SiteID:    siteID,
				Name:      "golang",
				SlugField: "golang",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateTag(ctx, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetTagByName(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)

	tests := []struct {
		name     string
		tagName  string
		wantErr  bool
		wantName string
	}{
		{
			name:     "gets tag by name",
			tagName:  "golang",
			wantErr:  false,
			wantName: "golang",
		},
		{
			name:    "fails with non-existent tag",
			tagName: "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetTagByName(ctx, tt.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Name != tt.wantName {
				t.Errorf("GetTagByName() got name = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
}

func TestClioRepoCreateParam(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tests := []struct {
		name    string
		param   *ssg.Param
		wantErr bool
	}{
		{
			name: "creates param successfully",
			param: &ssg.Param{
				ID:     uuid.New(),
				SiteID: siteID,
				Name:   "site.title",
				RefKey: "ssg.site.title",
				Value:  "My Site",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateParam(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateParam() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetParamByName(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	param := &ssg.Param{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "site.title",
		RefKey: "ssg.site.title",
		Value:  "My Site",
	}
	repo.CreateParam(ctx, param)

	tests := []struct {
		name      string
		paramName string
		wantErr   bool
		wantValue string
	}{
		{
			name:      "gets param by name",
			paramName: "site.title",
			wantErr:   false,
			wantValue: "My Site",
		},
		{
			name:      "fails with non-existent param",
			paramName: "nonexistent",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetParamByName(ctx, tt.paramName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParamByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value != tt.wantValue {
				t.Errorf("GetParamByName() got value = %v, want %v", got.Value, tt.wantValue)
			}
		})
	}
}

func TestClioRepoGetParamByRefKey(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	param := &ssg.Param{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "site.title",
		RefKey: "ssg.site.title",
		Value:  "My Site",
	}
	repo.CreateParam(ctx, param)

	tests := []struct {
		name      string
		refKey    string
		wantErr   bool
		wantValue string
	}{
		{
			name:      "gets param by refkey",
			refKey:    "ssg.site.title",
			wantErr:   false,
			wantValue: "My Site",
		},
		{
			name:    "fails with non-existent refkey",
			refKey:  "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetParamByRefKey(ctx, tt.refKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParamByRefKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Value != tt.wantValue {
				t.Errorf("GetParamByRefKey() got value = %v, want %v", got.Value, tt.wantValue)
			}
		})
	}
}

func TestClioRepoGetParam(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	param := &ssg.Param{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "site.title",
		RefKey: "ssg.site.title",
		Value:  "My Site",
	}
	repo.CreateParam(ctx, param)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing param",
			id:      param.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent param",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetParam(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetParam() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoListParams(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple params successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				param1 := &ssg.Param{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "site.title",
					RefKey: "ssg.site.title",
					Value:  "My Site",
				}
				param2 := &ssg.Param{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "site.author",
					RefKey: "ssg.site.author",
					Value:  "John Doe",
				}
				r.CreateParam(ctx, param1)
				r.CreateParam(ctx, param2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no params",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			params, err := repo.ListParams(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(params) != tt.wantCount {
				t.Errorf("ListParams() got %d params, want %d", len(params), tt.wantCount)
			}
		})
	}
}

func TestClioRepoUpdateParam(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	param := &ssg.Param{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "site.title",
		RefKey: "ssg.site.title",
		Value:  "Original Title",
	}
	repo.CreateParam(ctx, param)

	tests := []struct {
		name    string
		param   *ssg.Param
		wantErr bool
	}{
		{
			name: "updates param successfully",
			param: &ssg.Param{
				ID:     param.ID,
				SiteID: param.SiteID,
				Name:   param.Name,
				RefKey: param.RefKey,
				Value:  "Updated Title",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateParam(ctx, tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateParam() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetParam(ctx, tt.param.ID)
				if updated.Value != tt.param.Value {
					t.Errorf("UpdateParam() value = %v, want %v", updated.Value, tt.param.Value)
				}
			}
		})
	}
}

func TestClioRepoDeleteParam(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	param := &ssg.Param{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "site.title",
		RefKey: "ssg.site.title",
		Value:  "My Site",
	}
	repo.CreateParam(ctx, param)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes param successfully",
			id:      param.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteParam(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteParam() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetParam(ctx, tt.id)
				if err == nil {
					t.Error("DeleteParam() param still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoUpdateSection(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Original Section",
		Path:   "original-section",
	}
	repo.CreateSection(ctx, section)

	tests := []struct {
		name    string
		section ssg.Section
		wantErr bool
	}{
		{
			name: "updates section successfully",
			section: ssg.Section{
				ID:     section.ID,
				SiteID: section.SiteID,
				Name:   "Updated Section",
				Path:   "updated-section",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateSection(ctx, tt.section)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSection() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetSection(ctx, tt.section.ID)
				if updated.Name != tt.section.Name {
					t.Errorf("UpdateSection() name = %v, want %v", updated.Name, tt.section.Name)
				}
			}
		})
	}
}

func TestClioRepoDeleteSection(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Test Section",
		Path:   "test-section",
	}
	repo.CreateSection(ctx, section)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes section successfully",
			id:      section.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteSection(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteSection() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetSection(ctx, tt.id)
				if err == nil {
					t.Error("DeleteSection() section still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoGetTag(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing tag",
			id:      tag.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent tag",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetTag(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetTag() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoGetAllTags(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple tags successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				tag1 := ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "golang",
					SlugField: "golang",
				}
				tag2 := ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "rust",
					SlugField: "rust",
				}
				r.CreateTag(ctx, tag1)
				r.CreateTag(ctx, tag2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no tags",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			tags, err := repo.GetAllTags(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tags) != tt.wantCount {
				t.Errorf("GetAllTags() got %d tags, want %d", len(tags), tt.wantCount)
			}
		})
	}
}

func TestClioRepoUpdateTag(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)

	tests := []struct {
		name    string
		tag     ssg.Tag
		wantErr bool
	}{
		{
			name: "updates tag successfully",
			tag: ssg.Tag{
				ID:        tag.ID,
				SiteID:    tag.SiteID,
				Name:      "go-lang",
				SlugField: "go-lang",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateTag(ctx, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetTag(ctx, tt.tag.ID)
				if updated.Name != tt.tag.Name {
					t.Errorf("UpdateTag() name = %v, want %v", updated.Name, tt.tag.Name)
				}
			}
		})
	}
}

func TestClioRepoDeleteTag(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes tag successfully",
			id:      tag.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteTag(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTag() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetTag(ctx, tt.id)
				if err == nil {
					t.Error("DeleteTag() tag still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoCreateLayout(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tests := []struct {
		name    string
		layout  ssg.Layout
		wantErr bool
	}{
		{
			name: "creates layout successfully",
			layout: ssg.Layout{
				ID:     uuid.New(),
				SiteID: siteID,
				Name:   "Default Layout",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateLayout(ctx, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetAllLayouts(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple layouts successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				layout1 := ssg.Layout{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "Layout 1",
				}
				layout2 := ssg.Layout{
					ID:     uuid.New(),
					SiteID: siteID,
					Name:   "Layout 2",
				}
				r.CreateLayout(ctx, layout1)
				r.CreateLayout(ctx, layout2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no layouts",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			layouts, err := repo.GetAllLayouts(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllLayouts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(layouts) != tt.wantCount {
				t.Errorf("GetAllLayouts() got %d layouts, want %d", len(layouts), tt.wantCount)
			}
		})
	}
}

func TestClioRepoGetLayout(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	layout := ssg.Layout{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Test Layout",
	}
	repo.CreateLayout(ctx, layout)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing layout",
			id:      layout.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent layout",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetLayout(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLayout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetLayout() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoUpdateLayout(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	layout := ssg.Layout{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Original Layout",
	}
	repo.CreateLayout(ctx, layout)

	tests := []struct {
		name    string
		layout  ssg.Layout
		wantErr bool
	}{
		{
			name: "updates layout successfully",
			layout: ssg.Layout{
				ID:     layout.ID,
				SiteID: layout.SiteID,
				Name:   "Updated Layout",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateLayout(ctx, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLayout() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetLayout(ctx, tt.layout.ID)
				if updated.Name != tt.layout.Name {
					t.Errorf("UpdateLayout() name = %v, want %v", updated.Name, tt.layout.Name)
				}
			}
		})
	}
}

func TestClioRepoDeleteLayout(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	layout := ssg.Layout{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Test Layout",
	}
	repo.CreateLayout(ctx, layout)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes layout successfully",
			id:      layout.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteLayout(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLayout() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetLayout(ctx, tt.id)
				if err == nil {
					t.Error("DeleteLayout() layout still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoCreateImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	tests := []struct {
		name    string
		image   *ssg.Image
		wantErr bool
	}{
		{
			name: "creates image successfully",
			image: &ssg.Image{
				ID:       uuid.New(),
				SiteID:   siteID,
				FileName: "test.jpg",
				FilePath: "/images/test.jpg",
				AltText:  "Test Image",
			},
			wantErr: false,
		},
		{
			name: "creates image with all fields",
			image: &ssg.Image{
				ID:       uuid.New(),
				SiteID:   siteID,
				ShortID:  "img001",
				FileName: "full-test.jpg",
				FilePath: "/images/full-test.jpg",
				AltText:  "Full Test Image",
				Title:    "Test Title",
				Width:    1920,
				Height:   1080,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateImage(ctx, tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		FileName: "test.jpg",
		FilePath: "/images/test.jpg",
		AltText:  "Test Image",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "gets existing image",
			id:      image.ID,
			wantErr: false,
		},
		{
			name:    "fails with non-existent image",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetImage(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID != tt.id {
				t.Errorf("GetImage() got ID = %v, want %v", got.ID, tt.id)
			}
		})
	}
}

func TestClioRepoListImages(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*ClioRepo, uuid.UUID)
		wantCount  int
		wantErr    bool
	}{
		{
			name: "lists multiple images successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				image1 := &ssg.Image{
					ID:       uuid.New(),
					SiteID:   siteID,
					FileName: "test1.jpg",
					FilePath: "/images/test1.jpg",
				}
				image2 := &ssg.Image{
					ID:       uuid.New(),
					SiteID:   siteID,
					FileName: "test2.jpg",
					FilePath: "/images/test2.jpg",
				}
				r.CreateImage(ctx, image1)
				r.CreateImage(ctx, image2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no images",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			images, err := repo.ListImages(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(images) != tt.wantCount {
				t.Errorf("ListImages() got %d images, want %d", len(images), tt.wantCount)
			}
		})
	}
}

func TestClioRepoUpdateImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		FileName: "test.jpg",
		FilePath: "/images/test.jpg",
		AltText:  "Original Alt Text",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name    string
		image   *ssg.Image
		wantErr bool
	}{
		{
			name: "updates image successfully",
			image: &ssg.Image{
				ID:       image.ID,
				SiteID:   image.SiteID,
				FileName: image.FileName,
				FilePath: image.FilePath,
				AltText:  "Updated Alt Text",
			},
			wantErr: false,
		},
		{
			name: "updates image with title and dimensions",
			image: &ssg.Image{
				ID:       image.ID,
				SiteID:   image.SiteID,
				FileName: image.FileName,
				FilePath: image.FilePath,
				AltText:  "Updated Alt Text",
				Title:    "New Title",
				Width:    800,
				Height:   600,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateImage(ctx, tt.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				updated, _ := repo.GetImage(ctx, tt.image.ID)
				if updated.AltText != tt.image.AltText {
					t.Errorf("UpdateImage() altText = %v, want %v", updated.AltText, tt.image.AltText)
				}
			}
		})
	}
}

func TestClioRepoDeleteImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		FileName: "test.jpg",
		FilePath: "/images/test.jpg",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "deletes image successfully",
			id:      image.ID,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteImage(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteImage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				_, err := repo.GetImage(ctx, tt.id)
				if err == nil {
					t.Error("DeleteImage() image still exists after deletion")
				}
			}
		})
	}
}

func TestClioRepoAddTagToContent(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Test Content",
	}
	repo.CreateContent(ctx, content)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)

	tests := []struct {
		name      string
		contentID uuid.UUID
		tagID     uuid.UUID
		wantErr   bool
	}{
		{
			name:      "adds tag to content successfully",
			contentID: content.ID,
			tagID:     tag.ID,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.AddTagToContent(ctx, tt.contentID, tt.tagID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTagToContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoRemoveTagFromContent(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Test Content",
	}
	repo.CreateContent(ctx, content)

	tag := ssg.Tag{
		ID:        uuid.New(),
		SiteID:    siteID,
		Name:      "golang",
		SlugField: "golang",
	}
	repo.CreateTag(ctx, tag)
	repo.AddTagToContent(ctx, content.ID, tag.ID)

	tests := []struct {
		name      string
		contentID uuid.UUID
		tagID     uuid.UUID
		wantErr   bool
	}{
		{
			name:      "removes tag from content successfully",
			contentID: content.ID,
			tagID:     tag.ID,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.RemoveTagFromContent(ctx, tt.contentID, tt.tagID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveTagFromContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetTagsForContent(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID) uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple tags for content successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				content := &ssg.Content{
					ID:      uuid.New(),
					SiteID:  siteID,
					Heading: "Test Content",
				}
				r.CreateContent(ctx, content)

				tag1 := ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "golang",
					SlugField: "golang",
				}
				tag2 := ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "rust",
					SlugField: "rust",
				}
				r.CreateTag(ctx, tag1)
				r.CreateTag(ctx, tag2)
				r.AddTagToContent(ctx, content.ID, tag1.ID)
				r.AddTagToContent(ctx, content.ID, tag2.ID)
				return content.ID
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when content has no tags",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				content := &ssg.Content{
					ID:      uuid.New(),
					SiteID:  siteID,
					Heading: "Content Without Tags",
				}
				r.CreateContent(ctx, content)
				return content.ID
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			contentID := tt.setup(repo, siteID)

			tags, err := repo.GetTagsForContent(ctx, contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsForContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tags) != tt.wantCount {
				t.Errorf("GetTagsForContent() got %d tags, want %d", len(tags), tt.wantCount)
			}
		})
	}
}

func TestClioRepoGetAllContentWithMeta(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID)
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets multiple contents with meta successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				content1 := &ssg.Content{
					ID:      uuid.New(),
					SiteID:  siteID,
					Heading: "Content 1",
				}
				content2 := &ssg.Content{
					ID:      uuid.New(),
					SiteID:  siteID,
					Heading: "Content 2",
				}
				r.CreateContent(ctx, content1)
				r.CreateContent(ctx, content2)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no content",
			setup: func(r *ClioRepo, siteID uuid.UUID) {
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tt.setup(repo, siteID)

			contents, err := repo.GetAllContentWithMeta(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllContentWithMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(contents) != tt.wantCount {
				t.Errorf("GetAllContentWithMeta() got %d contents, want %d", len(contents), tt.wantCount)
			}
		})
	}
}

func TestClioRepoGetSiteBySlug(t *testing.T) {
	repo, _ := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := context.Background()

	tests := []struct {
		name     string
		slug     string
		wantSlug string
		wantErr  bool
	}{
		{
			name:     "gets site by slug successfully",
			slug:     "test-site",
			wantSlug: "test-site",
			wantErr:  false,
		},
		{
			name:    "returns error when site not found",
			slug:    "nonexistent-site",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site, err := repo.GetSiteBySlug(ctx, tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSiteBySlug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && site.Slug() != tt.wantSlug {
				t.Errorf("GetSiteBySlug() slug = %v, want %v", site.Slug(), tt.wantSlug)
			}
		})
	}
}

func TestClioRepoGetImageByShortID(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img123",
		FileName: "test.jpg",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name         string
		shortID      string
		wantShortID  string
		wantErr      bool
	}{
		{
			name:        "gets image by short ID successfully",
			shortID:     "img123",
			wantShortID: "img123",
			wantErr:     false,
		},
		{
			name:    "returns error when image not found",
			shortID: "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieved, err := repo.GetImageByShortID(ctx, tt.shortID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageByShortID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && retrieved.ShortID != tt.wantShortID {
				t.Errorf("GetImageByShortID() shortID = %v, want %v", retrieved.ShortID, tt.wantShortID)
			}
		})
	}
}

func TestClioRepoGetImageByContentHash(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img456",
		FileName: "test.jpg",
		FilePath: "images/test_hash.jpg",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name         string
		contentHash  string
		wantFilePath string
		wantErr      bool
	}{
		{
			name:         "gets image by content hash successfully",
			contentHash:  "images/test_hash.jpg",
			wantFilePath: "images/test_hash.jpg",
			wantErr:      false,
		},
		{
			name:        "returns error when image not found",
			contentHash: "nonexistent_hash",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieved, err := repo.GetImageByContentHash(ctx, tt.contentHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageByContentHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && retrieved.FilePath != tt.wantFilePath {
				t.Errorf("GetImageByContentHash() filePath = %v, want %v", retrieved.FilePath, tt.wantFilePath)
			}
		})
	}
}

func TestClioRepoGetContentForTag(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID) uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name: "gets content for tag successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				tag := &ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "test-tag",
					SlugField: "test-tag",
				}
				r.CreateTag(ctx, *tag)

				content := &ssg.Content{
					ID:      uuid.New(),
					SiteID:  siteID,
					Heading: "Tagged Content",
				}
				r.CreateContent(ctx, content)
				r.AddTagToContent(ctx, content.ID, tag.ID)
				return tag.ID
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "returns empty list when tag has no content",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				tag := &ssg.Tag{
					ID:        uuid.New(),
					SiteID:    siteID,
					Name:      "unused-tag",
					SlugField: "unused-tag",
				}
				r.CreateTag(ctx, *tag)
				return tag.ID
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			tagID := tt.setup(repo, siteID)

			contents, err := repo.GetContentForTag(ctx, tagID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentForTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(contents) != tt.wantCount {
				t.Errorf("GetContentForTag() got %d contents, want %d", len(contents), tt.wantCount)
			}
		})
	}
}

func TestClioRepoCreateImageVariant(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img789",
		FileName: "test.jpg",
	}
	repo.CreateImage(ctx, image)

	tests := []struct {
		name    string
		variant *ssg.ImageVariant
		wantErr bool
	}{
		{
			name: "creates thumbnail variant successfully",
			variant: &ssg.ImageVariant{
				ID:      uuid.New(),
				ImageID: image.ID,
				Kind:    "thumbnail",
				Width:   150,
				Height:  150,
				BlobRef: "test_thumb.jpg",
			},
			wantErr: false,
		},
		{
			name: "creates original variant with full fields",
			variant: &ssg.ImageVariant{
				ID:            uuid.New(),
				ShortID:       "var001",
				ImageID:       image.ID,
				Kind:          "original",
				Width:         1920,
				Height:        1080,
				FilesizeByte:  2048000,
				Mime:          "image/jpeg",
				BlobRef:       "test_original.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateImageVariant(ctx, tt.variant)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoGetImageVariant(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img101",
		FileName: "test.jpg",
	}
	repo.CreateImage(ctx, image)

	variant := &ssg.ImageVariant{
		ID:      uuid.New(),
		ImageID: image.ID,
		Kind:    "medium",
		Width:   500,
		Height:  500,
		BlobRef: "test_medium.jpg",
	}
	repo.CreateImageVariant(ctx, variant)

	tests := []struct {
		name     string
		id       uuid.UUID
		wantKind string
		wantErr  bool
	}{
		{
			name:     "gets image variant successfully",
			id:       variant.ID,
			wantKind: "medium",
			wantErr:  false,
		},
		{
			name:    "returns error when variant not found",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrieved, err := repo.GetImageVariant(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetImageVariant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && retrieved.Kind != tt.wantKind {
				t.Errorf("GetImageVariant() kind = %v, want %v", retrieved.Kind, tt.wantKind)
			}
		})
	}
}

func TestClioRepoListImageVariantsByImageID(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*ClioRepo, uuid.UUID) uuid.UUID
		wantCount int
		wantErr   bool
	}{
		{
			name: "lists multiple variants successfully",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				image := &ssg.Image{
					ID:       uuid.New(),
					SiteID:   siteID,
					ShortID:  "img202",
					FileName: "test.jpg",
				}
				r.CreateImage(ctx, image)

				variant1 := &ssg.ImageVariant{
					ID:      uuid.New(),
					ImageID: image.ID,
					Kind:    "small",
					BlobRef: "test_small.jpg",
				}
				variant2 := &ssg.ImageVariant{
					ID:      uuid.New(),
					ImageID: image.ID,
					Kind:    "large",
					BlobRef: "test_large.jpg",
				}
				r.CreateImageVariant(ctx, variant1)
				r.CreateImageVariant(ctx, variant2)
				return image.ID
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "returns empty list when no variants",
			setup: func(r *ClioRepo, siteID uuid.UUID) uuid.UUID {
				ctx := ssg.NewContextWithSite("test-site", siteID)
				image := &ssg.Image{
					ID:       uuid.New(),
					SiteID:   siteID,
					ShortID:  "img203",
					FileName: "test.jpg",
				}
				r.CreateImage(ctx, image)
				return image.ID
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, siteID := setupTestSsgRepo(t)
			defer repo.db.Close()
			ctx := ssg.NewContextWithSite("test-site", siteID)

			imageID := tt.setup(repo, siteID)

			variants, err := repo.ListImageVariantsByImageID(ctx, imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListImageVariantsByImageID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(variants) != tt.wantCount {
				t.Errorf("ListImageVariantsByImageID() got %d variants, want %d", len(variants), tt.wantCount)
			}
		})
	}
}

func TestClioRepoUpdateImageVariant(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img303",
		FileName: "test.jpg",
	}
	repo.CreateImage(ctx, image)

	variant := &ssg.ImageVariant{
		ID:      uuid.New(),
		ImageID: image.ID,
		Kind:    "original",
		BlobRef: "test_original.jpg",
	}
	repo.CreateImageVariant(ctx, variant)

	tests := []struct {
		name    string
		variant *ssg.ImageVariant
		wantErr bool
	}{
		{
			name: "updates variant kind successfully",
			variant: &ssg.ImageVariant{
				ID:      variant.ID,
				ImageID: variant.ImageID,
				Kind:    "updated",
				BlobRef: variant.BlobRef,
			},
			wantErr: false,
		},
		{
			name: "updates variant with full fields",
			variant: &ssg.ImageVariant{
				ID:            variant.ID,
				ImageID:       variant.ImageID,
				Kind:          "web",
				BlobRef:       "test_web.jpg",
				Width:         800,
				Height:        600,
				FilesizeByte:  512000,
				Mime:          "image/jpeg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateImageVariant(ctx, tt.variant)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateImageVariant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClioRepoDeleteImageVariant(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img404",
		FileName: "test.jpg",
	}
	repo.CreateImage(ctx, image)

	variant := &ssg.ImageVariant{
		ID:      uuid.New(),
		ImageID: image.ID,
		Kind:    "todelete",
		BlobRef: "test_delete.jpg",
	}
	repo.CreateImageVariant(ctx, variant)

	err := repo.DeleteImageVariant(ctx, variant.ID)
	if err != nil {
		t.Errorf("DeleteImageVariant() error = %v", err)
	}
}

func TestClioRepoCreateContentImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Content with image",
	}
	repo.CreateContent(ctx, content)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img505",
		FileName: "content.jpg",
	}
	repo.CreateImage(ctx, image)

	contentImage := &ssg.ContentImage{
		ID:        uuid.New(),
		ContentID: content.ID,
		ImageID:   image.ID,
	}

	err := repo.CreateContentImage(ctx, contentImage)
	if err != nil {
		t.Errorf("CreateContentImage() error = %v", err)
	}
}

func TestClioRepoGetContentImagesByContentID(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Content with images",
	}
	repo.CreateContent(ctx, content)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img606",
		FileName: "content2.jpg",
	}
	repo.CreateImage(ctx, image)

	contentImage := &ssg.ContentImage{
		ID:        uuid.New(),
		ContentID: content.ID,
		ImageID:   image.ID,
	}
	repo.CreateContentImage(ctx, contentImage)

	images, err := repo.GetContentImagesByContentID(ctx, content.ID)
	if err != nil {
		t.Errorf("GetContentImagesByContentID() error = %v", err)
		return
	}

	if len(images) == 0 {
		t.Error("GetContentImagesByContentID() returned no images")
	}
}

func TestClioRepoDeleteContentImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	content := &ssg.Content{
		ID:      uuid.New(),
		SiteID:  siteID,
		Heading: "Content",
	}
	repo.CreateContent(ctx, content)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img707",
		FileName: "content3.jpg",
	}
	repo.CreateImage(ctx, image)

	contentImage := &ssg.ContentImage{
		ID:        uuid.New(),
		ContentID: content.ID,
		ImageID:   image.ID,
	}
	repo.CreateContentImage(ctx, contentImage)

	err := repo.DeleteContentImage(ctx, contentImage.ID)
	if err != nil {
		t.Errorf("DeleteContentImage() error = %v", err)
	}
}

func TestClioRepoCreateSectionImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := &ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Section with image",
	}
	repo.CreateSection(ctx, *section)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img808",
		FileName: "section.jpg",
	}
	repo.CreateImage(ctx, image)

	sectionImage := &ssg.SectionImage{
		ID:        uuid.New(),
		SectionID: section.ID,
		ImageID:   image.ID,
	}

	err := repo.CreateSectionImage(ctx, sectionImage)
	if err != nil {
		t.Errorf("CreateSectionImage() error = %v", err)
	}
}

func TestClioRepoGetSectionImagesBySectionID(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := &ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Section with images",
	}
	repo.CreateSection(ctx, *section)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img909",
		FileName: "section2.jpg",
	}
	repo.CreateImage(ctx, image)

	sectionImage := &ssg.SectionImage{
		ID:        uuid.New(),
		SectionID: section.ID,
		ImageID:   image.ID,
	}
	repo.CreateSectionImage(ctx, sectionImage)

	images, err := repo.GetSectionImagesBySectionID(ctx, section.ID)
	if err != nil {
		t.Errorf("GetSectionImagesBySectionID() error = %v", err)
		return
	}

	if len(images) == 0 {
		t.Error("GetSectionImagesBySectionID() returned no images")
	}
}

func TestClioRepoDeleteSectionImage(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	section := &ssg.Section{
		ID:     uuid.New(),
		SiteID: siteID,
		Name:   "Section",
	}
	repo.CreateSection(ctx, *section)

	image := &ssg.Image{
		ID:       uuid.New(),
		SiteID:   siteID,
		ShortID:  "img000",
		FileName: "section3.jpg",
	}
	repo.CreateImage(ctx, image)

	sectionImage := &ssg.SectionImage{
		ID:        uuid.New(),
		SectionID: section.ID,
		ImageID:   image.ID,
	}
	repo.CreateSectionImage(ctx, sectionImage)

	err := repo.DeleteSectionImage(ctx, sectionImage.ID)
	if err != nil {
		t.Errorf("DeleteSectionImage() error = %v", err)
	}
}

func TestClioRepoGetContentWithPaginationAndSearch(t *testing.T) {
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
	ctx := ssg.NewContextWithSite("test-site", siteID)

	for i := 0; i < 5; i++ {
		content := &ssg.Content{
			ID:      uuid.New(),
			SiteID:  siteID,
			Heading: "Searchable Content",
		}
		repo.CreateContent(ctx, content)
	}

	contents, total, err := repo.GetContentWithPaginationAndSearch(ctx, 0, 10, "Searchable")
	if err != nil {
		t.Errorf("GetContentWithPaginationAndSearch() error = %v", err)
		return
	}

	if len(contents) == 0 {
		t.Error("GetContentWithPaginationAndSearch() returned no contents")
	}

	if total == 0 {
		t.Error("GetContentWithPaginationAndSearch() total = 0")
	}
}

func TestSanitizeURLPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "handles simple path",
			input:    "blog/my-post",
			expected: "blog/my-post",
		},
		{
			name:     "converts to lowercase",
			input:    "Blog/My-Post",
			expected: "blog/my-post",
		},
		{
			name:     "replaces special characters with hyphens",
			input:    "blog/my@post!test",
			expected: "blog/my-post-test",
		},
		{
			name:     "removes multiple consecutive hyphens",
			input:    "blog/my---post",
			expected: "blog/my-post",
		},
		{
			name:     "removes leading and trailing hyphens",
			input:    "blog/-my-post-",
			expected: "blog/my-post",
		},
		{
			name:     "preserves underscores and dots",
			input:    "blog/my_post.html",
			expected: "blog/my_post.html",
		},
		{
			name:     "handles empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "preserves leading slash",
			input:    "/blog/post",
			expected: "/blog/post",
		},
		{
			name:     "handles complex special characters",
			input:    "blog/post with spaces & symbols!",
			expected: "blog/post-with-spaces-symbols",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeURLPath(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeURLPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
