package ssg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
)

func TestAPIHandlerCreateImage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates image successfully",
			requestBody: map[string]interface{}{
				"fileName": "test.jpg",
				"filePath": "/images/test.jpg",
			},
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "fails with invalid JSON",
			requestBody:    "invalid json",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			requestBody: map[string]string{
				"fileName": "test.jpg",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createImageErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodPost, "/ssg/images", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetImage(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		imageID        string
		wantStatusCode int
	}{
		{
			name: "gets image successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.images[id] = Image{ID: id, FileName: "test.jpg"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			imageID:        "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when image not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			imageID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.imageID
			if idStr == "" {
				idStr = imageID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/images/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetImageByShortID(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		shortID        string
		wantStatusCode int
	}{
		{
			name: "gets image by short ID successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.imagesByShortID["abc123"] = Image{ShortID: "abc123"}
			},
			shortID:        "abc123",
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when image not found",
			setupRepo: func(m *mockServiceRepo) {
				m.getImageErr = fmt.Errorf("not found")
			},
			shortID:        "missing",
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/images/short/"+tt.shortID, nil)
			req.SetPathValue("short_id", tt.shortID)
			w := httptest.NewRecorder()

			handler.GetImageByShortID(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetImageByShortID() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerListImages(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "lists images successfully",
			setupRepo: func(m *mockServiceRepo) {
				m.images[uuid.New()] = Image{FileName: "img1.jpg"}
				m.images[uuid.New()] = Image{FileName: "img2.jpg"}
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "returns empty list when no images",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) {
				m.getImageErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodGet, "/ssg/images", nil)
			w := httptest.NewRecorder()

			handler.ListImages(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListImages() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateImage(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name: "updates image successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.images[id] = Image{ID: id, FileName: "old.jpg"}
				return id
			},
			requestBody: map[string]string{
				"fileName": "new.jpg",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateImageErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]string{
				"fileName": "test.jpg",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			imageID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/images/"+imageID.String(), bytes.NewReader(body))
			req.SetPathValue("id", imageID.String())
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteImage(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "deletes image successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.images[id] = Image{ID: id, FileName: "test.jpg"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteImageErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			imageID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodDelete, "/ssg/images/"+imageID.String(), nil)
			req.SetPathValue("id", imageID.String())
			w := httptest.NewRecorder()

			handler.DeleteImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerCreateImageVariant(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupRepo      func(*mockServiceRepo)
		wantStatusCode int
	}{
		{
			name: "creates image variant successfully",
			requestBody: map[string]interface{}{
				"imageID": uuid.New().String(),
				"kind":    "thumbnail",
				"width":   200,
			},
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:           "fails with invalid JSON",
			requestBody:    "invalid json",
			setupRepo:      func(m *mockServiceRepo) {},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			requestBody: map[string]string{
				"kind": "thumbnail",
			},
			setupRepo: func(m *mockServiceRepo) {
				m.createImageVariantErr = fmt.Errorf("db error")
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

			req := httptest.NewRequest(http.MethodPost, "/ssg/image-variants", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerGetImageVariant(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		variantID      string
		wantStatusCode int
	}{
		{
			name: "gets image variant successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.imageVariants[id] = ImageVariant{ID: id, Kind: "thumbnail"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			variantID:      "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when variant not found",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			variantID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.variantID
			if idStr == "" {
				idStr = variantID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/image-variants/"+idStr, nil)
			req.SetPathValue("id", idStr)
			w := httptest.NewRecorder()

			handler.GetImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("GetImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerListImageVariantsByImageID(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		imageID        string
		wantStatusCode int
	}{
		{
			name: "lists variants successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				imageID := uuid.New()
				// No need to populate variants for basic test
				return imageID
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with invalid UUID",
			imageID:        "invalid-uuid",
			setupRepo:      func(m *mockServiceRepo) uuid.UUID { return uuid.Nil },
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.getImageVariantErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			imageID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			idStr := tt.imageID
			if idStr == "" {
				idStr = imageID.String()
			}

			req := httptest.NewRequest(http.MethodGet, "/ssg/images/"+idStr+"/variants", nil)
			req.SetPathValue("image_id", idStr)
			w := httptest.NewRecorder()

			handler.ListImageVariantsByImageID(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListImageVariantsByImageID() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerUpdateImageVariant(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		requestBody    interface{}
		wantStatusCode int
	}{
		{
			name: "updates image variant successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.imageVariants[id] = ImageVariant{ID: id, Kind: "old"}
				return id
			},
			requestBody: map[string]string{
				"kind": "new",
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.updateImageVariantErr = fmt.Errorf("db error")
				return uuid.New()
			},
			requestBody: map[string]string{
				"kind": "test",
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			variantID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			body, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPut, "/ssg/image-variants/"+variantID.String(), bytes.NewReader(body))
			req.SetPathValue("id", variantID.String())
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestAPIHandlerDeleteImageVariant(t *testing.T) {
	tests := []struct {
		name           string
		setupRepo      func(*mockServiceRepo) uuid.UUID
		wantStatusCode int
	}{
		{
			name: "deletes image variant successfully",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				id := uuid.New()
				m.imageVariants[id] = ImageVariant{ID: id, Kind: "test"}
				return id
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "fails when service returns error",
			setupRepo: func(m *mockServiceRepo) uuid.UUID {
				m.deleteImageVariantErr = fmt.Errorf("db error")
				return uuid.New()
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newMockServiceRepo()
			variantID := tt.setupRepo(repo)
			svc := newTestService(repo)

			cfg := hm.NewConfig()
			handler := NewAPIHandler("test-api", svc, nil, hm.XParams{Cfg: cfg})

			req := httptest.NewRequest(http.MethodDelete, "/ssg/image-variants/"+variantID.String(), nil)
			req.SetPathValue("id", variantID.String())
			w := httptest.NewRecorder()

			handler.DeleteImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
