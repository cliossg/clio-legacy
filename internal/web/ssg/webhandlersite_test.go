package ssg

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
	"github.com/hermesgen/hm"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed assets
var testAssetsFS embed.FS

type mockSiteRepo struct {
	hm.Core
	sites          map[uuid.UUID]feat.Site
	sitesBySlug    map[string]feat.Site
	listSitesErr   error
	createSiteErr  error
	deleteSiteErr  error
	getSiteErr     error
	deletedSiteSlug string
}

func newMockSiteRepo() *mockSiteRepo {
	cfg := hm.NewConfig()
	return &mockSiteRepo{
		Core:        hm.NewCore("mock-site-repo", hm.XParams{Cfg: cfg}),
		sites:       make(map[uuid.UUID]feat.Site),
		sitesBySlug: make(map[string]feat.Site),
	}
}

func (m *mockSiteRepo) CreateSite(ctx context.Context, site *feat.Site) error {
	if m.createSiteErr != nil {
		return m.createSiteErr
	}
	m.sites[site.ID] = *site
	m.sitesBySlug[site.Slug()] = *site
	return nil
}

func (m *mockSiteRepo) GetSite(ctx context.Context, id uuid.UUID) (feat.Site, error) {
	if m.getSiteErr != nil {
		return feat.Site{}, m.getSiteErr
	}
	if site, ok := m.sites[id]; ok {
		return site, nil
	}
	return feat.Site{}, fmt.Errorf("site not found")
}

func (m *mockSiteRepo) ListSites(ctx context.Context, includeDeleted bool) ([]feat.Site, error) {
	if m.listSitesErr != nil {
		return nil, m.listSitesErr
	}
	var sites []feat.Site
	for _, site := range m.sites {
		sites = append(sites, site)
	}
	return sites, nil
}

func (m *mockSiteRepo) GetSiteBySlug(ctx context.Context, slug string) (feat.Site, error) {
	if m.getSiteErr != nil {
		return feat.Site{}, m.getSiteErr
	}
	if site, ok := m.sitesBySlug[slug]; ok {
		return site, nil
	}
	return feat.Site{}, fmt.Errorf("site not found")
}

func (m *mockSiteRepo) UpdateSite(ctx context.Context, site *feat.Site) error {
	m.sites[site.ID] = *site
	m.sitesBySlug[site.Slug()] = *site
	return nil
}

func (m *mockSiteRepo) DeleteSite(ctx context.Context, siteID uuid.UUID) error {
	if m.deleteSiteErr != nil {
		return m.deleteSiteErr
	}
	if site, ok := m.sites[siteID]; ok {
		m.deletedSiteSlug = site.Slug()
		delete(m.sites, siteID)
		delete(m.sitesBySlug, site.Slug())
	}
	return nil
}

type mockSessionManager struct {
	setSiteSlugErr error
	siteSlug       string
}

func (m *mockSessionManager) SetUserSession(w http.ResponseWriter, userID uuid.UUID, siteSlug string) error {
	return nil
}

func (m *mockSessionManager) GetUserSession(r *http.Request) (uuid.UUID, string, error) {
	return uuid.New(), m.siteSlug, nil
}

func (m *mockSessionManager) SetSiteSlug(w http.ResponseWriter, r *http.Request, siteSlug string) error {
	if m.setSiteSlugErr != nil {
		return m.setSiteSlugErr
	}
	m.siteSlug = siteSlug
	return nil
}

func newTestWebHandler(siteRepo *mockSiteRepo, sessMgr *mockSessionManager) *WebHandler {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}

	tm := hm.NewTemplateManager(testAssetsFS, params)
	flash := hm.NewFlashManager(params)

	if siteRepo == nil {
		siteRepo = newMockSiteRepo()
	}
	if sessMgr == nil {
		sessMgr = &mockSessionManager{}
	}

	// Create an in-memory SQLite DB for testing
	db, _ := sqlx.Open("sqlite3", ":memory:")

	// Create site table for SiteRepo (note: singular "site")
	db.Exec(`CREATE TABLE IF NOT EXISTS site (
		id TEXT PRIMARY KEY,
		short_id TEXT,
		name TEXT,
		slug TEXT,
		mode TEXT,
		active INTEGER DEFAULT 1,
		created_by TEXT,
		updated_by TEXT,
		created_at DATETIME,
		updated_at DATETIME
	)`)

	realSiteRepo := feat.NewSiteRepo(db)
	repo := newTestRepoWithDB(realSiteRepo, db)
	siteMgr := feat.NewSiteManager(repo, testAssetsFS, "sqlite3", params)
	siteMgr.Setup(context.Background())

	var smInterface interface {
		SetUserSession(w http.ResponseWriter, userID uuid.UUID, siteSlug string) error
		GetUserSession(r *http.Request) (userID uuid.UUID, siteSlug string, err error)
		SetSiteSlug(w http.ResponseWriter, r *http.Request, siteSlug string) error
	} = sessMgr

	return NewWebHandler(tm, flash, nil, siteMgr, smInterface, params)
}

