package ssg

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	feat "github.com/hermesgen/clio/internal/feat/ssg"
)

func TestContentFormFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		formData url.Values
		wantErr  bool
		checkFn  func(t *testing.T, form ContentForm)
	}{
		{
			name: "valid form with all fields",
			formData: url.Values{
				"id":                {"123e4567-e89b-12d3-a456-426614174000"},
				"user_id":           {"123e4567-e89b-12d3-a456-426614174001"},
				"section_id":        {"123e4567-e89b-12d3-a456-426614174002"},
				"kind":              {"article"},
				"heading":           {"Test Heading"},
				"body":              {"Test Body"},
				"image":             {"test.jpg"},
				"draft":             {"true"},
				"featured":          {"true"},
				"published_at":      {"2024-01-01T00:00:00Z"},
				"tags":              {"tag1,tag2"},
				"description":       {"Test Description"},
				"keywords":          {"test,keywords"},
				"robots":            {"index,follow"},
				"canonical_url":     {"https://example.com"},
				"sitemap":           {"true"},
				"table_of_contents": {"true"},
				"share":             {"true"},
				"comments":          {"true"},
			},
			wantErr: false,
			checkFn: func(t *testing.T, form ContentForm) {
				if form.Heading != "Test Heading" {
					t.Errorf("Heading = %v, want Test Heading", form.Heading)
				}
				if form.Body != "Test Body" {
					t.Errorf("Body = %v, want Test Body", form.Body)
				}
				if !form.Draft {
					t.Errorf("Draft = %v, want true", form.Draft)
				}
				if !form.Featured {
					t.Errorf("Featured = %v, want true", form.Featured)
				}
				if form.Description != "Test Description" {
					t.Errorf("Description = %v, want Test Description", form.Description)
				}
			},
		},
		{
			name: "minimal form",
			formData: url.Values{
				"heading": {"Minimal"},
				"body":    {"Content"},
			},
			wantErr: false,
			checkFn: func(t *testing.T, form ContentForm) {
				if form.Heading != "Minimal" {
					t.Errorf("Heading = %v, want Minimal", form.Heading)
				}
				if form.Body != "Content" {
					t.Errorf("Body = %v, want Content", form.Body)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			form, err := ContentFormFromRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ContentFormFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, form)
			}
		})
	}
}

