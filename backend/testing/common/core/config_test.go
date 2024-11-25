package core_test

import (
	"net/http"
	"wano-island/common/core"
	"wano-island/testing/testutils"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.uber.org/fx"
)

var _ = Describe("[common/core/config.go]", Ordered, func() {
	var (
		httpClient *resty.Client
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		testutils.ConfigureMinimumEnvVariables()

		httpClient = core.NewHTTPClient()
		httpmock.ActivateNonDefault(httpClient.GetClient())

		httpmock.RegisterResponder(
			"GET",
			"http://keycloak.test/.well-known/openid-configuration",
			httpmock.NewJsonResponderOrPanic(
				http.StatusOK,
				map[string]string{
					"issuer":   "http://localhost:8080/realms/momoino",
					"jwks_uri": "http://localhost:8080/realms/momoino/protocol/openid-connect/certs",
				}))

		DeferCleanup(func() {
			httpmock.DeactivateAndReset()
		})
	})

	Context("when initializing the config module", func() {
		It("should return an fx.Option", func() {
			Expect(core.NewConfigModule()).To(BeAssignableToTypeOf(fx.Module("")))
		})
	})

	Context("when the config not set in the environment", func() {
		It("should return default values when no configuration is provided", func() {
			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg).NotTo(BeNil())

			Expect(appCfg.GetAppVersion()).To(Equal(""))
			Expect(appCfg.GetCompatibleVersion()).To(Equal(""))
			Expect(appCfg.GetRevision()).To(Equal(""))
			Expect(appCfg.GetMode()).To(Equal("testing"))
			Expect(appCfg.IsDevelopment()).To(BeFalse())
			Expect(appCfg.IsProduction()).To(BeFalse())
			Expect(appCfg.IsTesting()).To(BeTrue())
			Expect(appCfg.GetDatabaseConfig()).To(PointTo(MatchFields(IgnoreMissing, Fields{
				"Host":         Equal("localhost"),
				"Port":         Equal(uint16(5432)),
				"DatabaseName": Equal("momoiro-wano"),
				"Username":     Equal("root"),
				"Password":     Equal("Keep!t5ecret"),
				"MaxAttempts":  Equal(uint(3)),
			})))
			Expect(appCfg.GetCorsConfig()).To(PointTo(MatchFields(IgnoreMissing, Fields{
				"AllowedOrigins":   Equal([]string{"*"}),
				"AllowedMethods":   Equal([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
				"AllowedHeaders":   BeEmpty(),
				"ExposedHeaders":   BeEmpty(),
				"AllowCredentials": BeFalse(),
				"MaxAge":           BeZero(),
			})))
		})
	})

	Context("when APP_MODE is set in the environment", func() {
		DescribeTable("should return value of APP_MODE if it is valid",
			func(mode string) {
				t := GinkgoT()
				t.Setenv("APP_MODE", mode)

				appCfg, err := core.NewAppConfig(httpClient)
				Expect(err).NotTo(HaveOccurred())
				Expect(appCfg.GetMode()).To(Equal(mode))
			},
			Entry(`should return "production"`, "production"),
			Entry(`should return "development"`, "development"),
			Entry(`should return "testing"`, "testing"))

		It("should return default value if the value of APP_MODE is invalid", func() {
			t := GinkgoT()
			t.Setenv("APP_MODE", "staging")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetMode()).To(Equal("development"))
		})

		It("IsDevelopment() should return true when APP_MODE is development", func() {
			t := GinkgoT()
			t.Setenv("APP_MODE", "development")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.IsDevelopment()).To(BeTrue())
		})

		It("IsProduction() should return true when APP_MODE is production", func() {
			t := GinkgoT()
			t.Setenv("APP_MODE", "production")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.IsProduction()).To(BeTrue())
		})

		It("IsTesting() should return true when APP_MODE is testing", func() {
			t := GinkgoT()
			t.Setenv("APP_MODE", "testing")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.IsTesting()).To(BeTrue())
		})
	})

	Context("when AppVesion is set", func() {
		It("should return the value", func() {
			core.AppVersion = "1.0.0"
			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetAppVersion()).To(Equal("1.0.0"))
		})
	})

	Context("when CompatibleVersion is set", func() {
		It("should return the value", func() {
			core.CompatibleVersion = "0.5.0"
			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetCompatibleVersion()).To(Equal("0.5.0"))
		})
	})

	Context("when AppRevision is set", func() {
		It("should return the value", func() {
			core.AppRevision = "12345678"
			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetRevision()).To(Equal("12345678"))
		})
	})

	Context("when user defines the db configuration in the envionment", func() {
		It("should use the db configuration in the envionment", func() {
			t := GinkgoT()
			t.Setenv("APP_DATABASE_HOST", "127.0.0.1")
			t.Setenv("APP_DATABASE_PORT", "5434")
			t.Setenv("APP_DATABASE_NAME", "testing-db")
			t.Setenv("APP_DATABASE_USERNAME", "admin")
			t.Setenv("APP_DATABASE_PASSWORD", "password")
			t.Setenv("APP_DATABASE_MAX_ATTEMPTS", "5")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetDatabaseConfig()).To(PointTo(MatchFields(IgnoreMissing, Fields{
				"Host":         Equal("127.0.0.1"),
				"Port":         Equal(uint16(5434)),
				"DatabaseName": Equal("testing-db"),
				"Username":     Equal("admin"),
				"Password":     Equal("password"),
				"MaxAttempts":  Equal(uint(5)),
			})))
		})
	})

	Context("when user defines the cors configuration in the envionment", func() {
		It("should use the cors configuration in the envionment", func() {
			t := GinkgoT()
			t.Setenv("APP_CORS_ALLOWED_ORIGINS", "http://localhost:3000 http://*.localhost")
			t.Setenv("APP_CORS_ALLOWED_METHODS", "get head post")
			t.Setenv("APP_CORS_ALLOWED_HEADERS", "X-Header-1 X-Header-2  X-Header-3")
			t.Setenv("APP_CORS_EXPOSED_HEADERS", "X-Exposed-Header-1  Y-Exposed-Header-2  ")
			t.Setenv("APP_CORS_ALLOW_CREDENTIALS", "true")
			t.Setenv("APP_CORS_MAX_AGE", "500")

			appCfg, err := core.NewAppConfig(httpClient)
			Expect(err).NotTo(HaveOccurred())
			Expect(appCfg.GetCorsConfig()).To(PointTo(MatchFields(IgnoreMissing, Fields{
				"AllowedOrigins":   Equal([]string{"http://localhost:3000", "http://*.localhost"}),
				"AllowedMethods":   Equal([]string{"get", "head", "post"}),
				"AllowedHeaders":   Equal([]string{"X-Header-1", "X-Header-2", "X-Header-3"}),
				"ExposedHeaders":   Equal([]string{"X-Exposed-Header-1", "Y-Exposed-Header-2"}),
				"AllowCredentials": BeTrue(),
				"MaxAge":           Equal(500),
			})))
		})
	})
})