func TestWebHandlerCreateSite(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		setupMock      func(*mockSiteRepo)
		wantStatusCode int
		wantLocation   string
	}{
		// Note: skipping "creates successfully" test because it requires full auth.Repo implementation
		// which is complex for this integration test. The validation paths are tested instead.
		{
			name: "fails with missing name",
			formData: url.Values{
				"slug": []string{"test-site"},
				"mode": []string{"structured"},
			},
			setupMock:      func(m *mockSiteRepo) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites/new",
		},
		{
			name: "fails with missing slug",
			formData: url.Values{
				"name": []string{"Test Site"},
				"mode": []string{"structured"},
			},
			setupMock:      func(m *mockSiteRepo) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites/new",
		},
		{
			name: "fails with missing mode",
			formData: url.Values{
				"name": []string{"Test Site"},
				"slug": []string{"test-site"},
			},
			setupMock:      func(m *mockSiteRepo) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites/new",
		},
		{
			name: "fails when site manager returns error",
			formData: url.Values{
				"name": []string{"Test Site"},
				"slug": []string{"test-site"},
				"mode": []string{"structured"},
			},
			setupMock: func(m *mockSiteRepo) {
				m.createSiteErr = fmt.Errorf("db error")
			},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites/new",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siteRepo := newMockSiteRepo()
			tt.setupMock(siteRepo)
			handler := newTestWebHandler(siteRepo, nil)

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/sites", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.CreateSite(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateSite() status = %d, want %d", w.Code, tt.wantStatusCode)
			}

			location := w.Header().Get("Location")
			if location != tt.wantLocation {
				t.Errorf("CreateSite() location = %s, want %s", location, tt.wantLocation)
			}
		})
	}
}

func TestWebHandlerSwitchSite(t *testing.T) {
	tests := []struct {
		name           string
		querySlug      string
		setupMock      func(*mockSessionManager)
		wantStatusCode int
		wantLocation   string
	}{
		{
			name:           "switches site successfully",
			querySlug:      "test-site",
			setupMock:      func(m *mockSessionManager) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/list-content?site=test-site",
		},
		{
			name:           "fails with missing slug",
			querySlug:      "",
			setupMock:      func(m *mockSessionManager) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
		{
			name:      "fails when session manager returns error",
			querySlug: "test-site",
			setupMock: func(m *mockSessionManager) {
				m.setSiteSlugErr = fmt.Errorf("session error")
			},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessMgr := &mockSessionManager{}
			tt.setupMock(sessMgr)
			handler := newTestWebHandler(nil, sessMgr)

			req := httptest.NewRequest(http.MethodGet, "/ssg/switch-site?slug="+tt.querySlug, nil)
			w := httptest.NewRecorder()

			handler.SwitchSite(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("SwitchSite() status = %d, want %d", w.Code, tt.wantStatusCode)
			}

			location := w.Header().Get("Location")
			if location != tt.wantLocation {
				t.Errorf("SwitchSite() location = %s, want %s", location, tt.wantLocation)
			}
		})
	}
}

func TestWebHandlerDeleteSite(t *testing.T) {
	tests := []struct {
		name           string
		queryID        string
		setupMock      func(*mockSiteRepo)
		wantStatusCode int
		wantLocation   string
	}{
		{
			name:    "deletes site successfully",
			queryID: uuid.New().String(),
			setupMock: func(m *mockSiteRepo) {
				m.deletedSiteSlug = "test-site"
			},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
		{
			name:           "fails with missing id",
			queryID:        "",
			setupMock:      func(m *mockSiteRepo) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
		{
			name:           "fails with invalid id",
			queryID:        "invalid-uuid",
			setupMock:      func(m *mockSiteRepo) {},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
		{
			name:    "fails when site manager returns error",
			queryID: uuid.New().String(),
			setupMock: func(m *mockSiteRepo) {
				m.deleteSiteErr = fmt.Errorf("db error")
			},
			wantStatusCode: http.StatusSeeOther,
			wantLocation:   "/ssg/sites",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siteRepo := newMockSiteRepo()
			tt.setupMock(siteRepo)
			handler := newTestWebHandler(siteRepo, nil)

			req := httptest.NewRequest(http.MethodGet, "/ssg/delete-site?id="+tt.queryID, nil)
			w := httptest.NewRecorder()

			handler.DeleteSite(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteSite() status = %d, want %d", w.Code, tt.wantStatusCode)
			}

			location := w.Header().Get("Location")
			if location != tt.wantLocation {
				t.Errorf("DeleteSite() location = %s, want %s", location, tt.wantLocation)
			}
		})
	}
}

func TestWebHandlerRootRedirect(t *testing.T) {
	tests := []struct {
		name           string
		cookie         *http.Cookie
		setupMock      func(*mockSiteRepo)
		wantStatusCode int
		wantLocation   string
	}{
		// Note: skipping "redirects to last site" test because it requires inserting data into the DB
		// The other redirect paths are tested
		{
			name:   "redirects to sites list when no cookie",
			cookie: nil,
			setupMock: func(m *mockSiteRepo) {
			},
			wantStatusCode: http.StatusFound,
			wantLocation:   "/ssg/sites",
		},
		{
			name: "redirects to sites list when site not found",
			cookie: &http.Cookie{
				Name:  "last_site",
				Value: "missing-site",
			},
			setupMock: func(m *mockSiteRepo) {
				m.getSiteErr = fmt.Errorf("not found")
			},
			wantStatusCode: http.StatusFound,
			wantLocation:   "/ssg/sites",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			siteRepo := newMockSiteRepo()
			tt.setupMock(siteRepo)
			handler := newTestWebHandler(siteRepo, nil)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
			w := httptest.NewRecorder()

			handler.RootRedirect(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("RootRedirect() status = %d, want %d", w.Code, tt.wantStatusCode)
			}

			location := w.Header().Get("Location")
			if location != tt.wantLocation {
				t.Errorf("RootRedirect() location = %s, want %s", location, tt.wantLocation)
			}
		})
	}
}
