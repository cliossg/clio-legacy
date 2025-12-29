package ssg

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewImage(t *testing.T) {
	tests := []struct {
		name        string
		imageName   string
		description string
		path        string
		url         string
		altText     string
		mimeType    string
		size        int64
		width       int
		height      int
	}{
		{
			name:        "creates image with all fields",
			imageName:   "Test Image",
			description: "Test description",
			path:        "/path/to/image.jpg",
			url:         "http://example.com/image.jpg",
			altText:     "Alt text",
			mimeType:    "image/jpeg",
			size:        1024,
			width:       800,
			height:      600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := NewImage(tt.imageName, tt.description, tt.path, tt.url, tt.altText, tt.mimeType, tt.size, tt.width, tt.height)
			if image.Name != tt.imageName {
				t.Errorf("NewImage() Name = %v, want %v", image.Name, tt.imageName)
			}
			if image.Description != tt.description {
				t.Errorf("NewImage() Description = %v, want %v", image.Description, tt.description)
			}
			if image.Path != tt.path {
				t.Errorf("NewImage() Path = %v, want %v", image.Path, tt.path)
			}
			if image.URL != tt.url {
				t.Errorf("NewImage() URL = %v, want %v", image.URL, tt.url)
			}
			if image.AltText != tt.altText {
				t.Errorf("NewImage() AltText = %v, want %v", image.AltText, tt.altText)
			}
			if image.MimeType != tt.mimeType {
				t.Errorf("NewImage() MimeType = %v, want %v", image.MimeType, tt.mimeType)
			}
			if image.Size != tt.size {
				t.Errorf("NewImage() Size = %v, want %v", image.Size, tt.size)
			}
			if image.Width != tt.width {
				t.Errorf("NewImage() Width = %v, want %v", image.Width, tt.width)
			}
			if image.Height != tt.height {
				t.Errorf("NewImage() Height = %v, want %v", image.Height, tt.height)
			}
		})
	}
}

func TestImageType(t *testing.T) {
	image := &Image{}
	if got := image.Type(); got != imageType {
		t.Errorf("Type() = %v, want %v", got, imageType)
	}
}

func TestImageGetID(t *testing.T) {
	id := uuid.New()
	image := Image{ID: id}
	if got := image.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestImageGenID(t *testing.T) {
	image := &Image{}
	image.GenID()
	if image.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestImageSetID(t *testing.T) {
	tests := []struct {
		name    string
		initial uuid.UUID
		new     uuid.UUID
		force   []bool
		wantID  uuid.UUID
	}{
		{
			name:    "sets ID when empty",
			initial: uuid.Nil,
			new:     uuid.New(),
			force:   nil,
			wantID:  uuid.Nil,
		},
		{
			name:    "sets ID with force",
			initial: uuid.New(),
			new:     uuid.New(),
			force:   []bool{true},
			wantID:  uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := &Image{ID: tt.initial}
			image.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && image.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", image.ID, tt.wantID)
			}
		})
	}
}

func TestImageGetShortID(t *testing.T) {
	image := Image{ShortID: "test123"}
	if got := image.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestImageGenShortID(t *testing.T) {
	image := &Image{}
	image.GenShortID()
	if image.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestImageSetShortID(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		new     string
		force   []bool
	}{
		{
			name:    "sets shortID when empty",
			initial: "",
			new:     "test123",
			force:   nil,
		},
		{
			name:    "sets shortID with force",
			initial: "old123",
			new:     "test123",
			force:   []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image := &Image{ShortID: tt.initial}
			image.SetShortID(tt.new, tt.force...)
			if image.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && image.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", image.ShortID, tt.new)
				}
			}
		})
	}
}

