package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewContentTag(t *testing.T) {
	tests := []struct {
		name      string
		contentID uuid.UUID
		tagID     uuid.UUID
	}{
		{
			name:      "creates content tag with valid IDs",
			contentID: uuid.New(),
			tagID:     uuid.New(),
		},
		{
			name:      "creates content tag with nil IDs",
			contentID: uuid.Nil,
			tagID:     uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewContentTag(tt.contentID, tt.tagID)

			if got.ContentID != tt.contentID {
				t.Errorf("ContentID = %v, want %v", got.ContentID, tt.contentID)
			}
			if got.TagID != tt.tagID {
				t.Errorf("TagID = %v, want %v", got.TagID, tt.tagID)
			}
		})
	}
}

func TestContentTagType(t *testing.T) {
	ct := ContentTag{}
	got := ct.Type()
	want := "content-tag"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestContentTagGetID(t *testing.T) {
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
			ct := ContentTag{ID: tt.id}
			got := ct.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestContentTagGenID(t *testing.T) {
	ct := ContentTag{}
	ct.GenID()

	if ct.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestContentTagSetID(t *testing.T) {
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
			ct := ContentTag{ID: tt.initial}
			ct.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if ct.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", ct.ID, want)
			}
		})
	}
}

func TestContentTagGetShortID(t *testing.T) {
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
			ct := ContentTag{ShortID: tt.shortID}
			got := ct.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestContentTagGenShortID(t *testing.T) {
	ct := ContentTag{}
	ct.GenShortID()

	if ct.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(ct.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(ct.ShortID))
	}
}

func TestContentTagSetShortID(t *testing.T) {
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
			ct := ContentTag{ShortID: tt.initial}
			ct.SetShortID(tt.newID, tt.force...)

			if ct.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", ct.ShortID, tt.expected)
			}
		})
	}
}

func TestContentTagGenCreateValues(t *testing.T) {
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
			ct := ContentTag{}
			beforeTime := time.Now()
			ct.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if ct.CreatedAt.Before(beforeTime) || ct.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", ct.CreatedAt)
			}

			if ct.UpdatedAt.Before(beforeTime) || ct.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", ct.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if ct.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", ct.CreatedBy, tt.userID[0])
				}
				if ct.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", ct.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestContentTagGenUpdateValues(t *testing.T) {
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
			ct := ContentTag{}
			beforeTime := time.Now()
			ct.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if ct.UpdatedAt.Before(beforeTime) || ct.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", ct.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if ct.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", ct.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestContentTagGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	ct := ContentTag{CreatedBy: userID}
	got := ct.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestContentTagGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	ct := ContentTag{UpdatedBy: userID}
	got := ct.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestContentTagGetCreatedAt(t *testing.T) {
	now := time.Now()
	ct := ContentTag{CreatedAt: now}
	got := ct.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestContentTagGetUpdatedAt(t *testing.T) {
	now := time.Now()
	ct := ContentTag{UpdatedAt: now}
	got := ct.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestContentTagSetCreatedAt(t *testing.T) {
	now := time.Now()
	ct := ContentTag{}
	ct.SetCreatedAt(now)

	if ct.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", ct.CreatedAt, now)
	}
}

func TestContentTagSetUpdatedAt(t *testing.T) {
	now := time.Now()
	ct := ContentTag{}
	ct.SetUpdatedAt(now)

	if ct.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", ct.UpdatedAt, now)
	}
}

func TestContentTagSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	ct := ContentTag{}
	ct.SetCreatedBy(userID)

	if ct.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", ct.CreatedBy, userID)
	}
}

func TestContentTagSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	ct := ContentTag{}
	ct.SetUpdatedBy(userID)

	if ct.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", ct.UpdatedBy, userID)
	}
}

func TestContentTagIsZero(t *testing.T) {
	tests := []struct {
		name string
		ct   ContentTag
		want bool
	}{
		{
			name: "returns true for uninitialized content tag",
			ct:   ContentTag{},
			want: true,
		},
		{
			name: "returns false for initialized content tag",
			ct:   ContentTag{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ct.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentTagSlug(t *testing.T) {
	tests := []struct {
		name    string
		shortID string
		want    string
	}{
		{
			name:    "generates slug from type and short ID",
			shortID: "abc123",
			want:    "content-tag-abc123",
		},
		{
			name:    "returns type with empty short ID",
			shortID: "",
			want:    "content-tag-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := ContentTag{ShortID: tt.shortID}
			got := ct.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentTagOptValue(t *testing.T) {
	id := uuid.New()
	ct := ContentTag{ID: id}
	got := ct.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestContentTagOptLabel(t *testing.T) {
	ct := ContentTag{ShortID: "abc123"}
	got := ct.OptLabel()
	want := "abc123"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}

func TestContentTagRef(t *testing.T) {
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
			ct := ContentTag{ref: tt.ref}
			got := ct.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestContentTagSetRef(t *testing.T) {
	ct := ContentTag{}
	ref := "new-ref"
	ct.SetRef(ref)

	if ct.ref != ref {
		t.Errorf("SetRef() set %v, want %v", ct.ref, ref)
	}
}
