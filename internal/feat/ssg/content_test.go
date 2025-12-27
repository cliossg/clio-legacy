package ssg

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewContent(t *testing.T) {
	tests := []struct {
		name    string
		heading string
		body    string
		want    Content
	}{
		{
			name:    "creates content with heading and body",
			heading: "Test Heading",
			body:    "Test body content",
			want: Content{
				Heading: "Test Heading",
				Body:    "Test body content",
				Draft:   true,
			},
		},
		{
			name:    "creates content with empty strings",
			heading: "",
			body:    "",
			want: Content{
				Heading: "",
				Body:    "",
				Draft:   true,
			},
		},
		{
			name:    "creates content with long body",
			heading: "Long Article",
			body:    strings.Repeat("Lorem ipsum ", 100),
			want: Content{
				Heading: "Long Article",
				Body:    strings.Repeat("Lorem ipsum ", 100),
				Draft:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewContent(tt.heading, tt.body)

			if got.Heading != tt.want.Heading {
				t.Errorf("Heading = %v, want %v", got.Heading, tt.want.Heading)
			}
			if got.Body != tt.want.Body {
				t.Errorf("Body = %v, want %v", got.Body, tt.want.Body)
			}
			if got.Draft != tt.want.Draft {
				t.Errorf("Draft = %v, want %v", got.Draft, tt.want.Draft)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	c := Content{}
	got := c.Type()
	want := "content"

	if got != want {
		t.Errorf("Type() = %v, want %v", got, want)
	}
}

func TestContentGetID(t *testing.T) {
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
			c := Content{ID: tt.id}
			got := c.GetID()

			if got != tt.id {
				t.Errorf("GetID() = %v, want %v", got, tt.id)
			}
		})
	}
}

func TestContentGenID(t *testing.T) {
	c := Content{}
	c.GenID()

	if c.ID == uuid.Nil {
		t.Error("GenID() did not generate a UUID")
	}
}

func TestContentSetID(t *testing.T) {
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
			expected: uuid.Nil, // Will be set to newID
		},
		{
			name:     "does not override existing ID without force",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    nil,
			expected: uuid.Nil, // Will keep initial
		},
		{
			name:     "overrides existing ID with force true",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{true},
			expected: uuid.Nil, // Will be set to newID
		},
		{
			name:     "does not override with force false",
			initial:  uuid.New(),
			newID:    uuid.New(),
			force:    []bool{false},
			expected: uuid.Nil, // Will keep initial
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Content{ID: tt.initial}
			c.SetID(tt.newID, tt.force...)

			// Determine expected based on logic
			var want uuid.UUID
			if tt.initial == uuid.Nil {
				want = tt.newID
			} else if len(tt.force) > 0 && tt.force[0] && tt.newID != uuid.Nil {
				want = tt.newID
			} else {
				want = tt.initial
			}

			if c.ID != want {
				t.Errorf("SetID() resulted in %v, want %v", c.ID, want)
			}
		})
	}
}

func TestContentGetShortID(t *testing.T) {
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
			c := Content{ShortID: tt.shortID}
			got := c.GetShortID()

			if got != tt.shortID {
				t.Errorf("GetShortID() = %v, want %v", got, tt.shortID)
			}
		})
	}
}

func TestContentGenShortID(t *testing.T) {
	c := Content{}
	c.GenShortID()

	if c.ShortID == "" {
		t.Error("GenShortID() did not generate a short ID")
	}

	if len(c.ShortID) != 12 {
		t.Errorf("GenShortID() generated ID of length %d, want 12", len(c.ShortID))
	}
}

func TestContentSetShortID(t *testing.T) {
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
			c := Content{ShortID: tt.initial}
			c.SetShortID(tt.newID, tt.force...)

			if c.ShortID != tt.expected {
				t.Errorf("SetShortID() resulted in %v, want %v", c.ShortID, tt.expected)
			}
		})
	}
}

func TestContentGenCreateValues(t *testing.T) {
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
			c := Content{}
			beforeTime := time.Now()
			c.GenCreateValues(tt.userID...)
			afterTime := time.Now()

			if c.CreatedAt.Before(beforeTime) || c.CreatedAt.After(afterTime) {
				t.Errorf("CreatedAt not set correctly: %v", c.CreatedAt)
			}

			if c.UpdatedAt.Before(beforeTime) || c.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", c.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if c.CreatedBy != tt.userID[0] {
					t.Errorf("CreatedBy = %v, want %v", c.CreatedBy, tt.userID[0])
				}
				if c.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", c.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestContentGenUpdateValues(t *testing.T) {
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
			c := Content{}
			beforeTime := time.Now()
			c.GenUpdateValues(tt.userID...)
			afterTime := time.Now()

			if c.UpdatedAt.Before(beforeTime) || c.UpdatedAt.After(afterTime) {
				t.Errorf("UpdatedAt not set correctly: %v", c.UpdatedAt)
			}

			if len(tt.userID) > 0 {
				if c.UpdatedBy != tt.userID[0] {
					t.Errorf("UpdatedBy = %v, want %v", c.UpdatedBy, tt.userID[0])
				}
			}
		})
	}
}

