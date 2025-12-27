package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewImageVariant(t *testing.T) {
	got := NewImageVariant()

	if got.ID == uuid.Nil {
		t.Error("NewImageVariant() did not generate a UUID")
	}
}

func TestImageVariantType(t *testing.T) {
	iv := ImageVariant{}
	got := iv.Type()
	want := "image-variant"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestImageVariantGetID(t *testing.T) {
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
			iv := ImageVariant{ID: tt.id}
			got := iv.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestImageVariantGenID(t *testing.T) {
	iv := ImageVariant{}
	iv.GenID()

	if iv.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestImageVariantSetID(t *testing.T) {
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
			iv := ImageVariant{ID: tt.initial}
			iv.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if iv.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", iv.ID, want)
			}
		})
	}
}

func TestImageVariantGetShortID(t *testing.T) {
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
			iv := ImageVariant{ShortID: tt.shortID}
			got := iv.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestImageVariantGenShortID(t *testing.T) {
	iv := ImageVariant{}
	iv.GenShortID()

	if iv.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(iv.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(iv.ShortID))
	}
}

func TestImageVariantSetShortID(t *testing.T) {
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
			iv := ImageVariant{ShortID: tt.initial}
			iv.SetShortID(tt.newID, tt.force...)

			if iv.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", iv.ShortID, tt.expected)
			}
		})
	}
}

func TestImageVariantGenCreateValues(t *testing.T) {
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
			iv := ImageVariant{}
			beforeTime := time.Now()
			iv.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if iv.CreatedAt.Before(beforeTime) || iv.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", iv.CreatedAt)
			}

			if iv.UpdatedAt.Before(beforeTime) || iv.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", iv.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if iv.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", iv.CreatedBy, tt.userID[0])
				}
				if iv.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", iv.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestImageVariantGenUpdateValues(t *testing.T) {
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
			iv := ImageVariant{}
			beforeTime := time.Now()
			iv.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if iv.UpdatedAt.Before(beforeTime) || iv.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", iv.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if iv.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", iv.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestImageVariantGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	iv := ImageVariant{CreatedBy: userID}
	got := iv.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestImageVariantGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	iv := ImageVariant{UpdatedBy: userID}
	got := iv.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestImageVariantGetCreatedAt(t *testing.T) {
	now := time.Now()
	iv := ImageVariant{CreatedAt: now}
	got := iv.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestImageVariantGetUpdatedAt(t *testing.T) {
	now := time.Now()
	iv := ImageVariant{UpdatedAt: now}
	got := iv.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestImageVariantSetCreatedAt(t *testing.T) {
	now := time.Now()
	iv := ImageVariant{}
	iv.SetCreatedAt(now)

	if iv.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", iv.CreatedAt, now)
	}
}

func TestImageVariantSetUpdatedAt(t *testing.T) {
	now := time.Now()
	iv := ImageVariant{}
	iv.SetUpdatedAt(now)

	if iv.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", iv.UpdatedAt, now)
	}
}

func TestImageVariantSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	iv := ImageVariant{}
	iv.SetCreatedBy(userID)

	if iv.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", iv.CreatedBy, userID)
	}
}

func TestImageVariantSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	iv := ImageVariant{}
	iv.SetUpdatedBy(userID)

	if iv.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", iv.UpdatedBy, userID)
	}
}

func TestImageVariantIsZero(t *testing.T) {
	tests := []struct {
		name string
		iv   ImageVariant
		want bool
	}{
		{
			name: "returns true for uninitialized image variant",
			iv:   ImageVariant{},
			want: true,
		},
		{
			name: "returns false for initialized image variant",
			iv:   ImageVariant{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.iv.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageVariantSlug(t *testing.T) {
	tests := []struct {
		name    string
		kind    string
		shortID string
		want    string
	}{
		{
			name:    "generates slug from kind and short ID",
			kind:    "thumbnail",
			shortID: "abc123",
			want:    "thumbnail-abc123",
		},
		{
			name:    "handles web kind",
			kind:    "web",
			shortID: "xyz789",
			want:    "web-xyz789",
		},
		{
			name:    "handles original kind",
			kind:    "original",
			shortID: "def456",
			want:    "original-def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iv := ImageVariant{
				Kind:    tt.kind,
				ShortID: tt.shortID,
			}
			got := iv.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImageVariantRef(t *testing.T) {
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
			iv := ImageVariant{ref: tt.ref}
			got := iv.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestImageVariantSetRef(t *testing.T) {
	iv := ImageVariant{}
	ref := "new-ref"
	iv.SetRef(ref)

	if iv.ref != ref {
		t.Errorf("SetRef() set %v, want %v", iv.ref, ref)
	}
}
