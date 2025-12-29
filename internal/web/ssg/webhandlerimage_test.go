package ssg

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestWebHandlerListImages(t *testing.T) {
	tests := []struct {
		name           string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name: "lists images successfully",
			getResp: map[string]interface{}{
				"images": []feat.Image{
					{FileName: "image1.jpg"},
					{FileName: "image2.jpg"},
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails when API returns error",
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-images", nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListImages(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListImages() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateImage(t *testing.T) {
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates image successfully",
			formData: url.Values{
				"name": []string{"Test Image"},
				"path": []string{"/images/test.jpg"},
			},
			postResp: map[string]interface{}{
				"image": map[string]interface{}{
					"id":   uuid.New().String(),
					"name": "Test Image",
					"path": "/images/test.jpg",
				},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"name": []string{"Test Image"},
				"path": []string{"/images/test.jpg"},
			},
			postErr:        fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, tt.postResp, tt.postErr, nil, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-image", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerEditImage(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows edit form successfully",
			queryID: imageID.String(),
			getResp: map[string]interface{}{
				"image": map[string]interface{}{
					"id":   imageID.String(),
					"name": "Test Image",
					"path": "/images/test.jpg",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        imageID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-image?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateImage(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates image successfully",
			formData: url.Values{
				"id":   []string{imageID.String()},
				"name": []string{"Updated Image"},
				"path": []string{"/images/updated.jpg"},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":   []string{imageID.String()},
				"name": []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":   []string{imageID.String()},
				"name": []string{"Updated Image"},
				"path": []string{"/images/updated.jpg"},
			},
			putErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, tt.putErr, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-image", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteImage(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes image successfully",
			formData: url.Values{
				"id": []string{imageID.String()},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name:           "fails with missing ID",
			formData:       url.Values{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id": []string{imageID.String()},
			},
			deleteErr:      fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, tt.deleteErr)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-image", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewImage(t *testing.T) {
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-image", nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewImage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewImage() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerShowImage(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows image successfully",
			queryID: imageID.String(),
			getResp: map[string]interface{}{
				"image": map[string]interface{}{
					"id":   imageID.String(),
					"name": "Test Image",
					"path": "/images/test.jpg",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        imageID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-image?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowImage(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowImage() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerListImageVariants(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "lists image variants successfully",
			queryID: imageID.String(),
			getResp: map[string]interface{}{
				"variants": []feat.ImageVariant{
					{Kind: "thumbnail"},
					{Kind: "large"},
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        imageID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/list-image-variants?image_id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ListImageVariants(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ListImageVariants() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerCreateImageVariant(t *testing.T) {
	imageID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		postResp       interface{}
		postErr        error
		wantStatusCode int
	}{
		{
			name: "creates image variant successfully",
			formData: url.Values{
				"image_id": []string{imageID.String()},
				"kind":     []string{"thumbnail"},
				"path":     []string{"/images/thumb.jpg"},
				"width":    []string{"150"},
				"height":   []string{"150"},
			},
			postResp: map[string]interface{}{
				"variant": map[string]interface{}{
					"id":       uuid.New().String(),
					"image_id": imageID.String(),
					"kind":     "thumbnail",
					"path":     "/images/thumb.jpg",
					"width":    150,
					"height":   150,
				},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"image_id": []string{imageID.String()},
				"kind":     []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"image_id": []string{imageID.String()},
				"kind":     []string{"thumbnail"},
				"path":     []string{"/images/thumb.jpg"},
			},
			postErr:        fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, tt.postResp, tt.postErr, nil, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/create-image-variant", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.CreateImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("CreateImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerEditImageVariant(t *testing.T) {
	variantID := uuid.New()
	imageID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows edit form successfully",
			queryID: variantID.String(),
			getResp: map[string]interface{}{
				"variant": map[string]interface{}{
					"id":       variantID.String(),
					"image_id": imageID.String(),
					"kind":     "thumbnail",
					"path":     "/images/thumb.jpg",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        variantID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/edit-image-variant?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.EditImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("EditImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerUpdateImageVariant(t *testing.T) {
	variantID := uuid.New()
	imageID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		putErr         error
		wantStatusCode int
	}{
		{
			name: "updates image variant successfully",
			formData: url.Values{
				"id":       []string{variantID.String()},
				"image_id": []string{imageID.String()},
				"kind":     []string{"large"},
				"path":     []string{"/images/large.jpg"},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name: "fails with invalid form data",
			formData: url.Values{
				"id":       []string{variantID.String()},
				"image_id": []string{imageID.String()},
				"kind":     []string{""},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id":       []string{variantID.String()},
				"image_id": []string{imageID.String()},
				"kind":     []string{"large"},
				"path":     []string{"/images/large.jpg"},
			},
			putErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, tt.putErr, nil)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/update-image-variant", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.UpdateImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("UpdateImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerDeleteImageVariant(t *testing.T) {
	variantID := uuid.New()
	tests := []struct {
		name           string
		formData       url.Values
		deleteErr      error
		wantStatusCode int
	}{
		{
			name: "deletes image variant successfully",
			formData: url.Values{
				"id": []string{variantID.String()},
			},
			wantStatusCode: http.StatusSeeOther,
		},
		{
			name:           "fails with missing ID",
			formData:       url.Values{},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "fails when API returns error",
			formData: url.Values{
				"id": []string{variantID.String()},
			},
			deleteErr:      fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, tt.deleteErr)
			defer server.Close()

			body := strings.NewReader(tt.formData.Encode())
			req := httptest.NewRequest(http.MethodPost, "/ssg/delete-image-variant", body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.DeleteImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("DeleteImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestWebHandlerNewImageVariant(t *testing.T) {
	imageID := uuid.New()
	handler, server := newTestWebHandlerWithMockAPI(nil, nil, nil, nil, nil, nil)
	defer server.Close()

	req := httptest.NewRequest(http.MethodGet, "/ssg/new-image-variant?image_id="+imageID.String(), nil)
	ctx := feat.NewContextWithSite("test-site", uuid.New())
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.NewImageVariant(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("NewImageVariant() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestWebHandlerShowImageVariant(t *testing.T) {
	variantID := uuid.New()
	imageID := uuid.New()
	tests := []struct {
		name           string
		queryID        string
		getResp        interface{}
		getErr         error
		wantStatusCode int
	}{
		{
			name:    "shows image variant successfully",
			queryID: variantID.String(),
			getResp: map[string]interface{}{
				"variant": map[string]interface{}{
					"id":       variantID.String(),
					"image_id": imageID.String(),
					"kind":     "thumbnail",
					"path":     "/images/thumb.jpg",
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "fails with missing ID",
			queryID:        "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "fails when API returns error",
			queryID:        variantID.String(),
			getErr:         fmt.Errorf("api error"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, server := newTestWebHandlerWithMockAPI(tt.getResp, tt.getErr, nil, nil, nil, nil)
			defer server.Close()

			req := httptest.NewRequest(http.MethodGet, "/ssg/show-image-variant?id="+tt.queryID, nil)
			ctx := feat.NewContextWithSite("test-site", uuid.New())
			req = req.WithContext(ctx)
			w := httptest.NewRecorder()

			handler.ShowImageVariant(w, req)

			if w.Code != tt.wantStatusCode {
				t.Errorf("ShowImageVariant() status = %d, want %d", w.Code, tt.wantStatusCode)
			}
		})
	}
}
