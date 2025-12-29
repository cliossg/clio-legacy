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
			template TEXT,
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
	repo, siteID := setupTestSsgRepo(t)
	defer repo.db.Close()
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
	repo.CreateSection(ctx, section1)
	repo.CreateSection(ctx, section2)

	sections, err := repo.GetSections(ctx)
	if err != nil {
		t.Errorf("GetSections() error = %v", err)
		return
	}

	if len(sections) != 2 {
		t.Errorf("GetSections() got %d sections, want 2", len(sections))
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
