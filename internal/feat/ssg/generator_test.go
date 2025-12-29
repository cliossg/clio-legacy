package ssg

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hermesgen/hm"
	"gopkg.in/yaml.v2"
)

func TestGeneratorGenerate(t *testing.T) {
	tests := []struct {
		name          string
		siteSlug      string
		contents      []Content
		wantFileCount int
		checkFiles    func(*testing.T, string, []Content)
	}{
		{
			name:     "generates single content file",
			siteSlug: "test-site",
			contents: []Content{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Heading:     "Test Article",
					Body:        "This is the content body.",
					SectionPath: "blog",
					SectionName: "blog",
					Draft:       false,
					Featured:    true,
					Meta: Meta{
						Description: "Test description",
						Keywords:    "test,article",
						Robots:      "index,follow",
					},
				},
			},
			wantFileCount: 1,
			checkFiles: func(t *testing.T, basePath string, contents []Content) {
				content := contents[0]
				expectedPath := filepath.Join(basePath, content.SectionPath, content.Slug()+".md")
				data, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				fileContent := string(data)
				if !strings.Contains(fileContent, "---") {
					t.Error("File should contain YAML frontmatter delimiters")
				}
				if !strings.Contains(fileContent, "title: Test Article") {
					t.Error("File should contain title in frontmatter")
				}
				if !strings.Contains(fileContent, "This is the content body.") {
					t.Error("File should contain body content")
				}
				if !strings.Contains(fileContent, "draft: false") {
					t.Error("File should contain draft status")
				}
				if !strings.Contains(fileContent, "featured: true") {
					t.Error("File should contain featured status")
				}
			},
		},
		{
			name:     "generates multiple content files",
			siteSlug: "multi-site",
			contents: []Content{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Heading:     "First Post",
					Body:        "First post content.",
					SectionPath: "blog",
					SectionName: "blog",
				},
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Heading:     "Second Post",
					Body:        "Second post content.",
					SectionPath: "articles",
					SectionName: "articles",
				},
			},
			wantFileCount: 2,
			checkFiles: func(t *testing.T, basePath string, contents []Content) {
				for _, content := range contents {
					expectedPath := filepath.Join(basePath, content.SectionPath, content.Slug()+".md")
					if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
						t.Errorf("Expected file not found: %s", expectedPath)
					}
				}
			},
		},
		{
			name:     "generates content with tags",
			siteSlug: "tags-site",
			contents: []Content{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Heading:     "Tagged Post",
					Body:        "Post with tags.",
					SectionPath: "blog",
					SectionName: "blog",
					Tags: []Tag{
						{Name: "golang"},
						{Name: "testing"},
					},
				},
			},
			wantFileCount: 1,
			checkFiles: func(t *testing.T, basePath string, contents []Content) {
				content := contents[0]
				expectedPath := filepath.Join(basePath, content.SectionPath, content.Slug()+".md")
				data, err := os.ReadFile(expectedPath)
				if err != nil {
					t.Fatalf("Failed to read generated file: %v", err)
				}

				fileContent := string(data)
				if !strings.Contains(fileContent, "tags:") {
					t.Error("File should contain tags field")
				}
				if !strings.Contains(fileContent, "golang") {
					t.Error("File should contain golang tag")
				}
				if !strings.Contains(fileContent, "testing") {
					t.Error("File should contain testing tag")
				}
			},
		},
		{
			name:     "generates content without section path",
			siteSlug: "no-section-site",
			contents: []Content{
				{
					ID:          uuid.New(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Heading:     "Root Post",
					Body:        "Post at root.",
					SectionPath: "",
					SectionName: "default",
				},
			},
			wantFileCount: 1,
			checkFiles: func(t *testing.T, basePath string, contents []Content) {
				content := contents[0]
				expectedPath := filepath.Join(basePath, content.Slug()+".md")
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("Expected file not found: %s", expectedPath)
				}
			},
		},
		{
			name:          "handles empty contents list",
			siteSlug:      "empty-site",
			contents:      []Content{},
			wantFileCount: 0,
			checkFiles: func(t *testing.T, basePath string, contents []Content) {
				// Nothing to check
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			cfg := hm.NewConfig()
			cfg.Set(SSGKey.SitesBasePath, tempDir)

			gen := NewGenerator(hm.XParams{Cfg: cfg})
			ctx := context.Background()

			err := gen.Generate(ctx, tt.siteSlug, tt.contents)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			basePath := GetSiteMarkdownPath(tempDir, tt.siteSlug)

			var fileCount int
			err = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(path, ".md") {
					fileCount++
				}
				return nil
			})

			if err != nil && tt.wantFileCount > 0 {
				t.Fatalf("Failed to walk directory: %v", err)
			}

			if fileCount != tt.wantFileCount {
				t.Errorf("Generated %d files, want %d", fileCount, tt.wantFileCount)
			}

			if tt.checkFiles != nil {
				tt.checkFiles(t, basePath, tt.contents)
			}
		})
	}
}

func TestGeneratorGenerateYAMLStructure(t *testing.T) {
	tempDir := t.TempDir()
	cfg := hm.NewConfig()
	cfg.Set(SSGKey.SitesBasePath, tempDir)

	gen := NewGenerator(hm.XParams{Cfg: cfg})

	publishedAt := time.Now()
	content := Content{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Heading:     "YAML Structure Test",
		Body:        "Testing YAML frontmatter structure.",
		SectionPath: "test",
		SectionName: "test",
		Draft:       true,
		Featured:    false,
		PublishedAt: &publishedAt,
		Meta: Meta{
			Description:      "Meta description",
			Keywords:         "yaml,test",
			Robots:           "noindex",
			CanonicalURL:     "https://example.com/test",
			Sitemap:          "priority:0.8",
			TableOfContents:  true,
			Comments:         false,
			Share:            true,
		},
		HeaderImageURL: "/images/header.jpg",
	}

	ctx := context.Background()
	err := gen.Generate(ctx, "yaml-test", []Content{content})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	filePath := filepath.Join(GetSiteMarkdownPath(tempDir, "yaml-test"), "test", content.Slug()+".md")
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	fileContent := string(data)
	parts := strings.Split(fileContent, "---")
	if len(parts) < 3 {
		t.Fatalf("Expected at least 3 parts (2 delimiters), got %d", len(parts))
	}

	yamlContent := parts[1]
	var frontMatter yaml.MapSlice
	err = yaml.Unmarshal([]byte(yamlContent), &frontMatter)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	expectedFields := map[string]bool{
		"title":             false,
		"slug":              false,
		"draft":             false,
		"featured":          false,
		"description":       false,
		"keywords":          false,
		"robots":            false,
		"canonical-url":     false,
		"sitemap":           false,
		"table-of-contents": false,
		"comments":          false,
		"share":             false,
		"image":             false,
	}

	for _, item := range frontMatter {
		if key, ok := item.Key.(string); ok {
			if _, exists := expectedFields[key]; exists {
				expectedFields[key] = true
			}
		}
	}

	for field, found := range expectedFields {
		if !found {
			t.Errorf("Expected field %q not found in frontmatter", field)
		}
	}

	bodyContent := strings.TrimSpace(parts[2])
	if bodyContent != content.Body {
		t.Errorf("Body content = %q, want %q", bodyContent, content.Body)
	}
}
