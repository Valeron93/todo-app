package view

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

var StaticHandler = http.FileServerFS(staticFS)

//go:embed static/**
var staticFS embed.FS

//go:embed templates/**
var templatesFS embed.FS

var pagesTemplates = make(map[string]*template.Template)
var partialTemplate = template.Must(template.ParseFS(templatesFS, "templates/partials/*.html"))

func init() {

	pages, _ := fs.Glob(templatesFS, "templates/pages/*.html")

	for _, pagePath := range pages {
		// extract html file name without extension
		page := path.Base(pagePath)
		page = strings.TrimSuffix(page, path.Ext(page))

		t := template.Must(
			template.ParseFS(
				templatesFS,
				"templates/layout.html",
				pagePath,
				"templates/partials/*.html",
			),
		)

		pagesTemplates[page] = t
	}

}

func renderPage(w http.ResponseWriter, page string, data any) error {

	t, ok := pagesTemplates[page]
	if !ok {
		return fmt.Errorf("view: no such page %#v", page)
	}

	return t.ExecuteTemplate(w, "layout", data)
}

func RenderPartial(w http.ResponseWriter, partial string, data any) error {
	return partialTemplate.ExecuteTemplate(w, partial, data)
}
