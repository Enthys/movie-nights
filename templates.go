package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

const (
	layoutsDir   = "view/layouts"
	templatesDir = "view"
	extension    = "/*.html"
)

var (
	// -- go:embed view/*
	files     embed.FS
	templates map[string]*template.Template
	tplFunc   template.FuncMap = template.FuncMap{}
)

func loadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}

	for _, tpl := range tplFiles {
		if tpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}

		templates[tpl.Name()] = pt.Funcs(tplFunc)
	}

	return nil
}

type TplUserData struct {
	Authenticated bool
	Name          string
	Data          interface{}
}

func render(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	t, ok := templates[templateName]
	if !ok {
		log.Panic("Template " + templateName + " not found")
	}

	userData, err := store.Get(r, sessionCookieKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to retrieve session"))
		return
	}

	tplData := TplUserData{}

	if userData.Values[sk_authenticated] == true {
		tplData.Authenticated = true
		tplData.Name = userData.Values[sk_name].(string)
	}

	tplData.Data = data

	fmt.Printf("%+v", tplData)

	if err := t.Execute(w, tplData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error": "could not execute. %s"}`, err)))
		return
	}
}
