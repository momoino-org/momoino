package httpsrv_test

import (
	"net/http"
	"net/http/httptest"
	"testing/fstest"
	"wano-island/common/core"
	"wano-island/console/modules/httpsrv"
	mockcore "wano-island/testing/mocks/common/core"
	"wano-island/testing/testutils"

	"github.com/go-chi/chi/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("withJwtMiddleware", func() {
	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
	})

	DescribeTable(
		"Localizer",
		func(path string, headers map[string]string, expectedResult string) {
			router := chi.NewRouter()
			i18nBundle, _ := core.NewI18nBundle(core.I18nBundleParams{
				LocaleFS: fstest.MapFS{
					"resources/trans/locale.en.yaml": &fstest.MapFile{
						Data: []byte("Language: English"),
					},
					"resources/trans/locale.vi.yaml": &fstest.MapFile{
						Data: []byte("Language: Vietnamese"),
					},
				},
			})

			config := mockcore.NewMockAppConfig(GinkgoT())
			config.EXPECT().GetJWTConfig().Return(testutils.GetJWTConfig())

			router.Use(httpsrv.WithI18nMiddleware(i18nBundle))
			router.Use(httpsrv.WithJwtMiddleware(i18nBundle, config, core.NewNoopLogger()))

			router.Get("/", func(w http.ResponseWriter, r *http.Request) {
				localizer := core.GetLocalizer(r)
				_, _ = w.Write([]byte(localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "Language",
					},
				})))
			})

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
				core.AuthorizationHeader: testutils.GenerateJWT(func(jc *core.JWTCustomClaims) {
					jc.Locale = "vi"
				}),
			},
			"Vietnamese",
		),
		Entry(
			//nolint:lll // This is a test case name, no need to fix
			`should always display the translated message based on the lang query parameter, even if the user is logged in with a different preferred language`,
			"/?lang=en",
			map[string]string{
				core.AuthorizationHeader: testutils.GenerateJWT(),
			},
			"English",
		),
	)
})
