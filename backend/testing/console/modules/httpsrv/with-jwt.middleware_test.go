package httpsrv_test

import (
	"net/http"
	"net/http/httptest"
	"testing/fstest"
	"wano-island/common/core"
	"wano-island/console/modules/httpsrv"

	"github.com/go-chi/chi/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	"go.uber.org/fx/fxtest"
)

var _ = Describe("withJwtMiddleware", func() {
	AfterEach(func() {
		Eventually(Goroutines).ShouldNot(HaveLeaked())
	})

	DescribeTable(
		"Localizer",
		func(path string, headers map[string]string, expectedResult string) {
			router := chi.NewRouter()
			appLifeCycle := fxtest.NewLifecycle(GinkgoT())
			i18nBundle := core.NewI18nBundle(core.I18nBundleParams{
				AppLifeCycle: appLifeCycle,
				LocaleFS: fstest.MapFS{
					"resources/trans/locale.en.yaml": &fstest.MapFile{
						Data: []byte("Language: English"),
					},
					"resources/trans/locale.vi.yaml": &fstest.MapFile{
						Data: []byte("Language: Vietnamese"),
					},
				},
			})

			router.Use(httpsrv.WithI18nMiddleware(i18nBundle))
			router.Use(httpsrv.WithJwtMiddleware(i18nBundle, core.NewNoopLogger()))

			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				localizer := core.GetLocalizer(r)
				_, _ = w.Write([]byte(localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "Language",
					},
				})))
			})

			appLifeCycle.RequireStart()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			for key, value := range headers {
				request.Header.Add(key, value)
			}

			router.ServeHTTP(recorder, request)

			Expect(recorder.Body.String()).To(Equal(expectedResult))
		},
		Entry(
			"should use the preferred language of the logged-in user if the request does not contain the lang query parameter",
			"/",
			map[string]string{
				//nolint:lll // No need to fix
				core.AuthorizationHeader: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMTkxYTM2NS04M2ViLTc1ZGYtYjMzMy0wZDdkNmEwYmNiZTYiLCJleHAiOjE3MjUwMjMzMzAsIm5iZiI6MTcyNTAyMzMzMCwiaWF0IjoxNzI1MDIzMzMwLCJsb2NhbGUiOiJ2aSIsInJvbGVzIjpbXSwicGVybWlzc2lvbnMiOltdfQ.C88NVUYWOR5KTltSpqw7-TmD57Rpia8n5IA8iKUueMYLNHiho0Q3fqRMOnQkai3eMrkvKa9hKIbLt9OJ8ndbbdO3hDMXeN0zS4-kPWJz6nXFqCcJrFqA1JSzsMs7470K_I7tUo0OUWqRkzbQttvHCUrWGrvO5FhyMyP3nof_3ciwXfkXj-yFWSl5Yu-8isU3EGKxKOeqtFDEkTmzN2pM0QR_jGhrPIABmg6up1zHxbX_lvY3uyhf4TgRzbIDbjFfVNw3TEi4RK7HQeCE92l8A4eZRa4q0qGezPeQsEG5UtJ9ZcYReNdJjX_0OJ1DJaYjWxLijvibQojRZqsfTMEInA",
			},
			"Vietnamese",
		),
		Entry(
			//nolint:lll // This is a test case name, no need to fix
			`should always display the translated message based on the lang query parameter, even if the user is logged in with a different preferred language`,
			"/?lang=en",
			map[string]string{
				//nolint:lll // No need to fix
				core.AuthorizationHeader: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMTkxYTM2NS04M2ViLTc1ZGYtYjMzMy0wZDdkNmEwYmNiZTYiLCJleHAiOjE3MjUwMjMzMzAsIm5iZiI6MTcyNTAyMzMzMCwiaWF0IjoxNzI1MDIzMzMwLCJsb2NhbGUiOiJ2aSIsInJvbGVzIjpbXSwicGVybWlzc2lvbnMiOltdfQ.C88NVUYWOR5KTltSpqw7-TmD57Rpia8n5IA8iKUueMYLNHiho0Q3fqRMOnQkai3eMrkvKa9hKIbLt9OJ8ndbbdO3hDMXeN0zS4-kPWJz6nXFqCcJrFqA1JSzsMs7470K_I7tUo0OUWqRkzbQttvHCUrWGrvO5FhyMyP3nof_3ciwXfkXj-yFWSl5Yu-8isU3EGKxKOeqtFDEkTmzN2pM0QR_jGhrPIABmg6up1zHxbX_lvY3uyhf4TgRzbIDbjFfVNw3TEi4RK7HQeCE92l8A4eZRa4q0qGezPeQsEG5UtJ9ZcYReNdJjX_0OJ1DJaYjWxLijvibQojRZqsfTMEInA",
			},
			"English",
		),
	)
})
