package ssg

import (
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewImageVariant(t *testing.T) {
	imageID := uuid.New()
	variant := NewImageVariant(imageID, "thumbnail", "/path/test.jpg", "http://test.jpg", "image/jpeg", 1024, 150, 150)

	if variant.ImageID != imageID {
		t.Errorf("ImageID = %v, want %v", variant.ImageID, imageID)
	}
	if variant.Name != "thumbnail" {
		t.Errorf("Name = %v, want thumbnail", variant.Name)
	}
	if variant.Width != 150 {
		t.Errorf("Width = %v, want 150", variant.Width)
	}
}

func TestImageVariantType(t *testing.T) {
	variant := &ImageVariant{}
	if variant.Type() != "imageVariant" {
		t.Errorf("Type() = %v, want imageVariant", variant.Type())
	}
}

func TestImageVariantGetID(t *testing.T) {
	id := uuid.New()
	variant := &ImageVariant{ID: id}

	if variant.GetID() != id {
		t.Errorf("GetID() = %v, want %v", variant.GetID(), id)
	}
}

func TestImageVariantGenID(t *testing.T) {
	variant := &ImageVariant{}
	variant.GenID()

	if variant.ID == uuid.Nil {
		t.Error("GenID() did not generate an ID")
	}
}

func TestImageVariantSetID(t *testing.T) {
	tests := []struct {
		name      string
		existing  uuid.UUID
		newID     uuid.UUID
		force     []bool
		expectSet bool
	}{
		{
			name:      "sets ID when nil",
			existing:  uuid.Nil,
			newID:     uuid.New(),
			force:     nil,
			expectSet: true,
		},
		{
			name:      "does not replace existing ID without force",
			existing:  uuid.New(),
			newID:     uuid.New(),
			force:     nil,
			expectSet: false,
		},
		{
			name:      "replaces existing ID with force",
			existing:  uuid.New(),
			newID:     uuid.New(),
			force:     []bool{true},
			expectSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variant := &ImageVariant{ID: tt.existing}
			variant.SetID(tt.newID, tt.force...)

			if tt.expectSet {
				if variant.ID != tt.newID {
					t.Errorf("SetID() ID = %v, want %v", variant.ID, tt.newID)
				}
			} else {
				if variant.ID != tt.existing {
					t.Errorf("SetID() unexpectedly changed ID from %v to %v", tt.existing, variant.ID)
				}
			}
		})
	}
}

func TestImageVariantGetShortID(t *testing.T) {
	variant := &ImageVariant{ShortID: "abc123"}

	if variant.GetShortID() != "abc123" {
		t.Errorf("GetShortID() = %v, want abc123", variant.GetShortID())
	}
}

func TestImageVariantGenShortID(t *testing.T) {
	variant := &ImageVariant{}
	variant.GenShortID()

	if variant.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}
}

func TestImageVariantSetShortID(t *testing.T) {
	tests := []struct {
		name      string
		existing  string
		newID     string
		force     []bool
		expectSet bool
	}{
		{
			name:      "sets short ID when empty",
			existing:  "",
			newID:     "abc123",
			force:     nil,
			expectSet: true,
		},
		{
			name:      "does not replace existing short ID without force",
			existing:  "xyz789",
			newID:     "abc123",
			force:     nil,
			expectSet: false,
		},
		{
			name:      "replaces existing short ID with force",
			existing:  "xyz789",
			newID:     "abc123",
			force:     []bool{true},
			expectSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variant := &ImageVariant{ShortID: tt.existing}
			variant.SetShortID(tt.newID, tt.force...)

			if tt.expectSet {
				if variant.ShortID != tt.newID {
					t.Errorf("SetShortID() ShortID = %v, want %v", variant.ShortID, tt.newID)
				}
			} else {
				if variant.ShortID != tt.existing {
					t.Errorf("SetShortID() unexpectedly changed ShortID from %v to %v", tt.existing, variant.ShortID)
				}
			}
		})
	}
}

func TestImageVariantTypeID(t *testing.T) {
	variant := &ImageVariant{ShortID: "abc123"}
	typeID := variant.TypeID()

	if typeID != "imagevariant-abc123" {
		t.Errorf("TypeID() = %v, want imagevariant-abc123", typeID)
	}
}

