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

var _ = Describe("withI18nMiddleware", func() {
	AfterEach(func() {
		Eventually(Goroutines).ShouldNot(HaveLeaked())
	})

	DescribeTable(
		"Localizer",
		func(path string, expectedResult string) {
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
			router.ServeHTTP(recorder, request)

			Expect(recorder.Body.String()).To(Equal(expectedResult))
		},
		Entry(
			//nolint:lll // This is a test case name, no need to fix
			`should return the default message in English when the user is not logged in and when the "lang" query parameter is not provided`,
			"/",
			"English",
		),
		Entry(
			`should return the translated message according to the language specified in the "lang" query parameter`,
			"/?lang=vi",
			"Vietnamese",
		),
		Entry(
			`should return the English message if the specified lang in the query parameter is not supported`,
			"/?lang=unknown",
			"English",
		),
	)
})
