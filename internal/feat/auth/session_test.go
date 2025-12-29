package auth

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

func setupSessionManager(t *testing.T) *SessionManager {
	cfg := hm.NewConfig()
	cfg.Set(hm.Key.SecHashKey, "01234567890123456789012345678901")
	cfg.Set(hm.Key.SecBlockKey, "01234567890123456789012345678901")

	log := hm.NewLogger("error")
	params := hm.XParams{Cfg: cfg, Log: log}

	sm := NewSessionManager(params)
	err := sm.Setup(context.Background())
	if err != nil {
		t.Fatalf("setup session manager: %v", err)
	}

	return sm
}

func TestNewSessionManager(t *testing.T) {
	cfg := hm.NewConfig()
	log := hm.NewLogger("error")
	params := hm.XParams{Cfg: cfg, Log: log}

	sm := NewSessionManager(params)
	if sm == nil {
		t.Fatal("NewSessionManager() returned nil")
	}
}

func TestSessionManagerSetup(t *testing.T) {
	sm := setupSessionManager(t)
	if sm == nil {
		t.Fatal("Setup() returned nil session manager")
	}
	if sm.encoder == nil {
		t.Fatal("Setup() did not initialize encoder")
	}
}

func TestSessionManagerClearUserSession(t *testing.T) {
	sm := setupSessionManager(t)
	w := httptest.NewRecorder()

	sm.ClearUserSession(w)

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("ClearUserSession() did not set cookie")
		return
	}

	found := false
	for _, cookie := range cookies {
		if cookie.Name == sessionCookieName {
			found = true
			if cookie.MaxAge != -1 {
				t.Errorf("ClearUserSession() MaxAge = %v, want -1", cookie.MaxAge)
			}
		}
	}

	if !found {
		t.Error("ClearUserSession() did not set session cookie")
	}
}

func TestSessionManagerSetUserSession(t *testing.T) {
	sm := setupSessionManager(t)
	w := httptest.NewRecorder()
	userID := uuid.New()
	siteSlug := "test-site"

	err := sm.SetUserSession(w, userID, siteSlug)
	if err != nil {
		t.Errorf("SetUserSession() error = %v", err)
		return
	}

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("SetUserSession() did not set cookie")
		return
	}

	found := false
	for _, cookie := range cookies {
		if cookie.Name == sessionCookieName {
			found = true
			if cookie.Value == "" {
				t.Error("SetUserSession() cookie value is empty")
			}
			if !cookie.HttpOnly {
				t.Error("SetUserSession() cookie is not HttpOnly")
			}
			if !cookie.Secure {
				t.Error("SetUserSession() cookie is not Secure")
			}
		}
	}

	if !found {
		t.Error("SetUserSession() did not set session cookie")
	}
}

func TestSessionManagerGetUserSession(t *testing.T) {
	sm := setupSessionManager(t)
	w := httptest.NewRecorder()
	expectedUserID := uuid.New()
	expectedSiteSlug := "test-site"

	err := sm.SetUserSession(w, expectedUserID, expectedSiteSlug)
	if err != nil {
		t.Fatalf("SetUserSession() error = %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	userID, siteSlug, err := sm.GetUserSession(req)
	if err != nil {
		t.Errorf("GetUserSession() error = %v", err)
		return
	}

	if userID != expectedUserID {
		t.Errorf("GetUserSession() userID = %v, want %v", userID, expectedUserID)
	}

	if siteSlug != expectedSiteSlug {
		t.Errorf("GetUserSession() siteSlug = %v, want %v", siteSlug, expectedSiteSlug)
	}
}

func TestSessionManagerSetSiteSlug(t *testing.T) {
	sm := setupSessionManager(t)

	t.Run("without existing session", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		siteSlug := "new-site"

		err := sm.SetSiteSlug(w, req, siteSlug)
		if err != nil {
			t.Errorf("SetSiteSlug() error = %v", err)
			return
		}

		cookies := w.Result().Cookies()
		if len(cookies) == 0 {
			t.Error("SetSiteSlug() did not set cookie")
		}
	})

	t.Run("with existing session", func(t *testing.T) {
		w1 := httptest.NewRecorder()
		userID := uuid.New()
		err := sm.SetUserSession(w1, userID, "old-site")
		if err != nil {
			t.Fatalf("SetUserSession() error = %v", err)
		}

		req := httptest.NewRequest("GET", "/", nil)
		for _, cookie := range w1.Result().Cookies() {
			req.AddCookie(cookie)
		}

		w2 := httptest.NewRecorder()
		newSiteSlug := "new-site"
		err = sm.SetSiteSlug(w2, req, newSiteSlug)
		if err != nil {
			t.Errorf("SetSiteSlug() error = %v", err)
			return
		}

		req2 := httptest.NewRequest("GET", "/", nil)
		for _, cookie := range w2.Result().Cookies() {
			req2.AddCookie(cookie)
		}

		retrievedUserID, retrievedSiteSlug, err := sm.GetUserSession(req2)
		if err != nil {
			t.Errorf("GetUserSession() error = %v", err)
			return
		}

		if retrievedUserID != userID {
			t.Errorf("SetSiteSlug() changed userID: got %v, want %v", retrievedUserID, userID)
		}

		if retrievedSiteSlug != newSiteSlug {
			t.Errorf("SetSiteSlug() siteSlug = %v, want %v", retrievedSiteSlug, newSiteSlug)
		}
	})
}

