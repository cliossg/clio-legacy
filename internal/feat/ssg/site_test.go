package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewSite(t *testing.T) {
	tests := []struct {
		name string
		sName string
		slug  string
		mode  string
		want  Site
	}{
		{
			name:  "creates site with all fields",
			sName: "Test Site",
			slug:  "test-site",
			mode:  "blog",
			want: Site{
				Name:      "Test Site",
				SlugValue: "test-site",
				Mode:      "blog",
				Active:    1,
			},
		},
		{
			name:  "creates site with structured mode",
			sName: "My Portfolio",
			slug:  "portfolio",
			mode:  "structured",
			want: Site{
				Name:      "My Portfolio",
				SlugValue: "portfolio",
				Mode:      "structured",
				Active:    1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSite(tt.sName, tt.slug, tt.mode)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.SlugValue != tt.want.SlugValue {
				t.Errorf("SlugValue = %v, want %v", got.SlugValue, tt.want.SlugValue)
			}
			if got.Mode != tt.want.Mode {
				t.Errorf("Mode = %v, want %v", got.Mode, tt.want.Mode)
			}
			if got.Active != tt.want.Active {
				t.Errorf("Active = %v, want %v", got.Active, tt.want.Active)
			}
		})
	}
}

func TestSiteType(t *testing.T) {
	s := Site{}
	got := s.Type()
	want := "site"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestSiteGetID(t *testing.T) {
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
			s := Site{ID: tt.id}
			got := s.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestSiteGenID(t *testing.T) {
	s := Site{}
	s.GenID()

	if s.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestSiteSetID(t *testing.T) {
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
			s := Site{ID: tt.initial}
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

func TestSiteGetShortID(t *testing.T) {
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
			s := Site{ShortID: tt.shortID}
			got := s.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestSiteGenShortID(t *testing.T) {
	s := Site{}
	s.GenShortID()

	if s.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(s.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(s.ShortID))
	}
}

func TestSiteSetShortID(t *testing.T) {
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
			s := Site{ShortID: tt.initial}
			s.SetShortID(tt.newID, tt.force...)

			if s.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", s.ShortID, tt.expected)
			}
		})
	}
}

func TestSiteSlug(t *testing.T) {
	tests := []struct {
		name      string
		slugValue string
		want      string
	}{
		{
			name:      "returns slug value",
			slugValue: "my-site",
			want:      "my-site",
		},
		{
			name:      "returns empty string when not set",
			slugValue: "",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Site{SlugValue: tt.slugValue}
			got := s.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSiteGenCreateValues(t *testing.T) {
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
			s := Site{}
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

func TestSiteGenUpdateValues(t *testing.T) {
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
			s := Site{}
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

func TestSiteGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	s := Site{CreatedBy: userID}
	got := s.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestSiteGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	s := Site{UpdatedBy: userID}
	got := s.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestSiteGetCreatedAt(t *testing.T) {
	now := time.Now()
	s := Site{CreatedAt: now}
	got := s.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestSiteGetUpdatedAt(t *testing.T) {
	now := time.Now()
	s := Site{UpdatedAt: now}
	got := s.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestSiteSetCreatedAt(t *testing.T) {
	now := time.Now()
	s := Site{}
	s.SetCreatedAt(now)

	if s.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", s.CreatedAt, now)
	}
}

func TestSiteSetUpdatedAt(t *testing.T) {
	now := time.Now()
	s := Site{}
	s.SetUpdatedAt(now)

	if s.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", s.UpdatedAt, now)
	}
}

func TestSiteSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	s := Site{}
	s.SetCreatedBy(userID)

	if s.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", s.CreatedBy, userID)
	}
}

func TestSiteSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	s := Site{}
	s.SetUpdatedBy(userID)

	if s.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", s.UpdatedBy, userID)
	}
}

func TestSiteIsZero(t *testing.T) {
	tests := []struct {
		name string
		s    Site
		want bool
	}{
		{
			name: "returns true for uninitialized site",
			s:    Site{},
			want: true,
		},
		{
			name: "returns false for initialized site",
			s:    Site{ID: uuid.New()},
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

func TestSiteRef(t *testing.T) {
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
			s := Site{ref: tt.ref}
			got := s.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestSiteSetRef(t *testing.T) {
	s := Site{}
	ref := "new-ref"
	s.SetRef(ref)

	if s.ref != ref {
		t.Errorf("SetRef() set %v, want %v", s.ref, ref)
	}
}

func TestSiteOptValue(t *testing.T) {
	id := uuid.New()
	s := Site{ID: id}
	got := s.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestSiteOptLabel(t *testing.T) {
	s := Site{Name: "My Site"}
	got := s.OptLabel()
	want := "My Site"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}
