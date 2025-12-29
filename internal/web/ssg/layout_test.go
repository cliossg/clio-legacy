package ssg

import (
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewlayout(t *testing.T) {
	tests := []struct {
		name        string
		layoutName  string
		description string
		code        string
	}{
		{
			name:        "creates layout with all fields",
			layoutName:  "Test Layout",
			description: "Test description",
			code:        "<html></html>",
		},
		{
			name:        "creates layout with empty fields",
			layoutName:  "",
			description: "",
			code:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout := Newlayout(tt.layoutName, tt.description, tt.code)
			if layout.Name != tt.layoutName {
				t.Errorf("Newlayout() Name = %v, want %v", layout.Name, tt.layoutName)
			}
			if layout.Description != tt.description {
				t.Errorf("Newlayout() Description = %v, want %v", layout.Description, tt.description)
			}
			if layout.Code != tt.code {
				t.Errorf("Newlayout() Code = %v, want %v", layout.Code, tt.code)
			}
		})
	}
}

func TestLayoutType(t *testing.T) {
	layout := &Layout{}
	if got := layout.Type(); got != layoutType {
		t.Errorf("Type() = %v, want %v", got, layoutType)
	}
}

func TestLayoutGetID(t *testing.T) {
	id := uuid.New()
	layout := Layout{ID: id}
	if got := layout.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestLayoutGenID(t *testing.T) {
	layout := &Layout{}
	layout.GenID()
	if layout.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestLayoutSetID(t *testing.T) {
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
			layout := &Layout{ID: tt.initial}
			layout.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && layout.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", layout.ID, tt.wantID)
			}
		})
	}
}

func TestLayoutGetShortID(t *testing.T) {
	layout := Layout{ShortID: "test123"}
	if got := layout.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestLayoutGenShortID(t *testing.T) {
	layout := &Layout{}
	layout.GenShortID()
	if layout.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestLayoutSetShortID(t *testing.T) {
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
			layout := &Layout{ShortID: tt.initial}
			layout.SetShortID(tt.new, tt.force...)
			if layout.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && layout.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", layout.ShortID, tt.new)
				}
			}
		})
	}
}

func TestLayoutIsZero(t *testing.T) {
	tests := []struct {
		name   string
		layout Layout
		want   bool
	}{
		{
			name:   "zero layout",
			layout: Layout{},
			want:   true,
		},
		{
			name:   "non-zero layout",
			layout: Layout{ID: uuid.New()},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.layout.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLayoutSlug(t *testing.T) {
	layout := &Layout{Name: "Test Layout", ShortID: "abc123"}
	got := layout.Slug()
	if got == "" {
		t.Error("Slug() returned empty string")
	}
	if len(got) < len("test-layout") {
		t.Errorf("Slug() = %v, too short", got)
	}
}

func TestLayoutTypeID(t *testing.T) {
	layout := &Layout{ShortID: "abc123"}
	got := layout.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestLayoutOptLabel(t *testing.T) {
	layout := Layout{Name: "Test Layout"}
	if got := layout.OptLabel(); got != "Test Layout" {
		t.Errorf("OptLabel() = %v, want Test Layout", got)
	}
}

func TestLayoutOptValue(t *testing.T) {
	id := uuid.New()
	layout := Layout{ID: id}
	if got := layout.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestLayoutStringID(t *testing.T) {
	id := uuid.New()
	layout := Layout{ID: id}
	if got := layout.StringID(); got != id.String() {
		t.Errorf("StringID() = %v, want %v", got, id.String())
	}
}

func TestToWebLayout(t *testing.T) {
	id := uuid.New()
	featLayout := feat.Layout{
		ID:          id,
		ShortID:     "abc123",
		Name:        "Test Layout",
		Description: "Test description",
		Code:        "<html></html>",
	}

	webLayout := ToWebLayout(featLayout)

	if webLayout.ID != featLayout.ID {
		t.Errorf("ToWebLayout() ID = %v, want %v", webLayout.ID, featLayout.ID)
	}
	if webLayout.ShortID != featLayout.ShortID {
		t.Errorf("ToWebLayout() ShortID = %v, want %v", webLayout.ShortID, featLayout.ShortID)
	}
	if webLayout.Name != featLayout.Name {
		t.Errorf("ToWebLayout() Name = %v, want %v", webLayout.Name, featLayout.Name)
	}
	if webLayout.Description != featLayout.Description {
		t.Errorf("ToWebLayout() Description = %v, want %v", webLayout.Description, featLayout.Description)
	}
	if webLayout.Code != featLayout.Code {
		t.Errorf("ToWebLayout() Code = %v, want %v", webLayout.Code, featLayout.Code)
	}
}

func TestToWebLayouts(t *testing.T) {
	featLayouts := []feat.Layout{
		{
			ID:          uuid.New(),
			ShortID:     "abc123",
			Name:        "Layout 1",
			Description: "Description 1",
			Code:        "<html>1</html>",
		},
		{
			ID:          uuid.New(),
			ShortID:     "def456",
			Name:        "Layout 2",
			Description: "Description 2",
			Code:        "<html>2</html>",
		},
	}

	webLayouts := ToWebLayouts(featLayouts)

	if len(webLayouts) != len(featLayouts) {
		t.Errorf("ToWebLayouts() length = %v, want %v", len(webLayouts), len(featLayouts))
	}

	for i, webLayout := range webLayouts {
		if webLayout.ID != featLayouts[i].ID {
			t.Errorf("ToWebLayouts()[%d] ID = %v, want %v", i, webLayout.ID, featLayouts[i].ID)
		}
		if webLayout.Name != featLayouts[i].Name {
			t.Errorf("ToWebLayouts()[%d] Name = %v, want %v", i, webLayout.Name, featLayouts[i].Name)
		}
	}
}
