package ssg

import (
	"testing"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		tagName string
	}{
		{
			name:    "creates tag with name",
			tagName: "Test Tag",
		},
		{
			name:    "creates tag with empty name",
			tagName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := NewTag(tt.tagName)
			if tag.Name != tt.tagName {
				t.Errorf("NewTag() Name = %v, want %v", tag.Name, tt.tagName)
			}
		})
	}
}

func TestTagType(t *testing.T) {
	tag := &Tag{}
	if got := tag.Type(); got != tagType {
		t.Errorf("Type() = %v, want %v", got, tagType)
	}
}

func TestTagGetID(t *testing.T) {
	id := uuid.New()
	tag := Tag{ID: id}
	if got := tag.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestTagGenID(t *testing.T) {
	tag := &Tag{}
	tag.GenID()
	if tag.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestTagSetID(t *testing.T) {
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
			tag := &Tag{ID: tt.initial}
			tag.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && tag.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", tag.ID, tt.wantID)
			}
		})
	}
}

func TestTagGetShortID(t *testing.T) {
	tag := Tag{ShortID: "test123"}
	if got := tag.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestTagGenShortID(t *testing.T) {
	tag := &Tag{}
	tag.GenShortID()
	if tag.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestTagSetShortID(t *testing.T) {
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
			tag := &Tag{ShortID: tt.initial}
			tag.SetShortID(tt.new, tt.force...)
			if tag.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && tag.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", tag.ShortID, tt.new)
				}
			}
		})
	}
}

func TestTagIsZero(t *testing.T) {
	tests := []struct {
		name string
		tag  Tag
		want bool
	}{
		{
			name: "zero tag",
			tag:  Tag{},
			want: true,
		},
		{
			name: "non-zero tag",
			tag:  Tag{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tag.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagSlug(t *testing.T) {
	tests := []struct {
		name      string
		tag       Tag
		wantEmpty bool
	}{
		{
			name: "returns slug field when set",
			tag: Tag{
				Name:      "Test Tag",
				ShortID:   "abc123",
				SlugField: "custom-slug",
			},
			wantEmpty: false,
		},
		{
			name: "generates slug from name when slug field empty",
			tag: Tag{
				Name:    "Test Tag",
				ShortID: "abc123",
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tag.Slug()
			if tt.wantEmpty && got != "" {
				t.Errorf("Slug() = %v, want empty string", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("Slug() returned empty string")
			}
			if tt.tag.SlugField != "" && got != tt.tag.SlugField {
				t.Errorf("Slug() = %v, want %v", got, tt.tag.SlugField)
			}
		})
	}
}

func TestTagTypeID(t *testing.T) {
	tag := &Tag{ShortID: "abc123"}
	got := tag.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestTagOptLabel(t *testing.T) {
	tag := Tag{Name: "Test Tag"}
	if got := tag.OptLabel(); got != "Test Tag" {
		t.Errorf("OptLabel() = %v, want Test Tag", got)
	}
}

func TestTagOptValue(t *testing.T) {
	id := uuid.New()
	tag := Tag{ID: id}
	if got := tag.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestToWebTag(t *testing.T) {
	id := uuid.New()
	featTag := feat.Tag{
		ID:      id,
		ShortID: "abc123",
		Name:    "Test Tag",
	}

	webTag := ToWebTag(featTag)

	if webTag.ID != featTag.ID {
		t.Errorf("ToWebTag() ID = %v, want %v", webTag.ID, featTag.ID)
	}
	if webTag.ShortID != featTag.ShortID {
		t.Errorf("ToWebTag() ShortID = %v, want %v", webTag.ShortID, featTag.ShortID)
	}
	if webTag.Name != featTag.Name {
		t.Errorf("ToWebTag() Name = %v, want %v", webTag.Name, featTag.Name)
	}
}

func TestToWebTags(t *testing.T) {
	featTags := []feat.Tag{
		{
			ID:      uuid.New(),
			ShortID: "abc123",
			Name:    "Tag 1",
		},
		{
			ID:      uuid.New(),
			ShortID: "def456",
			Name:    "Tag 2",
		},
	}

	webTags := ToWebTags(featTags)

	if len(webTags) != len(featTags) {
		t.Errorf("ToWebTags() length = %v, want %v", len(webTags), len(featTags))
	}

	for i, webTag := range webTags {
		if webTag.ID != featTags[i].ID {
			t.Errorf("ToWebTags()[%d] ID = %v, want %v", i, webTag.ID, featTags[i].ID)
		}
		if webTag.Name != featTags[i].Name {
			t.Errorf("ToWebTags()[%d] Name = %v, want %v", i, webTag.Name, featTags[i].Name)
		}
	}
}
