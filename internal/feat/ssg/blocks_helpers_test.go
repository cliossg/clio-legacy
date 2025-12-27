package ssg

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLimit(t *testing.T) {
	tests := []struct {
		name     string
		content  []Content
		max      int
		wantLen int
	}{
		{
			name: "limits content when exceeds max",
			content: []Content{
				{ID: uuid.New(), Heading: "First"},
				{ID: uuid.New(), Heading: "Second"},
				{ID: uuid.New(), Heading: "Third"},
				{ID: uuid.New(), Heading: "Fourth"},
				{ID: uuid.New(), Heading: "Fifth"},
			},
			max:     3,
			wantLen: 3,
		},
		{
			name: "returns all content when below max",
			content: []Content{
				{ID: uuid.New(), Heading: "First"},
				{ID: uuid.New(), Heading: "Second"},
			},
			max:     5,
			wantLen: 2,
		},
		{
			name:     "returns empty slice for empty input",
			content:  []Content{},
			max:      3,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := limit(tt.content, tt.max)

			if len(got) != tt.wantLen {
				t.Errorf("limit() returned %d items, want %d", len(got), tt.wantLen)
			}

			if len(tt.content) > tt.max && len(got) == tt.max {
				for i := 0; i < tt.max; i++ {
					if got[i].ID != tt.content[i].ID {
						t.Errorf("limit() item %d = %v, want %v", i, got[i].ID, tt.content[i].ID)
					}
				}
			}
		})
	}
}

func TestHasCommonTags(t *testing.T) {
	tag1 := Tag{ID: uuid.New(), Name: "go"}
	tag2 := Tag{ID: uuid.New(), Name: "testing"}
	tag3 := Tag{ID: uuid.New(), Name: "performance"}

	tests := []struct {
		name string
		c1   Content
		c2   Content
		want bool
	}{
		{
			name: "returns true when tags match",
			c1:   Content{Tags: []Tag{tag1, tag2}},
			c2:   Content{Tags: []Tag{tag1, tag3}},
			want: true,
		},
		{
			name: "returns false when no tags match",
			c1:   Content{Tags: []Tag{tag1}},
			c2:   Content{Tags: []Tag{tag2, tag3}},
			want: false,
		},
		{
			name: "returns false when one has no tags",
			c1:   Content{Tags: []Tag{tag1}},
			c2:   Content{Tags: []Tag{}},
			want: false,
		},
		{
			name: "returns false when both have no tags",
			c1:   Content{Tags: []Tag{}},
			c2:   Content{Tags: []Tag{}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasCommonTags(tt.c1, tt.c2)

			if got != tt.want {
				t.Errorf("hasCommonTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildBlogBlocks(t *testing.T) {
	sectionID := uuid.New()
	tagGo := Tag{ID: uuid.New(), Name: "go"}
	tagTest := Tag{ID: uuid.New(), Name: "testing"}

	now := time.Now()
	past1 := now.Add(-1 * time.Hour)
	past2 := now.Add(-2 * time.Hour)

	currentBlog := Content{
		ID:        uuid.New(),
		SectionID: sectionID,
		Kind:      "blog",
		Heading:   "Current Blog Post",
		Tags:      []Tag{tagGo},
	}

	allContent := []Content{
		currentBlog,
		{
			ID:          uuid.New(),
			SectionID:   sectionID,
			Kind:        "blog",
			Heading:     "Related Blog Post",
			Tags:        []Tag{tagGo},
			PublishedAt: &now,
		},
		{
			ID:          uuid.New(),
			SectionID:   sectionID,
			Kind:        "blog",
			Heading:     "Recent Blog Post 1",
			Tags:        []Tag{tagTest},
			PublishedAt: &past1,
		},
		{
			ID:          uuid.New(),
			SectionID:   sectionID,
			Kind:        "blog",
			Heading:     "Recent Blog Post 2",
			PublishedAt: &past2,
		},
		{
			ID:        uuid.New(),
			SectionID: uuid.New(), // Different section
			Kind:      "blog",
			Heading:   "Other Section Blog",
			Tags:      []Tag{tagGo},
		},
	}

	blocks := &GeneratedBlocks{}
	buildBlogBlocks(blocks, currentBlog, allContent, 5)

	if len(blocks.BlogTagRelated) != 1 {
		t.Errorf("BlogTagRelated count = %d, want 1", len(blocks.BlogTagRelated))
	}

	if len(blocks.BlogRecent) != 2 {
		t.Errorf("BlogRecent count = %d, want 2", len(blocks.BlogRecent))
	}

	if len(blocks.BlogRecent) > 0 && blocks.BlogRecent[0].Heading != "Recent Blog Post 1" {
		t.Errorf("BlogRecent[0] = %s, want 'Recent Blog Post 1'", blocks.BlogRecent[0].Heading)
	}
}

func TestBuildBlocksForBlog(t *testing.T) {
	sectionID := uuid.New()
	tagGo := Tag{ID: uuid.New(), Name: "go"}

	currentBlog := Content{
		ID:        uuid.New(),
		SectionID: sectionID,
		Kind:      "blog",
		Heading:   "Current Blog",
		Tags:      []Tag{tagGo},
	}

	allContent := []Content{
		currentBlog,
		{
			ID:        uuid.New(),
			SectionID: sectionID,
			Kind:      "blog",
			Heading:   "Other Blog",
			Tags:      []Tag{tagGo},
		},
	}

	blocks := BuildBlocks(currentBlog, allContent, 5)

	if blocks == nil {
		t.Fatal("BuildBlocks() returned nil")
	}

	if len(blocks.BlogTagRelated) != 1 {
		t.Errorf("BlogTagRelated count = %d, want 1", len(blocks.BlogTagRelated))
	}
}
