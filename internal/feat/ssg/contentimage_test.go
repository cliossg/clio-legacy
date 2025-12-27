package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewContentImage(t *testing.T) {
	tests := []struct {
		name      string
		contentID uuid.UUID
		imageID   uuid.UUID
		isHeader  bool
	}{
		{
			name:      "creates content image as header",
			contentID: uuid.New(),
			imageID:   uuid.New(),
			isHeader:  true,
		},
		{
			name:      "creates content image as non-header",
			contentID: uuid.New(),
			imageID:   uuid.New(),
			isHeader:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			got := NewContentImage(tt.contentID, tt.imageID, tt.isHeader)
			after := time.Now()

			if got == nil {
				t.Fatal("NewContentImage() returned nil")
			}

			if got.ID == uuid.Nil {
				t.Error("ID was not generated")
			}

			if got.ContentID != tt.contentID {
				t.Errorf("ContentID = %v, want %v", got.ContentID, tt.contentID)
			}

			if got.ImageID != tt.imageID {
				t.Errorf("ImageID = %v, want %v", got.ImageID, tt.imageID)
			}

			if got.IsHeader != tt.isHeader {
				t.Errorf("IsHeader = %v, want %v", got.IsHeader, tt.isHeader)
			}

			if got.IsFeatured != false {
				t.Errorf("IsFeatured = %v, want false", got.IsFeatured)
			}

			if got.OrderNum != 0 {
				t.Errorf("OrderNum = %v, want 0", got.OrderNum)
			}

			if got.CreatedAt.Before(before) || got.CreatedAt.After(after) {
				t.Errorf("CreatedAt not set correctly: %v", got.CreatedAt)
			}
		})
	}
}

func TestNewSectionImage(t *testing.T) {
	tests := []struct {
		name      string
		sectionID uuid.UUID
		imageID   uuid.UUID
		isHeader  bool
	}{
		{
			name:      "creates section image as header",
			sectionID: uuid.New(),
			imageID:   uuid.New(),
			isHeader:  true,
		},
		{
			name:      "creates section image as non-header",
			sectionID: uuid.New(),
			imageID:   uuid.New(),
			isHeader:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			got := NewSectionImage(tt.sectionID, tt.imageID, tt.isHeader)
			after := time.Now()

			if got == nil {
				t.Fatal("NewSectionImage() returned nil")
			}

			if got.ID == uuid.Nil {
				t.Error("ID was not generated")
			}

			if got.SectionID != tt.sectionID {
				t.Errorf("SectionID = %v, want %v", got.SectionID, tt.sectionID)
			}

			if got.ImageID != tt.imageID {
				t.Errorf("ImageID = %v, want %v", got.ImageID, tt.imageID)
			}

			if got.IsHeader != tt.isHeader {
				t.Errorf("IsHeader = %v, want %v", got.IsHeader, tt.isHeader)
			}

			if got.IsFeatured != false {
				t.Errorf("IsFeatured = %v, want false", got.IsFeatured)
			}

			if got.OrderNum != 0 {
				t.Errorf("OrderNum = %v, want 0", got.OrderNum)
			}

			if got.CreatedAt.Before(before) || got.CreatedAt.After(after) {
				t.Errorf("CreatedAt not set correctly: %v", got.CreatedAt)
			}
		})
	}
}
