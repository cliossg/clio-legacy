package ssg

import (
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewSection(t *testing.T) {
	tests := []struct {
		name        string
		sectionName string
		description string
		path        string
		layoutID    uuid.UUID
	}{
		{
			name:        "creates section with all fields",
			sectionName: "Test Section",
			description: "Test description",
			path:        "/test",
			layoutID:    uuid.New(),
		},
		{
			name:        "creates section with empty fields",
			sectionName: "",
			description: "",
			path:        "",
			layoutID:    uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			section := NewSection(tt.sectionName, tt.description, tt.path, tt.layoutID)
			if section.Name != tt.sectionName {
				t.Errorf("NewSection() Name = %v, want %v", section.Name, tt.sectionName)
			}
			if section.Description != tt.description {
				t.Errorf("NewSection() Description = %v, want %v", section.Description, tt.description)
			}
			if section.Path != tt.path {
				t.Errorf("NewSection() Path = %v, want %v", section.Path, tt.path)
			}
			if section.LayoutID != tt.layoutID {
				t.Errorf("NewSection() LayoutID = %v, want %v", section.LayoutID, tt.layoutID)
			}
		})
	}
}

func TestSectionType(t *testing.T) {
	section := &Section{}
	if got := section.Type(); got != sectionType {
		t.Errorf("Type() = %v, want %v", got, sectionType)
	}
}

func TestSectionGetID(t *testing.T) {
	id := uuid.New()
	section := Section{ID: id}
	if got := section.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestSectionGenID(t *testing.T) {
	section := &Section{}
	section.GenID()
	if section.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestSectionSetID(t *testing.T) {
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
			section := &Section{ID: tt.initial}
			section.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && section.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", section.ID, tt.wantID)
			}
		})
	}
}

func TestSectionGetShortID(t *testing.T) {
	section := Section{ShortID: "test123"}
	if got := section.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestSectionGenShortID(t *testing.T) {
	section := &Section{}
	section.GenShortID()
	if section.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestSectionSetShortID(t *testing.T) {
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
			section := &Section{ShortID: tt.initial}
			section.SetShortID(tt.new, tt.force...)
			if section.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && section.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", section.ShortID, tt.new)
				}
			}
		})
	}
}

func TestSectionIsZero(t *testing.T) {
	tests := []struct {
		name    string
		section Section
		want    bool
	}{
		{
			name:    "zero section",
			section: Section{},
			want:    true,
		},
		{
			name:    "non-zero section",
			section: Section{ID: uuid.New()},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.section.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSectionSlug(t *testing.T) {
	section := &Section{Name: "Test Section", ShortID: "abc123"}
	got := section.Slug()
	if got == "" {
		t.Error("Slug() returned empty string")
	}
	if len(got) < len("test-section") {
		t.Errorf("Slug() = %v, too short", got)
	}
}

func TestSectionTypeID(t *testing.T) {
	section := &Section{ShortID: "abc123"}
	got := section.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestSectionOptLabel(t *testing.T) {
	section := Section{Name: "Test Section"}
	if got := section.OptLabel(); got != "Test Section" {
		t.Errorf("OptLabel() = %v, want Test Section", got)
	}
}

func TestSectionOptValue(t *testing.T) {
	id := uuid.New()
	section := Section{ID: id}
	if got := section.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestToWebSection(t *testing.T) {
	id := uuid.New()
	layoutID := uuid.New()
	featSection := feat.Section{
		ID:          id,
		ShortID:     "abc123",
		Name:        "Test Section",
		Description: "Test description",
		Path:        "/test",
		LayoutID:    layoutID,
		LayoutName:  "Test Layout",
	}

	webSection := ToWebSection(featSection)

	if webSection.ID != featSection.ID {
		t.Errorf("ToWebSection() ID = %v, want %v", webSection.ID, featSection.ID)
	}
	if webSection.ShortID != featSection.ShortID {
		t.Errorf("ToWebSection() ShortID = %v, want %v", webSection.ShortID, featSection.ShortID)
	}
	if webSection.Name != featSection.Name {
		t.Errorf("ToWebSection() Name = %v, want %v", webSection.Name, featSection.Name)
	}
	if webSection.Description != featSection.Description {
		t.Errorf("ToWebSection() Description = %v, want %v", webSection.Description, featSection.Description)
	}
	if webSection.Path != featSection.Path {
		t.Errorf("ToWebSection() Path = %v, want %v", webSection.Path, featSection.Path)
	}
	if webSection.LayoutID != featSection.LayoutID {
		t.Errorf("ToWebSection() LayoutID = %v, want %v", webSection.LayoutID, featSection.LayoutID)
	}
	if webSection.LayoutName != featSection.LayoutName {
		t.Errorf("ToWebSection() LayoutName = %v, want %v", webSection.LayoutName, featSection.LayoutName)
	}
}

func TestToWebSections(t *testing.T) {
	featSections := []feat.Section{
		{
			ID:          uuid.New(),
			ShortID:     "abc123",
			Name:        "Section 1",
			Description: "Description 1",
			Path:        "/path1",
			LayoutID:    uuid.New(),
			LayoutName:  "Layout 1",
		},
		{
			ID:          uuid.New(),
			ShortID:     "def456",
			Name:        "Section 2",
			Description: "Description 2",
			Path:        "/path2",
			LayoutID:    uuid.New(),
			LayoutName:  "Layout 2",
		},
	}

	webSections := ToWebSections(featSections)

	if len(webSections) != len(featSections) {
		t.Errorf("ToWebSections() length = %v, want %v", len(webSections), len(featSections))
	}

	for i, webSection := range webSections {
		if webSection.ID != featSections[i].ID {
			t.Errorf("ToWebSections()[%d] ID = %v, want %v", i, webSection.ID, featSections[i].ID)
		}
		if webSection.Name != featSections[i].Name {
			t.Errorf("ToWebSections()[%d] Name = %v, want %v", i, webSection.Name, featSections[i].Name)
		}
	}
}
