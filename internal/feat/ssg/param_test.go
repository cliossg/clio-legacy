package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewParam(t *testing.T) {
	tests := []struct {
		name  string
		pName string
		value string
		want  Param
	}{
		{
			name:  "creates param with name and value",
			pName: "site_title",
			value: "My Website",
			want: Param{
				Name:  "site_title",
				Value: "My Website",
			},
		},
		{
			name:  "creates param with empty value",
			pName: "description",
			value: "",
			want: Param{
				Name:  "description",
				Value: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewParam(tt.pName, tt.value)

			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Value != tt.want.Value {
				t.Errorf("Value = %v, want %v", got.Value, tt.want.Value)
			}
		})
	}
}

func TestParamType(t *testing.T) {
	p := Param{}
	got := p.Type()
	want := "param"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestParamGetID(t *testing.T) {
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
			p := Param{ID: tt.id}
			got := p.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestParamGenID(t *testing.T) {
	p := Param{}
	p.GenID()

	if p.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestParamSetID(t *testing.T) {
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
			p := Param{ID: tt.initial}
			p.SetID(tt.newID, tt.force...)

			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if p.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", p.ID, want)
			}
		})
	}
}

func TestParamGetShortID(t *testing.T) {
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
			p := Param{ShortID: tt.shortID}
			got := p.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestParamGenShortID(t *testing.T) {
	p := Param{}
	p.GenShortID()

	if p.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(p.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(p.ShortID))
	}
}

func TestParamSetShortID(t *testing.T) {
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
			p := Param{ShortID: tt.initial}
			p.SetShortID(tt.newID, tt.force...)

			if p.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", p.ShortID, tt.expected)
			}
		})
	}
}

func TestParamGenCreateValues(t *testing.T) {
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
			p := Param{}
			beforeTime := time.Now()
			p.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if p.CreatedAt.Before(beforeTime) || p.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", p.CreatedAt)
			}

			if p.UpdatedAt.Before(beforeTime) || p.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", p.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if p.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", p.CreatedBy, tt.userID[0])
				}
				if p.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", p.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestParamGenUpdateValues(t *testing.T) {
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
			p := Param{}
			beforeTime := time.Now()
			p.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if p.UpdatedAt.Before(beforeTime) || p.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", p.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if p.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", p.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestParamGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	p := Param{CreatedBy: userID}
	got := p.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestParamGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	p := Param{UpdatedBy: userID}
	got := p.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestParamGetCreatedAt(t *testing.T) {
	now := time.Now()
	p := Param{CreatedAt: now}
	got := p.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestParamGetUpdatedAt(t *testing.T) {
	now := time.Now()
	p := Param{UpdatedAt: now}
	got := p.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestParamSetCreatedAt(t *testing.T) {
	now := time.Now()
	p := Param{}
	p.SetCreatedAt(now)

	if p.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", p.CreatedAt, now)
	}
}

func TestParamSetUpdatedAt(t *testing.T) {
	now := time.Now()
	p := Param{}
	p.SetUpdatedAt(now)

	if p.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", p.UpdatedAt, now)
	}
}

func TestParamSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	p := Param{}
	p.SetCreatedBy(userID)

	if p.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", p.CreatedBy, userID)
	}
}

func TestParamSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	p := Param{}
	p.SetUpdatedBy(userID)

	if p.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", p.UpdatedBy, userID)
	}
}

func TestParamIsZero(t *testing.T) {
	tests := []struct {
		name string
		p    Param
		want bool
	}{
		{
			name: "returns true for uninitialized param",
			p:    Param{},
			want: true,
		},
		{
			name: "returns false for initialized param",
			p:    Param{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamSlug(t *testing.T) {
	tests := []struct {
		name    string
		pName   string
		shortID string
		want    string
	}{
		{
			name:    "generates slug from name and short ID",
			pName:   "site_title",
			shortID: "abc123",
			want:    "site_title-abc123",
		},
		{
			name:    "handles special characters",
			pName:   "max_items",
			shortID: "xyz789",
			want:    "max_items-xyz789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Param{
				Name:    tt.pName,
				ShortID: tt.shortID,
			}
			got := p.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParamRef(t *testing.T) {
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
			p := Param{ref: tt.ref}
			got := p.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestParamSetRef(t *testing.T) {
	p := Param{}
	ref := "new-ref"
	p.SetRef(ref)

	if p.ref != ref {
		t.Errorf("SetRef() set %v, want %v", p.ref, ref)
	}
}