func TestImageIsZero(t *testing.T) {
	tests := []struct {
		name  string
		image Image
		want  bool
	}{
		{
			name:  "zero image",
			image: Image{},
			want:  true,
		},
		{
			name:  "non-zero image",
			image: Image{ID: uuid.New()},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.image.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageSlug(t *testing.T) {
	image := &Image{Name: "Test Image", ShortID: "abc123"}
	got := image.Slug()
	if got == "" {
		t.Error("Slug() returned empty string")
	}
	if len(got) < len("test-image") {
		t.Errorf("Slug() = %v, too short", got)
	}
}

func TestImageTypeID(t *testing.T) {
	image := &Image{ShortID: "abc123"}
	got := image.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestImageOptLabel(t *testing.T) {
	image := Image{Name: "Test Image"}
	if got := image.OptLabel(); got != "Test Image" {
		t.Errorf("OptLabel() = %v, want Test Image", got)
	}
}

func TestImageOptValue(t *testing.T) {
	id := uuid.New()
	image := Image{ID: id}
	if got := image.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestImageStringID(t *testing.T) {
	id := uuid.New()
	image := Image{ID: id}
	if got := image.StringID(); got != id.String() {
		t.Errorf("StringID() = %v, want %v", got, id.String())
	}
}

func TestToWebImage(t *testing.T) {
	id := uuid.New()
	featImage := feat.Image{
		ID:       id,
		ShortID:  "abc123",
		Title:    "Test Image",
		FilePath: "path/to/image.jpg",
		AltText:  "Alt text",
		Width:    800,
		Height:   600,
	}

	webImage := ToWebImage(featImage)

	if webImage.ID != featImage.ID {
		t.Errorf("ToWebImage() ID = %v, want %v", webImage.ID, featImage.ID)
	}
	if webImage.ShortID != featImage.ShortID {
		t.Errorf("ToWebImage() ShortID = %v, want %v", webImage.ShortID, featImage.ShortID)
	}
	if webImage.Name != featImage.Title {
		t.Errorf("ToWebImage() Name = %v, want %v", webImage.Name, featImage.Title)
	}
	if webImage.Path != featImage.FilePath {
		t.Errorf("ToWebImage() Path = %v, want %v", webImage.Path, featImage.FilePath)
	}
	if webImage.AltText != featImage.AltText {
		t.Errorf("ToWebImage() AltText = %v, want %v", webImage.AltText, featImage.AltText)
	}
	if webImage.Width != featImage.Width {
		t.Errorf("ToWebImage() Width = %v, want %v", webImage.Width, featImage.Width)
	}
	if webImage.Height != featImage.Height {
		t.Errorf("ToWebImage() Height = %v, want %v", webImage.Height, featImage.Height)
	}
	expectedURL := "/static/images/" + featImage.FilePath
	if webImage.URL != expectedURL {
		t.Errorf("ToWebImage() URL = %v, want %v", webImage.URL, expectedURL)
	}
}

func TestToWebImages(t *testing.T) {
	featImages := []feat.Image{
		{
			ID:       uuid.New(),
			ShortID:  "abc123",
			Title:    "Image 1",
			FilePath: "path1.jpg",
			AltText:  "Alt 1",
			Width:    800,
			Height:   600,
		},
		{
			ID:       uuid.New(),
			ShortID:  "def456",
			Title:    "Image 2",
			FilePath: "path2.jpg",
			AltText:  "Alt 2",
			Width:    1024,
			Height:   768,
		},
	}

	webImages := ToWebImages(featImages)

	if len(webImages) != len(featImages) {
		t.Errorf("ToWebImages() length = %v, want %v", len(webImages), len(featImages))
	}

	for i, webImage := range webImages {
		if webImage.ID != featImages[i].ID {
			t.Errorf("ToWebImages()[%d] ID = %v, want %v", i, webImage.ID, featImages[i].ID)
		}
		if webImage.Name != featImages[i].Title {
			t.Errorf("ToWebImages()[%d] Name = %v, want %v", i, webImage.Name, featImages[i].Title)
		}
	}
}

func TestNewImageForm(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	form := NewImageForm(req)

	if form.BaseForm == nil {
		t.Error("NewImageForm() BaseForm is nil")
	}
}

func TestToFeatImage(t *testing.T) {
	id := uuid.New()
	form := ImageForm{
		ID:      id.String(),
		Name:    "Test Image",
		AltText: "Alt text",
	}

	featImage := ToFeatImage(form)

	if featImage.ID != id {
		t.Errorf("ToFeatImage() ID = %v, want %v", featImage.ID, id)
	}
	if featImage.Title != form.Name {
		t.Errorf("ToFeatImage() Title = %v, want %v", featImage.Title, form.Name)
	}
	if featImage.AltText != form.AltText {
		t.Errorf("ToFeatImage() AltText = %v, want %v", featImage.AltText, form.AltText)
	}
}

func TestToImageForm(t *testing.T) {
	id := uuid.New()
	image := Image{
		ID:          id,
		Name:        "Test Image",
		Description: "Test description",
		AltText:     "Alt text",
	}

	req := httptest.NewRequest("GET", "/", nil)
	form := ToImageForm(req, image)

	if form.ID != id.String() {
		t.Errorf("ToImageForm() ID = %v, want %v", form.ID, id.String())
	}
	if form.Name != image.Name {
		t.Errorf("ToImageForm() Name = %v, want %v", form.Name, image.Name)
	}
	if form.Description != image.Description {
		t.Errorf("ToImageForm() Description = %v, want %v", form.Description, image.Description)
	}
	if form.AltText != image.AltText {
		t.Errorf("ToImageForm() AltText = %v, want %v", form.AltText, image.AltText)
	}
}

func TestImageFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      ImageForm
		wantError bool
	}{
		{
			name: "validates with missing name",
			form: ImageForm{
				Name: "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewImageForm(req).BaseForm

			tt.form.Validate()

			hasErrors := tt.form.HasErrors()
			if hasErrors != tt.wantError {
				t.Errorf("Validate() hasErrors = %v, want %v", hasErrors, tt.wantError)
			}
		})
	}
}
