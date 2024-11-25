package testutils

import (
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing/fstest"
	"wano-island/common/core"
	"wano-island/console/modules/httpsrv"

	. "github.com/onsi/ginkgo/v2"
	"go.uber.org/fx/fxtest"
)

// getWorkspaceDir retrieves the workspace directory path by executing the "go env GOWORK" command.
func getWorkspaceDir() string {
	cmd := exec.Command("go", "env", "GOWORK")
	output, _ := cmd.Output()

	return filepath.Dir(string(output))
}

// GetResourceFS creates a virtual file system (fs.FS).
func GetResourceFS() fs.FS {
	localeFile := filepath.Join(getWorkspaceDir(), "console", "resources", "trans", "locale.en.yaml")
	localeEnYamlData, _ := os.ReadFile(localeFile)

	return fstest.MapFS{
		"resources/trans/locale.en.yaml": &fstest.MapFile{
			Data: localeEnYamlData,
		},
	}
}

// CreateRouter sets up a test router for HTTP requests.
// It takes a function `opt` as a parameter, which is used to configure the route parameters.
func CreateRouter(opt func(*httpsrv.RouteParams)) http.Handler {
	i18nBundle, _ := core.NewI18nBundle(core.I18nBundleParams{
		LocaleFS: GetResourceFS(),
	})

	routeParams := httpsrv.RouteParams{
		Logger:     core.NewNoopLogger(),
		I18nBundle: i18nBundle,
	}

	opt(&routeParams)

	return httpsrv.NewRouter(routeParams)
}

// WithFxLifeCycle is a helper function that creates and manages an fxtest.Lifecycle instance.
// This function is useful for testing applications built using the fx framework.
//
//nolint:ireturn // We doesn't know the returned type, so it is okay to use generic type
func WithFxLifeCycle[T any](fn func(*fxtest.Lifecycle) T) T {
	appLifeCycle := fxtest.NewLifecycle(GinkgoT())

	result := fn(appLifeCycle)

	appLifeCycle.RequireStart()

	DeferCleanup(func() {
		appLifeCycle.RequireStop()
	})

	return result
}
