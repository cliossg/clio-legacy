package fake

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/hermesgen/hm"
)

type TemplateManager struct {
	hm.Core
	templates       map[string]*template.Template
	templatesPath   string
	customFuncMap   template.FuncMap
}

func NewTemplateManager() *TemplateManager {
	cfg := hm.NewConfig()
	params := hm.XParams{Cfg: cfg}
	core := hm.NewCore("fake-template-manager", params)

	return &TemplateManager{
		Core:          core,
		templates:     make(map[string]*template.Template),
		templatesPath: "assets/template",
		customFuncMap: make(template.FuncMap),
	}
}

func (tm *TemplateManager) Get(handler, action string) (*template.Template, error) {
	key := handler + ":" + action

	if tmpl, ok := tm.templates[key]; ok {
		return tmpl, nil
	}

	tmpl, err := tm.loadTemplate(handler, action)
	if err != nil || tmpl.Lookup("page") == nil {
		funcMap := make(template.FuncMap)
		for k, v := range tm.customFuncMap {
			funcMap[k] = v
		}
		tmpl = template.Must(template.New("page").Funcs(funcMap).Parse(`<div>{{template "content" .}}</div>`))
		template.Must(tmpl.New("content").Parse(`<div>test</div>`))
	}

	tm.templates[key] = tmpl
	return tmpl, nil
}

func (tm *TemplateManager) RegisterFunctions(funcs template.FuncMap) {
	for name, fn := range funcs {
		tm.customFuncMap[name] = fn
	}
}

func (tm *TemplateManager) loadTemplate(handler, action string) (*template.Template, error) {
	layoutPath := filepath.Join(tm.templatesPath, "layout", "layout.tmpl")
	handlerPath := filepath.Join(tm.templatesPath, "handler", handler, action+".tmpl")

	funcMap := make(template.FuncMap)
	for k, v := range tm.customFuncMap {
		funcMap[k] = v
	}

	tmpl := template.New("page").Funcs(funcMap)
	parsed := false

	if _, err := os.Stat(layoutPath); err == nil {
		tmpl, err = tmpl.ParseFiles(layoutPath)
		if err != nil {
			return nil, err
		}
		parsed = true
	}

	if _, err := os.Stat(handlerPath); err == nil {
		tmpl, err = tmpl.ParseFiles(handlerPath)
		if err != nil {
			return nil, err
		}
		parsed = true
	}

	if !parsed {
		return nil, os.ErrNotExist
	}

	return tmpl, nil
}
