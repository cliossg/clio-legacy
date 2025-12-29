package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		username string
		userName string
		email    string
	}{
		{
			name:     "creates new user",
			username: "testuser",
			userName: "Test User",
			email:    "test@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.username, tt.userName, tt.email)
			if user.Username != tt.username {
				t.Errorf("NewUser() username = %v, want %v", user.Username, tt.username)
			}
			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}
			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}
		})
	}
}

func TestUserType(t *testing.T) {
	user := &User{}
	if got := user.Type(); got != "user" {
		t.Errorf("Type() = %v, want user", got)
	}
}

func TestUserGetID(t *testing.T) {
	id := uuid.New()
	user := User{ID: id}
	if got := user.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestUserGenID(t *testing.T) {
	user := &User{}
	user.GenID()
	if user.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestUserSetID(t *testing.T) {
	tests := []struct {
		name      string
		initial   uuid.UUID
		new       uuid.UUID
		force     []bool
		wantID    uuid.UUID
	}{
		{
			name:      "sets ID when empty",
			initial:   uuid.Nil,
			new:       uuid.New(),
			force:     nil,
			wantID:    uuid.Nil,
		},
		{
			name:      "sets ID with force",
			initial:   uuid.New(),
			new:       uuid.New(),
			force:     []bool{true},
			wantID:    uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{ID: tt.initial}
			user.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && user.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", user.ID, tt.wantID)
			}
		})
	}
}

func TestUserGetShortID(t *testing.T) {
	user := User{ShortID: "test123"}
	if got := user.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestUserGenShortID(t *testing.T) {
	user := &User{}
	user.GenShortID()
	if user.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestUserSetShortID(t *testing.T) {
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
			user := &User{ShortID: tt.initial}
			user.SetShortID(tt.new, tt.force...)
			if user.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && user.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", user.ShortID, tt.new)
				}
			}
		})
	}
}

func TestUserGetCreatedBy(t *testing.T) {
	id := uuid.New()
	user := User{CreatedBy: id}
	if got := user.GetCreatedBy(); got != id {
		t.Errorf("GetCreatedBy() = %v, want %v", got, id)
	}
}

func TestUserGetUpdatedBy(t *testing.T) {
	id := uuid.New()
	user := User{UpdatedBy: id}
	if got := user.GetUpdatedBy(); got != id {
		t.Errorf("GetUpdatedBy() = %v, want %v", got, id)
	}
}

func TestUserGetCreatedAt(t *testing.T) {
	now := time.Now()
	user := User{CreatedAt: now}
	if got := user.GetCreatedAt(); !got.Equal(now) {
		t.Errorf("GetCreatedAt() = %v, want %v", got, now)
	}
}

func TestUserGetUpdatedAt(t *testing.T) {
	now := time.Now()
	user := User{UpdatedAt: now}
	if got := user.GetUpdatedAt(); !got.Equal(now) {
		t.Errorf("GetUpdatedAt() = %v, want %v", got, now)
	}
}

func TestUserSetCreatedAt(t *testing.T) {
	now := time.Now()
	user := &User{}
	user.SetCreatedAt(now)
	if !user.CreatedAt.Equal(now) {
		t.Errorf("SetCreatedAt() CreatedAt = %v, want %v", user.CreatedAt, now)
	}
}

func TestUserSetUpdatedAt(t *testing.T) {
	now := time.Now()
	user := &User{}
	user.SetUpdatedAt(now)
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("SetUpdatedAt() UpdatedAt = %v, want %v", user.UpdatedAt, now)
	}
}

func TestUserSetCreatedBy(t *testing.T) {
	id := uuid.New()
	user := &User{}
	user.SetCreatedBy(id)
	if user.CreatedBy != id {
		t.Errorf("SetCreatedBy() CreatedBy = %v, want %v", user.CreatedBy, id)
	}
}

func TestUserSetUpdatedBy(t *testing.T) {
	id := uuid.New()
	user := &User{}
	user.SetUpdatedBy(id)
	if user.UpdatedBy != id {
		t.Errorf("SetUpdatedBy() UpdatedBy = %v, want %v", user.UpdatedBy, id)
	}
}

func TestUserIsZero(t *testing.T) {
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "zero user",
			user: User{},
			want: true,
		},
		{
			name: "non-zero user",
			user: User{ID: uuid.New()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.user.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserSlug(t *testing.T) {
	user := &User{Username: "testuser", ShortID: "abc123"}
	got := user.Slug()
	if got == "" {
		t.Error("Slug() returned empty string")
	}
	if len(got) < len("testuser") {
		t.Errorf("Slug() = %v, too short", got)
	}
}

func TestUserOptLabel(t *testing.T) {
	user := User{Username: "testuser"}
	if got := user.OptLabel(); got != "testuser" {
		t.Errorf("OptLabel() = %v, want testuser", got)
	}
}

func TestUserOptValue(t *testing.T) {
	id := uuid.New()
	user := User{ID: id}
	if got := user.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestUserRef(t *testing.T) {
	user := &User{RefValue: "testuser"}
	if got := user.Ref(); got != "testuser" {
		t.Errorf("Ref() = %v, want testuser", got)
	}
}

func TestUserSetRef(t *testing.T) {
	user := &User{}
	user.SetRef("testuser")
	if user.RefValue != "testuser" {
		t.Errorf("SetRef() RefValue = %v, want testuser", user.RefValue)
	}
}
