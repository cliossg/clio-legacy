package ssg

import (
	"testing"
)

func TestGetContentPath(t *testing.T) {
	tests := []struct {
		name    string
		content Content
		mode    string
		want    string
	}{
		{
			name: "blog mode returns root path",
			content: Content{
				Heading:     "Test Post",
				ShortID:     "abc123",
				SectionPath: "articles",
			},
			mode: "blog",
			want: "/test-post-abc123",
		},
		{
			name: "structured mode with section path",
			content: Content{
				Heading:     "My Article",
				ShortID:     "xyz789",
				SectionPath: "docs/guides",
			},
			mode: "structured",
			want: "/docs/guides/my-article-xyz789",
		},
		{
			name: "structured mode with root section",
			content: Content{
				Heading:     "Home Page",
				ShortID:     "def456",
				SectionPath: "/",
			},
			mode: "structured",
			want: "/home-page-def456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContentPath(tt.content, tt.mode)

			if got != tt.want {
				t.Errorf("GetContentPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIndexPath(t *testing.T) {
	tests := []struct {
		name        string
		sectionPath string
		contentType string
		mode        string
		want        string
	}{
		{
			name:        "blog mode returns root",
			sectionPath: "articles",
			contentType: "blog",
			mode:        "blog",
			want:        "/",
		},
		{
			name:        "structured mode blog type non-root section",
			sectionPath: "/tech",
			contentType: "blog",
			mode:        "structured",
			want:        "/tech/blog/",
		},
		{
			name:        "structured mode blog type root section",
			sectionPath: "/",
			contentType: "blog",
			mode:        "structured",
			want:        "/blog/",
		},
		{
			name:        "structured mode non-blog type",
			sectionPath: "/docs/guides",
			contentType: "page",
			mode:        "structured",
			want:        "/docs/guides/",
		},
		{
			name:        "structured mode non-blog type root",
			sectionPath: "/",
			contentType: "page",
			mode:        "structured",
			want:        "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIndexPath(tt.sectionPath, tt.contentType, tt.mode)

			if got != tt.want {
				t.Errorf("GetIndexPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPaginationPath(t *testing.T) {
	tests := []struct {
		name      string
		indexPath string
		page      int
		mode      string
		want      string
	}{
		{
			name:      "first page at root",
			indexPath: "",
			page:      1,
			mode:      "blog",
			want:      "/",
		},
		{
			name:      "first page with path",
			indexPath: "/blog",
			page:      1,
			mode:      "structured",
			want:      "/blog/",
		},
		{
			name:      "second page at root",
			indexPath: "",
			page:      2,
			mode:      "blog",
			want:      "/page/2/",
		},
		{
			name:      "second page with path",
			indexPath: "/blog",
			page:      2,
			mode:      "structured",
			want:      "/blog/page/2/",
		},
		{
			name:      "third page with nested path",
			indexPath: "/docs/guides",
			page:      3,
			mode:      "structured",
			want:      "/docs/guides/page/3/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPaginationPath(tt.indexPath, tt.page, tt.mode)

			if got != tt.want {
				t.Errorf("GetPaginationPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContentFilePath(t *testing.T) {
	tests := []struct {
		name     string
		htmlPath string
		content  Content
		mode     string
		want     string
	}{
		{
			name:     "blog mode file path",
			htmlPath: "/var/www/html",
			content: Content{
				Heading:     "Test Post",
				ShortID:     "abc123",
				SectionPath: "ignored",
			},
			mode: "blog",
			want: "/var/www/html/test-post-abc123/index.html",
		},
		{
			name:     "structured mode with section",
			htmlPath: "/var/www/html",
			content: Content{
				Heading:     "Guide",
				ShortID:     "xyz789",
				SectionPath: "docs/tutorials",
			},
			mode: "structured",
			want: "/var/www/html/docs/tutorials/guide-xyz789/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContentFilePath(tt.htmlPath, tt.content, tt.mode)

			if got != tt.want {
				t.Errorf("GetContentFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIndexFilePath(t *testing.T) {
	tests := []struct {
		name      string
		htmlPath  string
		indexPath string
		want      string
	}{
		{
			name:      "root index",
			htmlPath:  "/var/www/html",
			indexPath: "/",
			want:      "/var/www/html/index.html",
		},
		{
			name:      "nested index",
			htmlPath:  "/var/www/html",
			indexPath: "blog",
			want:      "/var/www/html/blog/index.html",
		},
		{
			name:      "deep nested index",
			htmlPath:  "/var/www/html",
			indexPath: "docs/guides/advanced",
			want:      "/var/www/html/docs/guides/advanced/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetIndexFilePath(tt.htmlPath, tt.indexPath)

			if got != tt.want {
				t.Errorf("GetIndexFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPaginationFilePath(t *testing.T) {
	tests := []struct {
		name      string
		htmlPath  string
		indexPath string
		page      int
		want      string
	}{
		{
			name:      "first page at root",
			htmlPath:  "/var/www/html",
			indexPath: "/",
			page:      1,
			want:      "/var/www/html/index.html",
		},
		{
			name:      "second page at root",
			htmlPath:  "/var/www/html",
			indexPath: "/",
			page:      2,
			want:      "/var/www/html/page/2/index.html",
		},
		{
			name:      "first page with path",
			htmlPath:  "/var/www/html",
			indexPath: "blog",
			page:      1,
			want:      "/var/www/html/blog/index.html",
		},
		{
			name:      "third page with path",
			htmlPath:  "/var/www/html",
			indexPath: "blog",
			page:      3,
			want:      "/var/www/html/blog/page/3/index.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPaginationFilePath(tt.htmlPath, tt.indexPath, tt.page)

			if got != tt.want {
				t.Errorf("GetPaginationFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteBasePath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard site path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteBasePath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteDBPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard db path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/db/clio.db",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/db/clio.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteDBPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteDBPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteDBDSN(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard DSN",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "file:_workspace/sites/my-blog/db/clio.db?cache=shared&mode=rwc",
		},
		{
			name:          "custom base path DSN",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "file:/var/clio/sites/portfolio/db/clio.db?cache=shared&mode=rwc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteDBDSN(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteDBDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteDocsPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard docs path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/documents",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/documents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteDocsPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteDocsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteMarkdownPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard markdown path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/documents/markdown",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/documents/markdown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteMarkdownPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteMarkdownPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteHTMLPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard HTML path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/documents/html",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/documents/html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteHTMLPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteHTMLPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteAssetsPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard assets path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/documents/assets",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/documents/assets",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteAssetsPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteAssetsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSiteImagesPath(t *testing.T) {
	tests := []struct {
		name          string
		sitesBasePath string
		siteSlug      string
		want          string
	}{
		{
			name:          "standard images path",
			sitesBasePath: "_workspace/sites",
			siteSlug:      "my-blog",
			want:          "_workspace/sites/my-blog/documents/assets/images",
		},
		{
			name:          "custom base path",
			sitesBasePath: "/var/clio/sites",
			siteSlug:      "portfolio",
			want:          "/var/clio/sites/portfolio/documents/assets/images",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSiteImagesPath(tt.sitesBasePath, tt.siteSlug)

			if got != tt.want {
				t.Errorf("GetSiteImagesPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
