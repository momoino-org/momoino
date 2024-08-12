package filesystem

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type createFileSystemHandler struct {
	staticFiles embed.FS
}

var _ core.HTTPRoute = (*createFileSystemHandler)(nil)

func newCreateFileSystemHandler(staticFiles embed.FS) *createFileSystemHandler {
	return &createFileSystemHandler{
		staticFiles: staticFiles,
	}
}

func (handler *createFileSystemHandler) Pattern() string {
	return "GET /static/*"
}

func (handler *createFileSystemHandler) IsPrivateRoute() bool {
	return false
}

func (handler *createFileSystemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fSys, _ := fs.Sub(handler.staticFiles, "static")
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(fSys)))
	fs.ServeHTTP(w, r)
}

func NewFileSystemModule(staticFiles embed.FS) fx.Option {
	return fx.Module(
		"File System Module",
		fx.Provide(core.AsRoute(func() *createFileSystemHandler {
			return newCreateFileSystemHandler(staticFiles)
		})),
	)
}
