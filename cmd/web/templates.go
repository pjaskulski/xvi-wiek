package main

import (
	"html/template"
	"path/filepath"
	"strings"
)

func createTemplateCache(path string) (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	var functions = template.FuncMap{
		"dayFact":     dayFact,
		"firstLetter": firstLetter,
	}

	templateFiles, err := filepath.Glob(filepath.Join(path, "*.page.gohtml"))
	if err != nil {
		return nil, err
	}

	for _, tFile := range templateFiles {
		name := filepath.Base(tFile)
		var ts *template.Template

		ts, err = template.New(name).Funcs(functions).ParseFiles(tFile)

		//ts, err = template.ParseFiles(tFile)
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

// firstLetter - helper func, return first letter of word
func firstLetter(word string) string {
	var result string = ""

	if len(word) > 0 {
		result = strings.ToLower(word)[:1]
	}

	return result
}
