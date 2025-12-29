package ssg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hermesgen/hm"
)

func TestAPIHandlerPublish(t *testing.T) {
	repo := newMockServiceRepo()
	svc := newTestService(repo)
	apiHandler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: hm.NewConfig()})

	publishReq := PublishRequest{
		Message: "Test publish",
	}
	body, _ := json.Marshal(publishReq)

	req := httptest.NewRequest("POST", "/publish", bytes.NewReader(body))
	w := httptest.NewRecorder()

	apiHandler.Publish(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Publish() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAPIHandlerGenerateMarkdown(t *testing.T) {
	repo := newMockServiceRepo()
	svc := newTestService(repo)
	apiHandler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: hm.NewConfig()})

	req := httptest.NewRequest("POST", "/generate-markdown", nil)
	w := httptest.NewRecorder()

	apiHandler.GenerateMarkdown(w, req)

	if w.Code == http.StatusOK || w.Code == http.StatusInternalServerError {
		t.Logf("GenerateMarkdown() status = %d", w.Code)
	}
}

func TestAPIHandlerGenerateHTML(t *testing.T) {
	repo := newMockServiceRepo()
	svc := newTestService(repo)
	apiHandler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: hm.NewConfig()})

	req := httptest.NewRequest("POST", "/generate-html", nil)
	w := httptest.NewRecorder()

	apiHandler.GenerateHTML(w, req)

	if w.Code == http.StatusOK || w.Code == http.StatusInternalServerError {
		t.Logf("GenerateHTML() status = %d", w.Code)
	}
}
