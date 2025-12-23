package assets

import (
	"embed"
	"net/http"
)

var StaticHandler = http.FileServerFS(staticFS)

//go:embed static/**
var staticFS embed.FS
