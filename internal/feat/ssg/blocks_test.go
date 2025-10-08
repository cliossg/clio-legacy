package ssg_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/clio/internal/feat/ssg"
)

func TestBlockBuilder(t *testing.T) {
	allContent := setupBlockBuilderTestData()

	// Find specific content items by a unique property for testing
	findContentByHeading := func(heading string) ssg.Content {
		for _, c := range allContent {
			if c.Heading == heading {
				return c
			}
		}
		return ssg.Content{}
	}

	// TODO: Add extensive test cases here
	tests := []struct {
		name           string
		currentContent ssg.Content
		validate       func(t *testing.T, blocks *ssg.GeneratedBlocks)
	}{
		{
			name:           "Article in tech section with 'go' tag",
			currentContent: findContentByHeading("Advanced Go Generics"),
			validate: func(t *testing.T, blocks *ssg.GeneratedBlocks) {
				if len(blocks.ArticleTagRelatedSameSection) != 1 {
					t.Errorf("Expected 1 tag-related article in same section, got %d", len(blocks.ArticleTagRelatedSameSection))
				}
				if len(blocks.ArticleRecentSameSection) != 1 {
					t.Errorf("Expected 2 recent articles in same section, got %d", len(blocks.ArticleRecentSameSection))
				}
			},
		},
		{
			name:           "Series post (part 3)",
			currentContent: findContentByHeading("Part 3: Concurrency"),
			validate: func(t *testing.T, blocks *ssg.GeneratedBlocks) {
				if blocks.SeriesPrev == nil {
					t.Fatal("SeriesPrev should not be nil")
				}
				if blocks.SeriesPrev.Heading != "Part 2: Interfaces" {
					t.Errorf("Incorrect SeriesPrev: got %q, want %q", blocks.SeriesPrev.Heading, "Part 2: Interfaces")
				}
				if blocks.SeriesNext == nil {
					t.Fatal("SeriesNext should not be nil")
				}
				if blocks.SeriesNext.Heading != "Part 4: Testing" {
					t.Errorf("Incorrect SeriesNext: got %q, want %q", blocks.SeriesNext.Heading, "Part 4: Testing")
				}
				if len(blocks.SeriesIndexForward) != 2 {
					t.Errorf("Expected 2 forward series posts, got %d", len(blocks.SeriesIndexForward))
				}
				if len(blocks.SeriesIndexBackward) != 2 {
					t.Errorf("Expected 2 backward series posts, got %d", len(blocks.SeriesIndexBackward))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks := ssg.BuildBlocks(tt.currentContent, allContent, 5)
			tt.validate(t, blocks)
		})
	}
}

// setupBlockBuilderTestData is at the end of the file to reduce noise.
func setupBlockBuilderTestData() []ssg.Content {
	// --- UUIDs for entities ---
	// Sections
	secRootID := uuid.New()
	secTechID := uuid.New()
	secTravelID := uuid.New()
	secCookingID := uuid.New()

	// Tags
	tagGoID := uuid.New()
	tagTestingID := uuid.New()
	tagPerfID := uuid.New()
	tagItalyID := uuid.New()
	tagFoodID := uuid.New()
	tagPastaID := uuid.New()
	tagBakingID := uuid.New()

	// --- Entities ---
	// Sections
	sections := map[string]ssg.Section{
		"root":    {ID: secRootID, Name: "root", Path: "/"},
		"tech":    {ID: secTechID, Name: "tech", Path: "/tech"},
		"travel":  {ID: secTravelID, Name: "travel", Path: "/travel"},
		"cooking": {ID: secCookingID, Name: "cooking", Path: "/cooking"},
	}

	// Tags
	tags := map[string]ssg.Tag{
		"go":      {ID: tagGoID, Name: "go"},
		"testing": {ID: tagTestingID, Name: "testing"},
		"perf":    {ID: tagPerfID, Name: "performance"},
		"italy":   {ID: tagItalyID, Name: "italy"},
		"food":    {ID: tagFoodID, Name: "food"},
		"pasta":   {ID: tagPastaID, Name: "pasta"},
		"baking":  {ID: tagBakingID, Name: "baking"},
	}

	now := time.Now()
	content := []ssg.Content{
		// === Tech Section ===
		{ID: uuid.New(), SectionID: secTechID, Kind: "page", Heading: "About Tech"},
		{ID: uuid.New(), SectionID: secTechID, Kind: "article", Heading: "Advanced Go Generics", Tags: []ssg.Tag{tags["go"], tags["perf"]}, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTechID, Kind: "article", Heading: "Introduction to Go", Tags: []ssg.Tag{tags["go"]}, PublishedAt: func() *time.Time { t := now.Add(-1 * time.Hour); return &t }()},
		{ID: uuid.New(), SectionID: secTechID, Kind: "article", Heading: "Why I Don't Like Java", PublishedAt: func() *time.Time { t := now.Add(-2 * time.Hour); return &t }()},
		{ID: uuid.New(), SectionID: secTechID, Kind: "blog", Heading: "Tech Blog Post 1", Tags: []ssg.Tag{tags["go"]}, PublishedAt: &now},
		// Series: Learning Go
		{ID: uuid.New(), SectionID: secTechID, Kind: "series", Heading: "Part 1: Introduction", Series: "Learning Go", SeriesOrder: 1, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTechID, Kind: "series", Heading: "Part 2: Interfaces", Series: "Learning Go", SeriesOrder: 2, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTechID, Kind: "series", Heading: "Part 3: Concurrency", Series: "Learning Go", SeriesOrder: 3, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTechID, Kind: "series", Heading: "Part 4: Testing", Series: "Learning Go", SeriesOrder: 4, Tags: []ssg.Tag{tags["testing"]}, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTechID, Kind: "series", Heading: "Part 5: Performance", Series: "Learning Go", SeriesOrder: 5, Tags: []ssg.Tag{tags["perf"]}, PublishedAt: &now},

		// === Travel Section ===
		{ID: uuid.New(), SectionID: secTravelID, Kind: "page", Heading: "Travel Destinations"},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "article", Heading: "A Trip to Rome", Tags: []ssg.Tag{tags["italy"]}, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "blog", Heading: "Travel Blog 1", PublishedAt: &now},
		// Series: Italian Tour
		{ID: uuid.New(), SectionID: secTravelID, Kind: "series", Heading: "Day 1: Rome", Series: "Italian Tour", SeriesOrder: 1, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "series", Heading: "Day 2: Florence", Series: "Italian Tour", SeriesOrder: 2, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "series", Heading: "Day 3: Venice", Series: "Italian Tour", SeriesOrder: 3, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "series", Heading: "Day 4: Milan", Series: "Italian Tour", SeriesOrder: 4, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secTravelID, Kind: "series", Heading: "Day 5: Cinque Terre", Series: "Italian Tour", SeriesOrder: 5, PublishedAt: &now},

		// === Cooking Section ===
		{ID: uuid.New(), SectionID: secCookingID, Kind: "page", Heading: "Recipes"},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "article", Heading: "The Art of Pasta", Tags: []ssg.Tag{tags["food"], tags["pasta"], tags["italy"]}, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "blog", Heading: "Cooking Blog 1", Tags: []ssg.Tag{tags["food"]}, PublishedAt: &now},
		// Series: Pasta Making
		{ID: uuid.New(), SectionID: secCookingID, Kind: "series", Heading: "Step 1: The Dough", Series: "Pasta Making", SeriesOrder: 1, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "series", Heading: "Step 2: Rolling", Series: "Pasta Making", SeriesOrder: 2, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "series", Heading: "Step 3: Cutting", Series: "Pasta Making", SeriesOrder: 3, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "series", Heading: "Step 4: Cooking", Series: "Pasta Making", SeriesOrder: 4, PublishedAt: &now},
		{ID: uuid.New(), SectionID: secCookingID, Kind: "series", Heading: "Step 5: The Sauce", Series: "Pasta Making", SeriesOrder: 5, PublishedAt: &now},

		// === Root Section ===
		{ID: uuid.New(), SectionID: secRootID, Kind: "page", Heading: "Welcome Home"},
		{ID: uuid.New(), SectionID: secRootID, Kind: "article", Heading: "Site Philosophy", PublishedAt: &now},
	}

	// Assign section paths to content
	for i, c := range content {
		for _, s := range sections {
			if c.SectionID == s.ID {
				content[i].SectionPath = s.Path
				content[i].SectionName = s.Name
				break
			}
		}
	}

	return content
}
