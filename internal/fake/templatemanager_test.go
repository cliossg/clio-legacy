package fake

import (
	"bytes"
	"html/template"
	"testing"
)

func TestNewTemplateManager(t *testing.T) {
	tm := NewTemplateManager()
	if tm == nil {
		t.Fatal("NewTemplateManager() returned nil")
	}
}

func TestTemplateManagerGet(t *testing.T) {
	tests := []struct {
		name         string
		feature      string
		templateName string
		wantErr      bool
	}{
		{
			name:         "returns template for ssg feature",
			feature:      "ssg",
			templateName: "new-tag",
			wantErr:      false,
		},
		{
			name:         "returns template for auth feature",
			feature:      "auth",
			templateName: "login",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTemplateManager()

			tmpl, err := tm.Get(tt.feature, tt.templateName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateManager.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tmpl == nil {
				t.Error("TemplateManager.Get() returned nil template")
			}
		})
	}
}

func TestTemplateManagerExecute(t *testing.T) {
	tm := NewTemplateManager()
	tmpl, err := tm.Get("ssg", "test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	if buf.Len() == 0 {
		t.Error("Execute() produced no output")
	}
}

func TestTemplateManagerCaching(t *testing.T) {
	tm := NewTemplateManager()

	tmpl1, err := tm.Get("ssg", "test")
	if err != nil {
		t.Fatalf("First Get() error = %v", err)
	}

	tmpl2, err := tm.Get("ssg", "test")
	if err != nil {
		t.Fatalf("Second Get() error = %v", err)
	}

	if tmpl1 != tmpl2 {
		t.Error("Get() should return same template instance for same key")
	}
}

func TestTemplateManagerRegisterFunctions(t *testing.T) {
	tests := []struct {
		name      string
		funcs     template.FuncMap
		wantFuncs int
	}{
		{
			name: "registers single function",
			funcs: template.FuncMap{
				"upper": func(s string) string { return s },
			},
			wantFuncs: 1,
		},
		{
			name: "registers multiple functions",
			funcs: template.FuncMap{
				"upper": func(s string) string { return s },
				"lower": func(s string) string { return s },
			},
			wantFuncs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := NewTemplateManager()
			tm.RegisterFunctions(tt.funcs)

			if len(tm.customFuncMap) != tt.wantFuncs {
				t.Errorf("expected %d functions, got %d", tt.wantFuncs, len(tm.customFuncMap))
			}

			for name := range tt.funcs {
				if _, ok := tm.customFuncMap[name]; !ok {
					t.Errorf("function %q not registered", name)
				}
			}
		})
	}
}
