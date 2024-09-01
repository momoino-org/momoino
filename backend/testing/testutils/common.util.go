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
	return WithFxLifeCycle(func(l *fxtest.Lifecycle) http.Handler {
		routeParams := httpsrv.RouteParams{
			Logger: core.NewNoopLogger(),
			I18nBundle: core.NewI18nBundle(core.I18nBundleParams{
				AppLifeCycle: l,
				LocaleFS:     GetResourceFS(),
			}),
		}

		opt(&routeParams)

		return httpsrv.NewRouter(routeParams)
	})
}

// WithFakeJWT adds a JWT token to the given HTTP request header.
// The function sets the "Authorization" header with a Bearer token containing a sample JWT.
// This token is used for testing purposes and should not be used in production.
func WithFakeJWT(r *http.Request) *http.Request {
	//nolint:lll // No need to fix
	r.Header.Add(core.AuthorizationHeader, "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMTkxMzVmNy02MjY1LTdlZjgtODkyMC01NzI4MDczNmY2YzAiLCJleHAiOjE3MjQ4NDc3MTksIm5iZiI6MTcyNDg0NzcxOSwiaWF0IjoxNzI0ODQ3NzE5LCJyb2xlcyI6W10sInBlcm1pc3Npb25zIjpbXX0.I9-Kr2ArmW3V-eUN9KKxKShmV9oDWefKBzaXo5BJCqV6fqVtddNFSxnmGzj72WMykCXSTrz92NDGtH8M-lZWwBsNOJY7XCZFoDdYKHk_OyGR9Nk-lRvburgMgaNChw6lD-zjZTb2xfJhmdj4IMbZOcDMB6bdo5bAz_M_3iiPw1gMX9Jkd5yXIwchjOWwVasVO0ycZZ3qFz-mBrSn1FyG8T_ox6avcEHFdiDiBUR6YBaXZwIpiFqhy0aDdvGz8MCvT95b5keTO6jcNLwHZrm1YnZD-lPz5xJQL14n-FnKOvi0UVpEbmkkmyfQz4IH5kdzaRaEdHEYsSyjpNJ1Xaq5lA")
	return r
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
