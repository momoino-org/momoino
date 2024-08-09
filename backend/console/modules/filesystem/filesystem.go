package filesystem

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"wano-island/common/core"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type createFileSystemHandler struct {
}

var _ core.HTTPRoute = (*createFileSystemHandler)(nil)

func newCreateFileSystemHandler() *createFileSystemHandler {
	return &createFileSystemHandler{}
}

func (handler *createFileSystemHandler) Pattern() string {
	return "GET /static/*"
}

func (handler *createFileSystemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
	fs.ServeHTTP(w, r)
}

func NewFileSystemModule() fx.Option {
	return fx.Module(
		"File System Module",
		fx.Provide(core.AsRoute(newCreateFileSystemHandler)),
	)
}
