package main

import (
	"html/template"
	"path/filepath"
)

func createTemplateCache(path string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	templateFiles, err := filepath.Glob(filepath.Join(path, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	for _, tFile := range templateFiles {
		name := filepath.Base(tFile)
		var ts *template.Template

		ts, err = template.ParseFiles(tFile)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(filepath.Join(path, "base.layout.gohtml"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