func TestContentGetCreatedBy(t *testing.T) {
	userID := uuid.New()
	c := Content{CreatedBy: userID}
	got := c.GetCreatedBy()

	if got != userID {
		t.Errorf("GetCreatedBy() = %v, want %v", got, userID)
	}
}

func TestContentGetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	c := Content{UpdatedBy: userID}
	got := c.GetUpdatedBy()

	if got != userID {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, userID)
	}
}

func TestContentGetCreatedAt(t *testing.T) {
	now := time.Now()
	c := Content{CreatedAt: now}
	got := c.GetCreatedAt()

	if got != now {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestContentGetUpdatedAt(t *testing.T) {
	now := time.Now()
	c := Content{UpdatedAt: now}
	got := c.GetUpdatedAt()

	if got != now {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestContentSetCreatedAt(t *testing.T) {
	now := time.Now()
	c := Content{}
	c.SetCreatedAt(now)

	if c.CreatedAt != now {
		t.Errorf("SetCreatedAt() set %v, want %v", c.CreatedAt, now)
	}
}

func TestContentSetUpdatedAt(t *testing.T) {
	now := time.Now()
	c := Content{}
	c.SetUpdatedAt(now)

	if c.UpdatedAt != now {
		t.Errorf("SetUpdatedAt() set %v, want %v", c.UpdatedAt, now)
	}
}

func TestContentSetCreatedBy(t *testing.T) {
	userID := uuid.New()
	c := Content{}
	c.SetCreatedBy(userID)

	if c.CreatedBy != userID {
		t.Errorf("SetCreatedBy() set %v, want %v", c.CreatedBy, userID)
	}
}

func TestContentSetUpdatedBy(t *testing.T) {
	userID := uuid.New()
	c := Content{}
	c.SetUpdatedBy(userID)

	if c.UpdatedBy != userID {
		t.Errorf("SetUpdatedBy() set %v, want %v", c.UpdatedBy, userID)
	}
}

func TestContentIsZero(t *testing.T) {
	tests := []struct {
		name string
		c    Content
		want bool
	}{
		{
			name: "returns true for uninitialized content",
			c:    Content{},
			want: true,
		},
		{
			name: "returns false for initialized content",
			c:    Content{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.IsZero()

			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentSlug(t *testing.T) {
	tests := []struct {
		name    string
		heading string
		shortID string
		want    string
	}{
		{
			name:    "generates slug from heading and short ID",
			heading: "My First Post",
			shortID: "abc123",
			want:    "my-first-post-abc123",
		},
		{
			name:    "handles special characters",
			heading: "Hello, World!",
			shortID: "xyz789",
			want:    "hello,-world!-xyz789",
		},
		{
			name:    "handles multiple spaces",
			heading: "Multiple   Spaces",
			shortID: "def456",
			want:    "multiple---spaces-def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Content{
				Heading: tt.heading,
				ShortID: tt.shortID,
			}
			got := c.Slug()

			if got != tt.want {
				t.Errorf("Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentOptValue(t *testing.T) {
	id := uuid.New()
	c := Content{ID: id}
	got := c.OptValue()
	want := id.String()

	if got != want {
		t.Errorf("OptValue() = %v, want %v", got, want)
	}
}

func TestContentOptLabel(t *testing.T) {
	c := Content{Heading: "Test Heading"}
	got := c.OptLabel()
	want := "Test Heading"

	if got != want {
		t.Errorf("OptLabel() = %v, want %v", got, want)
	}
}

func TestContentRef(t *testing.T) {
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
			c := Content{ref: tt.ref}
			got := c.Ref()

			if got != tt.ref {
				t.Errorf("Ref() = %v, want %v", got, tt.ref)
			}
		})
	}
}

func TestContentSetRef(t *testing.T) {
	c := Content{}
	ref := "new-ref"
	c.SetRef(ref)

	if c.ref != ref {
		t.Errorf("SetRef() set %v, want %v", c.ref, ref)
	}
}

func TestContentString(t *testing.T) {
	tests := []struct {
		name     string
		content  Content
		contains []string
	}{
		{
			name: "formats content with short body",
			content: Content{
				ID:      uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				Heading: "Test Post",
				Body:    "Short body",
				Draft:   true,
			},
			contains: []string{"12345678", "Test Post", "Short body", "true"},
		},
		{
			name: "truncates long body",
			content: Content{
				ID:      uuid.MustParse("12345678-1234-1234-1234-123456789012"),
				Heading: "Long Post",
				Body:    strings.Repeat("a", 100),
				Draft:   false,
			},
			contains: []string{"12345678", "Long Post", "...", "false"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.content.String()

			for _, substr := range tt.contains {
				if !strings.Contains(got, substr) {
					t.Errorf("String() = %v, should contain %v", got, substr)
				}
			}

			// Check body truncation for long body
			if len(tt.content.Body) > 50 {
				if strings.Contains(got, strings.Repeat("a", 100)) {
					t.Error("String() did not truncate long body")
				}
			}
		})
	}
}
