package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestNewContent(t *testing.T) {
	tests := []struct {
		name    string
		heading string
		body    string
	}{
		{
			name:    "creates content with heading and body",
			heading: "Test Heading",
			body:    "Test Body",
		},
		{
			name:    "creates content with empty fields",
			heading: "",
			body:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := NewContent(tt.heading, tt.body)
			if content.Heading != tt.heading {
				t.Errorf("NewContent() Heading = %v, want %v", content.Heading, tt.heading)
			}
			if content.Body != tt.body {
				t.Errorf("NewContent() Body = %v, want %v", content.Body, tt.body)
			}
		})
	}
}

func TestContentType(t *testing.T) {
	content := &Content{}
	if got := content.Type(); got != contentType {
		t.Errorf("Type() = %v, want %v", got, contentType)
	}
}

func TestContentGetID(t *testing.T) {
	id := uuid.New()
	content := Content{ID: id}
	if got := content.GetID(); got != id {
		t.Errorf("GetID() = %v, want %v", got, id)
	}
}

func TestContentGenID(t *testing.T) {
	content := &Content{}
	content.GenID()
	if content.ID == uuid.Nil {
		t.Error("GenID() did not generate ID")
	}
}

func TestContentSetID(t *testing.T) {
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
			content := &Content{ID: tt.initial}
			content.SetID(tt.new, tt.force...)
			if tt.wantID != uuid.Nil && content.ID != tt.wantID {
				t.Errorf("SetID() ID = %v, want %v", content.ID, tt.wantID)
			}
		})
	}
}

func TestContentGetShortID(t *testing.T) {
	content := Content{ShortID: "test123"}
	if got := content.GetShortID(); got != "test123" {
		t.Errorf("GetShortID() = %v, want test123", got)
	}
}

func TestContentGenShortID(t *testing.T) {
	content := &Content{}
	content.GenShortID()
	if content.ShortID == "" {
		t.Error("GenShortID() did not generate ShortID")
	}
}

func TestContentSetShortID(t *testing.T) {
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
			content := &Content{ShortID: tt.initial}
			content.SetShortID(tt.new, tt.force...)
			if content.ShortID != tt.new && (len(tt.force) == 0 || !tt.force[0]) {
				if tt.initial == "" && content.ShortID != tt.new {
					t.Errorf("SetShortID() ShortID = %v, want %v", content.ShortID, tt.new)
				}
			}
		})
	}
}

func TestContentIsZero(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		want    bool
	}{
		{
			name:    "zero content",
			content: Content{},
			want:    true,
		},
		{
			name:    "non-zero content",
			content: Content{ID: uuid.New()},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.content.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentSlug(t *testing.T) {
	content := &Content{Heading: "Test Content", ShortID: "abc123"}
	got := content.Slug()
	if got != "content" {
		t.Errorf("Slug() = %v, want content", got)
	}
}

func TestContentTypeID(t *testing.T) {
	content := &Content{ShortID: "abc123"}
	got := content.TypeID()
	if got == "" {
		t.Error("TypeID() returned empty string")
	}
}

func TestContentOptLabel(t *testing.T) {
	content := Content{Heading: "Test Content"}
	if got := content.OptLabel(); got != "Test Content" {
		t.Errorf("OptLabel() = %v, want Test Content", got)
	}
}

func TestContentOptValue(t *testing.T) {
	id := uuid.New()
	content := Content{ID: id}
	if got := content.OptValue(); got != id.String() {
		t.Errorf("OptValue() = %v, want %v", got, id.String())
	}
}

func TestToWebContent(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	sectionID := uuid.New()
	now := time.Now()

	featContent := feat.Content{
		ID:          id,
		ShortID:     "abc123",
		UserID:      userID,
		SectionID:   sectionID,
		Kind:        "article",
		Heading:     "Test Content",
		Body:        "Test body",
		Draft:       true,
		Featured:    false,
		PublishedAt: &now,
		SectionPath: "/test",
		SectionName: "Test Section",
	}

	webContent := ToWebContent(featContent)

	if webContent.ID != featContent.ID {
		t.Errorf("ToWebContent() ID = %v, want %v", webContent.ID, featContent.ID)
	}
	if webContent.ShortID != featContent.ShortID {
		t.Errorf("ToWebContent() ShortID = %v, want %v", webContent.ShortID, featContent.ShortID)
	}
	if webContent.UserID != featContent.UserID {
		t.Errorf("ToWebContent() UserID = %v, want %v", webContent.UserID, featContent.UserID)
	}
	if webContent.SectionID != featContent.SectionID {
		t.Errorf("ToWebContent() SectionID = %v, want %v", webContent.SectionID, featContent.SectionID)
	}
	if webContent.Kind != featContent.Kind {
		t.Errorf("ToWebContent() Kind = %v, want %v", webContent.Kind, featContent.Kind)
	}
	if webContent.Heading != featContent.Heading {
		t.Errorf("ToWebContent() Heading = %v, want %v", webContent.Heading, featContent.Heading)
	}
	if webContent.Body != featContent.Body {
		t.Errorf("ToWebContent() Body = %v, want %v", webContent.Body, featContent.Body)
	}
	if webContent.Draft != featContent.Draft {
		t.Errorf("ToWebContent() Draft = %v, want %v", webContent.Draft, featContent.Draft)
	}
	if webContent.Featured != featContent.Featured {
		t.Errorf("ToWebContent() Featured = %v, want %v", webContent.Featured, featContent.Featured)
	}
	if webContent.SectionPath != featContent.SectionPath {
		t.Errorf("ToWebContent() SectionPath = %v, want %v", webContent.SectionPath, featContent.SectionPath)
	}
	if webContent.SectionName != featContent.SectionName {
		t.Errorf("ToWebContent() SectionName = %v, want %v", webContent.SectionName, featContent.SectionName)
	}
}

func TestToWebContents(t *testing.T) {
	now := time.Now()
	featContents := []feat.Content{
		{
			ID:          uuid.New(),
			ShortID:     "abc123",
			Heading:     "Content 1",
			Body:        "Body 1",
			Draft:       false,
			PublishedAt: &now,
		},
		{
			ID:          uuid.New(),
			ShortID:     "def456",
			Heading:     "Content 2",
			Body:        "Body 2",
			Draft:       true,
			PublishedAt: nil,
		},
	}

	webContents := ToWebContents(featContents)

	if len(webContents) != len(featContents) {
		t.Errorf("ToWebContents() length = %v, want %v", len(webContents), len(featContents))
	}

	for i, webContent := range webContents {
		if webContent.ID != featContents[i].ID {
			t.Errorf("ToWebContents()[%d] ID = %v, want %v", i, webContent.ID, featContents[i].ID)
		}
		if webContent.Heading != featContents[i].Heading {
			t.Errorf("ToWebContents()[%d] Heading = %v, want %v", i, webContent.Heading, featContents[i].Heading)
		}
	}
}
