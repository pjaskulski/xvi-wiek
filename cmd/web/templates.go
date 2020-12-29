package main

import (
	"html/template"
	"path/filepath"
)

func createTemplateCache(path string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	templateFiles, err := filepath.Glob(filepath.Join(path, "*.page.tmpl.html"))
	if err != nil {
		return nil, err
	}

	for _, tFile := range templateFiles {
		name := filepath.Base(tFile)

		ts, err := template.ParseFiles(tFile)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(filepath.Join(path, "base.layout.tmpl.html"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}

	return cache, nil
}
