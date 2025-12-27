package ssg

import (
	"strings"
	"testing"
)

func TestNewMarkdownProcessor(t *testing.T) {
	p := NewMarkdownProcessor()

	if p == nil {
		t.Fatal("NewMarkdownProcessor() returned nil")
	}

	if p.parser == nil {
		t.Error("parser was not initialized")
	}
}

func TestNewMarkdownProcessorWithImageContext(t *testing.T) {
	tests := []struct {
		name         string
		imageContext *ImageContext
	}{
		{
			name: "creates processor with empty image context",
			imageContext: &ImageContext{
				Images: make(map[string]ImageMetadata),
			},
		},
		{
			name: "creates processor with populated image context",
			imageContext: &ImageContext{
				Images: map[string]ImageMetadata{
					"img1": {AltText: "Test alt", Title: "Test Image"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMarkdownProcessorWithImageContext(tt.imageContext)

			if p == nil {
				t.Fatal("NewMarkdownProcessorWithImageContext() returned nil")
			}

			if p.parser == nil {
				t.Error("parser was not initialized")
			}
		})
	}
}

func TestProcessorToHTML(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "converts simple markdown to HTML",
			markdown: "# Heading\n\nParagraph text.",
			want:     "<h1>Heading</h1>\n<p>Paragraph text.</p>\n",
		},
		{
			name:     "converts bold text",
			markdown: "**bold text**",
			want:     "<p><strong>bold text</strong></p>\n",
		},
		{
			name:     "converts italic text",
			markdown: "*italic text*",
			want:     "<p><em>italic text</em></p>\n",
		},
		{
			name:     "converts links",
			markdown: "[link text](https://example.com)",
			want:     "<p><a href=\"https://example.com\">link text</a></p>\n",
		},
		{
			name:     "converts code blocks",
			markdown: "```\ncode\n```",
			want:     "<pre><code>code\n</code></pre>\n",
		},
		{
			name:     "handles empty markdown",
			markdown: "",
			want:     "",
		},
		{
			name:     "converts unordered lists",
			markdown: "- Item 1\n- Item 2",
			want:     "<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n</ul>\n",
		},
		{
			name:     "converts ordered lists",
			markdown: "1. First\n2. Second",
			want:     "<ol>\n<li>First</li>\n<li>Second</li>\n</ol>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMarkdownProcessor()
			got, err := p.ToHTML([]byte(tt.markdown))

			if err != nil {
				t.Errorf("ToHTML() error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("ToHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProcessorToHTMLWithImageContext(t *testing.T) {
	tests := []struct {
		name         string
		markdown     string
		imageContext *ImageContext
		wantContains string
	}{
		{
			name:     "converts markdown with image and adds class",
			markdown: "![alt text](image.jpg)",
			imageContext: &ImageContext{
				Images: make(map[string]ImageMetadata),
			},
			wantContains: `class="prose-img"`,
		},
		{
			name:     "converts image with caption using triple pipe separator",
			markdown: "![alt text|||This is a caption](image.jpg)",
			imageContext: &ImageContext{
				Images: make(map[string]ImageMetadata),
			},
			wantContains: `<figcaption class="prose-figcaption">This is a caption</figcaption>`,
		},
		{
			name:     "handles markdown without images",
			markdown: "# Just a heading",
			imageContext: &ImageContext{
				Images: make(map[string]ImageMetadata),
			},
			wantContains: "<h1>Just a heading</h1>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMarkdownProcessor()
			got, err := p.ToHTMLWithImageContext([]byte(tt.markdown), tt.imageContext)

			if err != nil {
				t.Errorf("ToHTMLWithImageContext() error = %v", err)
				return
			}

			if !strings.Contains(got, tt.wantContains) {
				t.Errorf("ToHTMLWithImageContext() = %q, want to contain %q", got, tt.wantContains)
			}
		})
	}
}

func TestEnhanceImagesInHTML(t *testing.T) {
	tests := []struct {
		name         string
		html         string
		imageContext *ImageContext
		want         string
	}{
		{
			name:         "adds prose-img class to simple image",
			html:         `<img src="test.jpg" alt="test">`,
			imageContext: &ImageContext{Images: make(map[string]ImageMetadata)},
			want:         `<img src="test.jpg" alt="test" class="prose-img">`,
		},
		{
			name:         "wraps image with caption in figure",
			html:         `<img src="test.jpg" alt="alt text|||caption text">`,
			imageContext: &ImageContext{Images: make(map[string]ImageMetadata)},
			want:         `<figure class="prose-figure"><img src="test.jpg" alt="alt text" class="prose-img"><figcaption class="prose-figcaption">caption text</figcaption></figure>`,
		},
		{
			name:         "handles HTML without images",
			html:         `<p>No images here</p>`,
			imageContext: &ImageContext{Images: make(map[string]ImageMetadata)},
			want:         `<p>No images here</p>`,
		},
		{
			name:         "handles multiple images",
			html:         `<img src="1.jpg" alt="first"><img src="2.jpg" alt="second">`,
			imageContext: &ImageContext{Images: make(map[string]ImageMetadata)},
			want:         `<img src="1.jpg" alt="first" class="prose-img"><img src="2.jpg" alt="second" class="prose-img">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := enhanceImagesInHTML(tt.html, tt.imageContext)

			if got != tt.want {
				t.Errorf("enhanceImagesInHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}
