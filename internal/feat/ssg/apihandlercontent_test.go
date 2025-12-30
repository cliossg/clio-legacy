package ssg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

func TestAPIHandlerGetAllContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "gets all content successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.contents[uuid.New()] = Content{ID: uuid.New(), Heading: "Content 1"}
				m.contents[uuid.New()] = Content{ID: uuid.New(), Heading: "Content 2"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no content",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getContentErr = fmt.Errorf("db error")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodGet, "/ssg/contents", nil)
			w := httptest.NewRecorder()

			handler.GetAllContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetAllContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		contentID      string
		wantStatusCode int
	}{
		{
			name: "gets content successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test Content"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			contentID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when content not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.contentID
			if idStr == "" {
				idStr = contentID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/contents/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerCreateContent(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		setupContext   func() *http.Request
		wantStatusCode int
	}{
		{
			name: "creates content successfully",
			requestBody: map[string]interface{}{
				"heading": "Test Content",
				"body":    "Test body",
			},
			setupRepo: func(m *mockServiceRepo) {},
			setupContext: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/ssg/contents", nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:        "fails with invalid JSON",
			requestBody: "invalid json",
			setupRepo:   func(m *mockServiceRepo) {},
			setupContext: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/ssg/contents", nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails without site ID in context",
			requestBody: map[string]interface{}{
				"heading": "Test Content",
			},
			setupRepo: func(m *mockServiceRepo) {},
			setupContext: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/ssg/contents", nil)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			requestBody: map[string]interface{}{
				"heading": "Test Content",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createContentErr = fmt.Errorf("db error")
			},
			setupContext: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/ssg/contents", nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "creates content with tags",
			requestBody: map[string]interface{}{
				"heading": "Test Content",
				"tags": []map[string]string{
					{"name": "tag1"},
				},
			},
			setupRepo: func(m *mockServiceRepo) {},
			setupContext: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/ssg/contents", nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			req := tt.setupContext()
			req.Body = http.NoBody
			req = httptest.NewRequest(http.MethodPost, "/ssg/contents", bytes.NewReader(body))
			req = req.WithContext(tt.setupContext().Context())
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		setupContext   func(uuid.UUID) *http.Request
		contentID      string
		wantStatusCode int
	}{
		{
			name: "updates content successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Old Heading"}
				return id
			},
			requestBody: map[string]interface{}{
				"heading": "New Heading",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:      "fails with invalid UUID",
			contentID: "invalid-uuid",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.Nil
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/invalid-uuid", nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails with invalid JSON body",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			requestBody: "invalid json",
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when no site selected",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				return httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateContentErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fails when cannot get existing tags",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test"}
				m.getTagsForContentErr = fmt.Errorf("cannot get tags")
				return id
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fails when cannot remove tag from content",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				tagID := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test"}
				m.contentTags[id] = []Tag{{ID: tagID, Name: "test-tag"}}
				m.removeTagFromContentErr = fmt.Errorf("cannot remove tag")
				return id
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fails when cannot add tag to content",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test"}
				m.addTagToContentErr = fmt.Errorf("cannot add tag")
				return id
			},
			requestBody: map[string]interface{}{
				"heading": "Test",
				"tags": []map[string]string{
					{"name": "new-tag"},
				},
			},
			setupContext: func(id uuid.UUID) *http.Request {
				req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+id.String(), nil)
				ctx := addSiteIDToContext(req.Context(), uuid.New())
				return req.WithContext(ctx)
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			idStr := tt.contentID
			if idStr == "" {
				idStr = contentID.String()
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/contents/"+idStr, bytes.NewReader(body))
			req = req.WithContext(tt.setupContext(contentID).Context())
			req.SetPathValue("id", idStr)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		contentID      string
		wantStatusCode int
	}{
		{
			name: "deletes content successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id, Heading: "Test"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			contentID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteContentErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.contentID
			if idStr == "" {
				idStr = contentID.String()
			}

			req := httptest.NewRequest(http.MethodDelete, "/ssg/contents/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.DeleteContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerAddTagToContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		contentID      string
		wantStatusCode int
	}{
		{
			name: "adds tag to content successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id}
				return id
			},
			requestBody: map[string]string{
				"name": "test-tag",
			},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "fails with invalid content UUID",
			contentID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			requestBody:    map[string]string{"name": "test"},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails with invalid JSON body",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.New() },
			requestBody:    "invalid json",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.addTagToContentErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]string{
				"name": "test-tag",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				var err error
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			idStr := tt.contentID
			if idStr == "" {
				idStr = contentID.String()
			}

			req := httptest.NewRequest(http.MethodPost, "/ssg/contents/"+idStr+"/tags", bytes.NewReader(body))
			req.SetPathValue("content_id", idStr)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AddTagToContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("AddTagToContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerRemoveTagFromContent(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) (uuid.UUID, uuid.UUID)
		contentID      string
		tagID          string
		wantStatusCode int
	}{
		{
			name: "removes tag from content successfully",
			setupRepo: func(m *mockServiceRepo) (uuid.UUID, uuid.UUID) {
				contentID := uuid.New()
				tagID := uuid.New()
				m.contents[contentID] = Content{ID: contentID}
				m.tags[tagID] = Tag{ID: tagID}
				return contentID, tagID
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid content UUID",
			contentID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) (uuid.UUID, uuid.UUID) { return uuid.Nil, uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails with invalid tag UUID",
			tagID:          "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) (uuid.UUID, uuid.UUID) { return uuid.New(), uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) (uuid.UUID, uuid.UUID) {
				m.removeTagFromContentErr = fmt.Errorf("db error")
				return uuid.New(), uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID, tagID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			cidStr := tt.contentID
			if cidStr == "" {
				cidStr = contentID.String()
			}
			tidStr := tt.tagID
			if tidStr == "" {
				tidStr = tagID.String()
			}

			req := httptest.NewRequest(http.MethodDelete, "/ssg/contents/"+cidStr+"/tags/"+tidStr, nil)
			req.SetPathValue("content_id", cidStr)
			req.SetPathValue("tag_id", tidStr)
			w := httptest.NewRecorder()

			handler.RemoveTagFromContent(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("RemoveTagFromContent() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetContentImages(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		contentID      string
		wantStatusCode int
	}{
		{
			name: "gets content images successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.contents[id] = Content{ID: id}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid content UUID",
			contentID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.getContentImagesByContentIDErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			contentID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.contentID
			if idStr == "" {
				idStr = contentID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/contents/"+idStr+"/images", nil)
			req.SetPathValue("content_id", idStr)
			w := httptest.NewRecorder()

			handler.GetContentImages(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetContentImages() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

// Note: TestAPIHandlerDeleteContentImage and TestAPIHandlerUploadContentImage are skipped
// because they require ImageManager which is a concrete struct and difficult to mock.
// These handlers will be tested via integration tests or when ImageManager is refactored to use an interface.

// Helper function to add site ID to context
func addSiteIDToContext(ctx context.Context, siteID uuid.UUID) context.Context {
	return context.WithValue(ctx, siteIDKey, siteID)
}