func TestToFeatContent(t *testing.T) {
	userID := uuid.New()
	sectionID := uuid.New()

	tests := []struct {
		name    string
		form    ContentForm
		checkFn func(t *testing.T, content feat.Content)
	}{
		{
			name: "basic conversion",
			form: ContentForm{
				Heading: "Test Heading",
				Body:    "Test Body",
				Draft:   true,
			},
			checkFn: func(t *testing.T, content feat.Content) {
				if content.Heading != "Test Heading" {
					t.Errorf("Heading = %v, want Test Heading", content.Heading)
				}
				if content.Body != "Test Body" {
					t.Errorf("Body = %v, want Test Body", content.Body)
				}
				if !content.Draft {
					t.Errorf("Draft = %v, want true", content.Draft)
				}
			},
		},
		{
			name: "with IDs",
			form: ContentForm{
				UserID:    userID.String(),
				SectionID: sectionID.String(),
				Heading:   "Test",
				Body:      "Body",
			},
			checkFn: func(t *testing.T, content feat.Content) {
				if content.UserID != userID {
					t.Errorf("UserID = %v, want %v", content.UserID, userID)
				}
				if content.SectionID != sectionID {
					t.Errorf("SectionID = %v, want %v", content.SectionID, sectionID)
				}
			},
		},
		{
			name: "with published date",
			form: ContentForm{
				Heading:     "Test",
				Body:        "Body",
				PublishedAt: "2024-01-01T00:00:00Z",
			},
			checkFn: func(t *testing.T, content feat.Content) {
				if content.PublishedAt == nil {
					t.Error("PublishedAt is nil, want time value")
					return
				}
				expected, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
				if !content.PublishedAt.Equal(expected) {
					t.Errorf("PublishedAt = %v, want %v", content.PublishedAt, expected)
				}
			},
		},
		{
			name: "with comma-separated tags",
			form: ContentForm{
				Heading: "Test",
				Body:    "Body",
				Tags:    "tag1,tag2,tag3",
			},
			checkFn: func(t *testing.T, content feat.Content) {
				if len(content.Tags) != 3 {
					t.Errorf("len(Tags) = %v, want 3", len(content.Tags))
					return
				}
				if content.Tags[0].Name != "tag1" {
					t.Errorf("Tags[0].Name = %v, want tag1", content.Tags[0].Name)
				}
			},
		},
		{
			name: "with meta fields",
			form: ContentForm{
				Heading:         "Test",
				Body:            "Body",
				Description:     "Test Description",
				Keywords:        "test,keywords",
				TableOfContents: true,
			},
			checkFn: func(t *testing.T, content feat.Content) {
				if content.Meta.Description != "Test Description" {
					t.Errorf("Meta.Description = %v, want Test Description", content.Meta.Description)
				}
				if content.Meta.Keywords != "test,keywords" {
					t.Errorf("Meta.Keywords = %v, want test,keywords", content.Meta.Keywords)
				}
				if !content.Meta.TableOfContents {
					t.Errorf("Meta.TableOfContents = %v, want true", content.Meta.TableOfContents)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := ToFeatContent(tt.form)
			tt.checkFn(t, content)
		})
	}
}

func TestToContentForm(t *testing.T) {
	pubTime := time.Now()
	content := feat.Content{
		Heading:     "Test Heading",
		Body:        "Test Body",
		Draft:       true,
		Featured:    true,
		PublishedAt: &pubTime,
		Tags: []feat.Tag{
			{Name: "tag1"},
			{Name: "tag2"},
		},
		Meta: feat.Meta{
			Description:     "Test Description",
			Keywords:        "test,keywords",
			TableOfContents: true,
		},
	}
	content.ID = uuid.New()
	content.UserID = uuid.New()
	content.SectionID = uuid.New()

	req := httptest.NewRequest("GET", "/", nil)
	form := ToContentForm(req, content)

	if form.Heading != "Test Heading" {
		t.Errorf("Heading = %v, want Test Heading", form.Heading)
	}
	if form.Body != "Test Body" {
		t.Errorf("Body = %v, want Test Body", form.Body)
	}
	if !form.Draft {
		t.Errorf("Draft = %v, want true", form.Draft)
	}
	if !form.Featured {
		t.Errorf("Featured = %v, want true", form.Featured)
	}
	if form.Description != "Test Description" {
		t.Errorf("Description = %v, want Test Description", form.Description)
	}
	if form.Tags != "tag1,tag2" {
		t.Errorf("Tags = %v, want tag1,tag2", form.Tags)
	}
}

func TestContentFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      ContentForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: ContentForm{
				Heading: "Test",
			},
			wantValid: true,
		},
		{
			name:      "empty heading",
			form:      ContentForm{},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewContentForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestLayoutFormFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		formData url.Values
		wantErr  bool
		checkFn  func(t *testing.T, form LayoutForm)
	}{
		{
			name: "valid form",
			formData: url.Values{
				"name":        {"Test Layout"},
				"description": {"Test Description"},
				"code":        {"<div>Test</div>"},
			},
			wantErr: false,
			checkFn: func(t *testing.T, form LayoutForm) {
				if form.Name != "Test Layout" {
					t.Errorf("Name = %v, want Test Layout", form.Name)
				}
				if form.Description != "Test Description" {
					t.Errorf("Description = %v, want Test Description", form.Description)
				}
				if form.Code != "<div>Test</div>" {
					t.Errorf("Code = %v, want <div>Test</div>", form.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			form, err := LayoutFormFromRequest(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("LayoutFormFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFn != nil {
				tt.checkFn(t, form)
			}
		})
	}
}

func TestToFeatLayout(t *testing.T) {
	form := LayoutForm{
		Name:        "Test Layout",
		Description: "Test Description",
		Code:        "<div>Test</div>",
	}

	layout := ToFeatLayout(form)

	if layout.Name != "Test Layout" {
		t.Errorf("Name = %v, want Test Layout", layout.Name)
	}
	if layout.Description != "Test Description" {
		t.Errorf("Description = %v, want Test Description", layout.Description)
	}
	if layout.Code != "<div>Test</div>" {
		t.Errorf("Code = %v, want <div>Test</div>", layout.Code)
	}
}

func TestToLayoutForm(t *testing.T) {
	layout := feat.Layout{
		Name:        "Test Layout",
		Description: "Test Description",
		Code:        "<div>Test</div>",
	}
	layout.ID = uuid.New()

	req := httptest.NewRequest("GET", "/", nil)
	form := ToLayoutForm(req, layout)

	if form.Name != "Test Layout" {
		t.Errorf("Name = %v, want Test Layout", form.Name)
	}
	if form.Description != "Test Description" {
		t.Errorf("Description = %v, want Test Description", form.Description)
	}
}

func TestLayoutFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      LayoutForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: LayoutForm{
				Name: "Test",
				Code: "<div>Test</div>",
			},
			wantValid: true,
		},
		{
			name: "missing name",
			form: LayoutForm{
				Code: "<div>Test</div>",
			},
			wantValid: false,
		},
		{
			name: "missing code",
			form: LayoutForm{
				Name: "Test",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewLayoutForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestSectionFormFromRequest(t *testing.T) {
	layoutID := uuid.New()
	formData := url.Values{
		"name":        {"Test Section"},
		"description": {"Test Description"},
		"path":        {"/test"},
		"layout_id":   {layoutID.String()},
	}

	req := httptest.NewRequest("POST", "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := SectionFormFromRequest(req)
	if err != nil {
		t.Errorf("SectionFormFromRequest() error = %v", err)
		return
	}

	if form.Name != "Test Section" {
		t.Errorf("Name = %v, want Test Section", form.Name)
	}
	if form.Path != "/test" {
		t.Errorf("Path = %v, want /test", form.Path)
	}
}

func TestToFeatSection(t *testing.T) {
	layoutID := uuid.New()
	form := SectionForm{
		Name:        "Test Section",
		Description: "Test Description",
		Path:        "/test",
		LayoutID:    layoutID.String(),
	}

	section := ToFeatSection(form)

	if section.Name != "Test Section" {
		t.Errorf("Name = %v, want Test Section", section.Name)
	}
	if section.Path != "/test" {
		t.Errorf("Path = %v, want /test", section.Path)
	}
	if section.LayoutID != layoutID {
		t.Errorf("LayoutID = %v, want %v", section.LayoutID, layoutID)
	}
}

func TestToSectionForm(t *testing.T) {
	layoutID := uuid.New()
	section := feat.Section{
		Name:        "Test Section",
		Description: "Test Description",
		Path:        "/test",
		LayoutID:    layoutID,
	}
	section.ID = uuid.New()

	req := httptest.NewRequest("GET", "/", nil)
	form := ToSectionForm(req, section)

	if form.Name != "Test Section" {
		t.Errorf("Name = %v, want Test Section", form.Name)
	}
	if form.Path != "/test" {
		t.Errorf("Path = %v, want /test", form.Path)
	}
}

func TestSectionFormValidate(t *testing.T) {
	layoutID := uuid.New()
	tests := []struct {
		name      string
		form      SectionForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: SectionForm{
				Name:     "Test",
				Path:     "/test",
				LayoutID: layoutID.String(),
			},
			wantValid: true,
		},
		{
			name: "missing name",
			form: SectionForm{
				Path:     "/test",
				LayoutID: layoutID.String(),
			},
			wantValid: false,
		},
		{
			name: "missing path",
			form: SectionForm{
				Name:     "Test",
				LayoutID: layoutID.String(),
			},
			wantValid: false,
		},
		{
			name: "missing layout",
			form: SectionForm{
				Name: "Test",
				Path: "/test",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewSectionForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestTagFormFromRequest(t *testing.T) {
	formData := url.Values{
		"name": {"Test Tag"},
	}

	req := httptest.NewRequest("POST", "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := TagFormFromRequest(req)
	if err != nil {
		t.Errorf("TagFormFromRequest() error = %v", err)
		return
	}

	if form.Name != "Test Tag" {
		t.Errorf("Name = %v, want Test Tag", form.Name)
	}
}

func TestToFeatTag(t *testing.T) {
	form := TagForm{
		Name: "Test Tag",
	}

	tag := ToFeatTag(form)

	if tag.Name != "Test Tag" {
		t.Errorf("Name = %v, want Test Tag", tag.Name)
	}
}

func TestToTagForm(t *testing.T) {
	tag := feat.Tag{
		Name: "Test Tag",
	}
	tag.ID = uuid.New()

	req := httptest.NewRequest("GET", "/", nil)
	form := ToTagForm(req, tag)

	if form.Name != "Test Tag" {
		t.Errorf("Name = %v, want Test Tag", form.Name)
	}
}

func TestTagFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      TagForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: TagForm{
				Name: "Test",
			},
			wantValid: true,
		},
		{
			name:      "empty name",
			form:      TagForm{},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewTagForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestParamFormFromRequest(t *testing.T) {
	formData := url.Values{
		"name":        {"test_param"},
		"description": {"Test Description"},
		"value":       {"test_value"},
		"ref_key":     {"test_ref"},
	}

	req := httptest.NewRequest("POST", "/", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := ParamFormFromRequest(req)
	if err != nil {
		t.Errorf("ParamFormFromRequest() error = %v", err)
		return
	}

	if form.Name != "test_param" {
		t.Errorf("Name = %v, want test_param", form.Name)
	}
	if form.Value != "test_value" {
		t.Errorf("Value = %v, want test_value", form.Value)
	}
}

func TestToFeatParam(t *testing.T) {
	form := ParamForm{
		Name:        "test_param",
		Description: "Test Description",
		Value:       "test_value",
		RefKey:      "test_ref",
	}

	param := ToFeatParam(form)

	if param.Name != "test_param" {
		t.Errorf("Name = %v, want test_param", param.Name)
	}
	if param.Value != "test_value" {
		t.Errorf("Value = %v, want test_value", param.Value)
	}
	if param.RefKey != "test_ref" {
		t.Errorf("RefKey = %v, want test_ref", param.RefKey)
	}
}

func TestToParamForm(t *testing.T) {
	param := feat.Param{
		Name:        "test_param",
		Description: "Test Description",
		Value:       "test_value",
		RefKey:      "test_ref",
	}
	param.ID = uuid.New()

	req := httptest.NewRequest("GET", "/", nil)
	form := ToParamForm(req, param)

	if form.Name != "test_param" {
		t.Errorf("Name = %v, want test_param", form.Name)
	}
	if form.Value != "test_value" {
		t.Errorf("Value = %v, want test_value", form.Value)
	}
}

func TestParamFormValidate(t *testing.T) {
	tests := []struct {
		name      string
		form      ParamForm
		wantValid bool
	}{
		{
			name: "valid form",
			form: ParamForm{
				Name:  "test",
				Value: "value",
			},
			wantValid: true,
		},
		{
			name: "missing name",
			form: ParamForm{
				Value: "value",
			},
			wantValid: false,
		},
		{
			name: "missing value",
			form: ParamForm{
				Name: "test",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.form.BaseForm = NewParamForm(req).BaseForm
			tt.form.Validate()
			isValid := tt.form.Validation().IsValid()
			if isValid != tt.wantValid {
				t.Errorf("Validate() isValid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}