func TestImageVariantIsZero(t *testing.T) {
	tests := []struct {
		name     string
		variant  ImageVariant
		expected bool
	}{
		{
			name:     "returns true for zero value",
			variant:  ImageVariant{},
			expected: true,
		},
		{
			name:     "returns false for initialized value",
			variant:  ImageVariant{ID: uuid.New()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.variant.IsZero(); got != tt.expected {
				t.Errorf("IsZero() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestImageVariantSlug(t *testing.T) {
	variant := &ImageVariant{
		Name:    "Thumbnail Image",
		ShortID: "abc123",
	}

	slug := variant.Slug()
	if slug != "thumbnail-image-abc123" {
		t.Errorf("Slug() = %v, want thumbnail-image-abc123", slug)
	}
}

func TestImageVariantOptValue(t *testing.T) {
	id := uuid.New()
	variant := &ImageVariant{ID: id}

	if variant.OptValue() != id.String() {
		t.Errorf("OptValue() = %v, want %v", variant.OptValue(), id.String())
	}
}

func TestImageVariantOptLabel(t *testing.T) {
	variant := &ImageVariant{Name: "Thumbnail"}

	if variant.OptLabel() != "Thumbnail" {
		t.Errorf("OptLabel() = %v, want Thumbnail", variant.OptLabel())
	}
}

func TestImageVariantStringID(t *testing.T) {
	id := uuid.New()
	variant := &ImageVariant{ID: id}

	if variant.StringID() != id.String() {
		t.Errorf("StringID() = %v, want %v", variant.StringID(), id.String())
	}
}

func TestToWebImageVariant(t *testing.T) {
	featVariant := feat.ImageVariant{
		ID:            uuid.New(),
		ShortID:       "abc123",
		ImageID:       uuid.New(),
		Kind:          "thumbnail",
		BlobRef:       "/path/test.jpg",
		Mime:          "image/jpeg",
		FilesizeByte:  1024,
		Width:         150,
		Height:        150,
	}

	webVariant := ToWebImageVariant(featVariant)

	if webVariant.ID != featVariant.ID {
		t.Errorf("ID = %v, want %v", webVariant.ID, featVariant.ID)
	}
	if webVariant.Name != featVariant.Kind {
		t.Errorf("Name = %v, want %v", webVariant.Name, featVariant.Kind)
	}
	if webVariant.Path != featVariant.BlobRef {
		t.Errorf("Path = %v, want %v", webVariant.Path, featVariant.BlobRef)
	}
	if webVariant.Width != 150 {
		t.Errorf("Width = %v, want 150", webVariant.Width)
	}
}

func TestToWebImageVariants(t *testing.T) {
	featVariants := []feat.ImageVariant{
		{
			ID:      uuid.New(),
			ShortID: "abc1",
			Kind:    "thumbnail",
		},
		{
			ID:      uuid.New(),
			ShortID: "abc2",
			Kind:    "medium",
		},
	}

	webVariants := ToWebImageVariants(featVariants)

	if len(webVariants) != 2 {
		t.Errorf("len(webVariants) = %v, want 2", len(webVariants))
		return
	}
	if webVariants[0].Name != "thumbnail" {
		t.Errorf("webVariants[0].Name = %v, want thumbnail", webVariants[0].Name)
	}
	if webVariants[1].Name != "medium" {
		t.Errorf("webVariants[1].Name = %v, want medium", webVariants[1].Name)
	}
}

func TestNewImageVariantForm(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	form := NewImageVariantForm(req)

	if form.BaseForm == nil {
		t.Error("NewImageVariantForm() BaseForm is nil")
	}
}

func TestToFeatImageVariant(t *testing.T) {
	id := uuid.New()
	imageID := uuid.New()

	form := ImageVariantForm{
		ID:      id.String(),
		ImageID: imageID.String(),
		Name:    "thumbnail",
	}

	featVariant := ToFeatImageVariant(form)

	if featVariant.ID != id {
		t.Errorf("ID = %v, want %v", featVariant.ID, id)
	}
	if featVariant.ImageID != imageID {
		t.Errorf("ImageID = %v, want %v", featVariant.ImageID, imageID)
	}
	if featVariant.Kind != "thumbnail" {
		t.Errorf("Kind = %v, want thumbnail", featVariant.Kind)
	}
}

func TestToImageVariantForm(t *testing.T) {
	id := uuid.New()
	imageID := uuid.New()

	variant := ImageVariant{
		ID:      id,
		ImageID: imageID,
		Name:    "thumbnail",
	}

	req := httptest.NewRequest("GET", "/", nil)
	form := ToImageVariantForm(req, variant)

	if form.ID != id.String() {
		t.Errorf("ID = %v, want %v", form.ID, id.String())
	}
	if form.ImageID != imageID.String() {
		t.Errorf("ImageID = %v, want %v", form.ImageID, imageID.String())
	}
	if form.Name != "thumbnail" {
		t.Errorf("Name = %v, want thumbnail", form.Name)
	}
}

func TestImageVariantFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      ImageVariantForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: ImageVariantForm{
				ImageID: uuid.New().String(),
				Name:    "thumbnail",
			},
			wantValid: true,
		},
		{
			name: "missing image ID",
			form: ImageVariantForm{
				Name: "thumbnail",
			},
			wantValid: false,
		},
		{
			name: "missing name",
			form: ImageVariantForm{
				ImageID: uuid.New().String(),
			},
			wantValid: false,
		},
		{
			name:      "missing both",
			form:      ImageVariantForm{},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewImageVariantForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}
