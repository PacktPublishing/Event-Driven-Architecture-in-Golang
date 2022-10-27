package rest

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed index.html
//go:embed api.swagger.json
var swaggerUI embed.FS

func RegisterSwagger(mux *chi.Mux) error {
	const specRoot = "/baskets-spec/"

	// mount the swagger specification
	mux.Mount(specRoot, http.StripPrefix(specRoot, http.FileServer(http.FS(swaggerUI))))

	return nil
}
