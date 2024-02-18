package main

import (
	"embed"
	"html/template"
	"io/fs"
)

const (
	layoutsDir   = "view/layouts"
	templatesDir = "view"
	extension    = "/*.html"
)

var (
	// Template variables
	//go:embed view/*
	files     embed.FS
	templates map[string]*template.Template
)

func loadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}

	return nil
}
