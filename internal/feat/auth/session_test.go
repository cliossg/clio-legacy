package auth

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/hermesgen/hm"
)

func setupSessionManager(t *testing.T) *SessionManager {
	cfg := hm.NewConfig()
	hashKey := []byte("01234567890123456789012345678901")
	blockKey := []byte("01234567890123456789012345678901")
	cfg.Set(hm.Key.SecHashKey, hashKey)
	cfg.Set(hm.Key.SecBlockKey, blockKey)

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

