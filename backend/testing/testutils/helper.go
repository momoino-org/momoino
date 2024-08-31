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

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	"go.uber.org/fx/fxtest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DetectLeakyGoroutines is a helper function that detects and prevents leaked goroutines during testing.
// It uses the gleak library to track goroutine leaks and asserts that no goroutines are leaked.
// This function should be called at the beginning of each test case that involves goroutines.
// It uses ginkgo's DeferCleanup function to ensure that the cleanup code is executed after each test case.
func DetectLeakyGoroutines() {
	// Capture the initial set of goroutines
	nonLeakyGoroutines := Goroutines()

	// Register a cleanup function to assert that no goroutines are leaked
	DeferCleanup(func() {
		// Assert that no goroutines are leaked
		Eventually(Goroutines).ShouldNot(HaveLeaked(nonLeakyGoroutines))
	})
}

// CreateTestDBInstance creates a test database instance using sqlmock and gorm.DB.
// It returns a pointer to the gorm.DB instance and a sqlmock.Sqlmock instance for mocking database operations.
// The function also includes a cleanup mechanism using DeferCleanup to ensure that the database connection is closed
// and that all expectations set on the mock are met.
func CreateTestDBInstance() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	Expect(err).NotTo(HaveOccurred())

	gormDB, err := core.OpenDatabase(
		core.NewNoopLogger(),
		func(postgresCfg *postgres.Config, gormCfg *gorm.Config) {
			postgresCfg.Conn = db
			// Need to disable PrepareStmt to make mocking queries easier
			gormCfg.PrepareStmt = false
		},
	)
	Expect(err).NotTo(HaveOccurred())

	DeferCleanup(func() {
		mock.ExpectClose()

		sqlDB, getDBErr := gormDB.DB()
		Expect(getDBErr).NotTo(HaveOccurred())

		closeDBErr := sqlDB.Close()
		Expect(closeDBErr).NotTo(HaveOccurred())

		Expect(mock.ExpectationsWereMet()).NotTo(HaveOccurred())
	})

	return gormDB, mock
}

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

// SetupTestRouter sets up a test router for HTTP requests.
// It takes a function `opt` as a parameter, which is used to configure the route parameters.
func SetupTestRouter(opt func(*httpsrv.RouteParams)) http.Handler {
	appLifeCycle := fxtest.NewLifecycle(GinkgoT())
	routeParams := httpsrv.RouteParams{
		Logger: core.NewNoopLogger(),
		I18nBundle: core.NewI18nBundle(core.I18nBundleParams{
			AppLifeCycle: appLifeCycle,
			LocaleFS:     GetResourceFS(),
		}),
	}

	opt(&routeParams)

	router := httpsrv.NewRouter(routeParams)

	appLifeCycle.RequireStart()

	DeferCleanup(func() {
		appLifeCycle.RequireStop()
	})

	return router
}

// SetJWTToken adds a JWT token to the given HTTP request header.
// The function sets the "Authorization" header with a Bearer token containing a sample JWT.
// This token is used for testing purposes and should not be used in production.
func SetJWTToken(r *http.Request) {
	//nolint:lll // No need to fix
	r.Header.Add(core.AuthorizationHeader, "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMTkxMzVmNy02MjY1LTdlZjgtODkyMC01NzI4MDczNmY2YzAiLCJleHAiOjE3MjQ4NDc3MTksIm5iZiI6MTcyNDg0NzcxOSwiaWF0IjoxNzI0ODQ3NzE5LCJyb2xlcyI6W10sInBlcm1pc3Npb25zIjpbXX0.I9-Kr2ArmW3V-eUN9KKxKShmV9oDWefKBzaXo5BJCqV6fqVtddNFSxnmGzj72WMykCXSTrz92NDGtH8M-lZWwBsNOJY7XCZFoDdYKHk_OyGR9Nk-lRvburgMgaNChw6lD-zjZTb2xfJhmdj4IMbZOcDMB6bdo5bAz_M_3iiPw1gMX9Jkd5yXIwchjOWwVasVO0ycZZ3qFz-mBrSn1FyG8T_ox6avcEHFdiDiBUR6YBaXZwIpiFqhy0aDdvGz8MCvT95b5keTO6jcNLwHZrm1YnZD-lPz5xJQL14n-FnKOvi0UVpEbmkkmyfQz4IH5kdzaRaEdHEYsSyjpNJ1Xaq5lA")
}
