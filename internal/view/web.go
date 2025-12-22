package view

import (
	"embed"
	"html/template"
	"net/http"
)

var Templates = template.Must(template.ParseFS(templatesFS, "templates/*.html"))
var StaticHandler = http.FileServerFS(staticFS)

//go:embed static/**
var staticFS embed.FS

//go:embed templates/*.html
var templatesFS embed.FS
