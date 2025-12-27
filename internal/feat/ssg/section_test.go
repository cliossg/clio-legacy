package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewSection(t *testing.T) {
	tests := []struct {
		name        string
		sName       string
		description string
		path        string
		layoutID    uuid.UUID
		want        Section
	}{
		{
			name:        "creates section with all fields",
			sName:       "Blog Posts",
			description: "Main blog section",
			path:        "/blog",
			layoutID:    uuid.New(),
			want: Section{
				Name:        "Blog Posts",
				Description: "Main blog section",
				Path:        "/blog",
			},
		},
		{
			name:        "creates section with empty description",
			sName:       "Portfolio",
			description: "",
			path:        "/portfolio",
			layoutID:    uuid.New(),
			want: Section{
				Name:        "Portfolio",
				Description: "",
				Path:        "/portfolio",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSection(tt.sName, tt.description, tt.path, tt.layoutID)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %v, want %v", got.Description, tt.want.Description)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.LayoutID != tt.layoutID {
				t.Errorf("LayoutID = %v, want %v", got.LayoutID, tt.layoutID)
			}
		})
	}
}

func TestSectionType(t *testing.T) {
	s := Section{}
	got := s.Type()
	want := "section"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestSectionGetID(t *testing.T) {
	tests := []struct {
		name string
		id   uuid.UUID
	}{
		{
			name: "returns set ID",
			id:   uuid.New(),
		},
		{
			name: "returns nil UUID when not set",
			id:   uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{ID: tt.id}
			got := s.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestSectionGenID(t *testing.T) {
	s := Section{}
	s.GenID()

	if s.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestSectionSetID(t *testing.T) {
	tests := []struct {
		name     string
		initial  uuid.UUID
		newID    uuid.UUID
		force    []bool
		expected uuid.UUID
	}{
		{
			name:     "sets ID when nil",
			initial:  uuid.Nil,
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil,
		},
		{
			name:     "does not override existing ID without force",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil,
		},
		{
			name:     "overrides existing ID with force true",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{true},
			expected: uuid.Nil,
		},
		{
			name:     "does not override with force false",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{false},
			expected: uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{ID: tt.initial}
			s.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if s.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", s.ID, want)
			}
		})
	}
}

func TestSectionGetShortID(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
	}{
		{
			name:    "returns set short ID",
			shortID: "abc123",
		},
		{
			name:    "returns empty string when not set",
			shortID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{ShortID: tt.shortID}
			got := s.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestSectionGenShortID(t *testing.T) {
	s := Section{}
	s.GenShortID()

	if s.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(s.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(s.ShortID))
	}
}

func TestSectionSetShortID(t *testing.T) {
	tests := []struct {
		name     string
		initial  string
		newID    string
		force    []bool
		expected string
	}{
		{
			name:     "sets short ID when empty",
			initial:  "",
			newID:    "xyz789",
			force:    nil,
			expected: "xyz789",
		},
		{
			name:     "does not override existing short ID without force",
			initial:  "abc123",
			newID:    "xyz789",
			force:    nil,
			expected: "abc123",
		},
		{
			name:     "overrides existing short ID with force true",
			initial:  "abc123",
			newID:    "xyz789",
			force:    []bool{true},
			expected: "xyz789",
		},
		{
			name:     "does not override with force false",
			initial:  "abc123",
			newID:    "xyz789",
			force:    []bool{false},
			expected: "abc123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{ShortID: tt.initial}
			s.SetShortID(tt.newID, tt.force...)

			if s.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", s.ShortID, tt.expected)
			}
		})
	}
}

func TestSectionGenCreateValues(t *testing.T) {
	tests := []struct {
		name   string
		userID []uuid.UUID
	}{
		{
			name:   "sets create values with user ID",
			userID: []uuid.UUID{uuid.New()},
		},
		{
			name:   "sets create values without user ID",
			userID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{}
			beforeTime := time.Now()
			s.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if s.CreatedAt.Before(beforeTime) || s.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", s.CreatedAt)
			}

			if s.UpdatedAt.Before(beforeTime) || s.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", s.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if s.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", s.CreatedBy, tt.userID[0])
				}
				if s.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", s.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestSectionGenUpdateValues(t *testing.T) {
	tests := []struct {
		name   string
		userID []uuid.UUID
	}{
		{
			name:   "sets update values with user ID",
			userID: []uuid.UUID{uuid.New()},
		},
		{
			name:   "sets update values without user ID",
			userID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{}
			beforeTime := time.Now()
			s.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if s.UpdatedAt.Before(beforeTime) || s.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", s.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if s.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", s.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestSectionGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	s := Section{CreatedBy: userID}
	got := s.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestSectionGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	s := Section{UpdatedBy: userID}
	got := s.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestSectionGetCreatedAt(t *testing.T) {
	now := time.Now()
	s := Section{CreatedAt: now}
	got := s.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestSectionGetUpdatedAt(t *testing.T) {
	now := time.Now()
	s := Section{UpdatedAt: now}
	got := s.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestSectionSetCreatedAt(t *testing.T) {
	now := time.Now()
	s := Section{}
	s.SetCreatedAt(now)

	if s.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", s.CreatedAt, now)
	}
}

func TestSectionSetUpdatedAt(t *testing.T) {
	now := time.Now()
	s := Section{}
	s.SetUpdatedAt(now)

	if s.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", s.UpdatedAt, now)
	}
}

func TestSectionSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	s := Section{}
	s.SetCreatedBy(userID)

	if s.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", s.CreatedBy, userID)
	}
}

func TestSectionSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	s := Section{}
	s.SetUpdatedBy(userID)

	if s.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", s.UpdatedBy, userID)
	}
}

func TestSectionIsZero(t *testing.T) {
	tests := []struct {
		name string
		s    Section
		want bool
	}{
		{
			name: "returns true for uninitialized section",
			s:    Section{},
			want: true,
		},
		{
			name: "returns false for initialized section",
			s:    Section{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSectionSlug(t *testing.T) {
	tests := []struct {
		name    string
		sName   string
		shortID string
		want    string
	}{
		{
			name:    "generates slug from name and short ID",
			sName:   "Blog Posts",
			shortID: "abc123",
			want:    "blog-posts-abc123",
		},
		{
			name:    "handles special characters",
			sName:   "News & Updates",
			shortID: "xyz789",
			want:    "news-&-updates-xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{
				Name:    tt.sName,
				ShortID: tt.shortID,
			}
			got := s.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSectionOptValue(t *testing.T) {
	id := uuid.New()
	s := Section{ID: id}
	got := s.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestSectionOptLabel(t *testing.T) {
	s := Section{Name: "Blog Posts"}
	got := s.OptLabel()
	want := "Blog Posts"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}

func TestSectionRef(t *testing.T) {
	tests := []struct {
		name string
		ref  string
	}{
		{
			name: "returns set ref",
			ref:  "test-ref",
		},
		{
			name: "returns empty string when not set",
			ref:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Section{ref: tt.ref}
			got := s.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestSectionSetRef(t *testing.T) {
	s := Section{}
	ref := "new-ref"
	s.SetRef(ref)

	if s.ref != ref {
		t.Errorf("SetRef() set %v, want %v", s.ref, ref)
	}
}
