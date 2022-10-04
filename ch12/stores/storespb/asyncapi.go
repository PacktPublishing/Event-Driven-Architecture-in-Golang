package storespb

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed index.html
//go:embed css/*
//go:embed js/*
var asyncAPI embed.FS

func RegisterAsyncAPI(mux *chi.Mux) error {
	const specRoot = "/stores-asyncapi/"

	// mount the swagger specifications
	mux.Mount(specRoot, http.StripPrefix(specRoot, http.FileServer(http.FS(asyncAPI))))

	return nil
}
