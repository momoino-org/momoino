package usermgt_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
	"wano-island/common/core"
	"wano-island/common/usermgt"
	"wano-island/console/modules/httpsrv"
	mockcore "wano-island/testing/mocks/common/core"
	"wano-island/testing/testutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gorm.io/gorm"
)

var _ = Describe("handler.profile.go", func() {
	var (
		db     *gorm.DB
		config *mockcore.MockAppConfig
		router http.Handler
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		db, _ = testutils.CreateTestDBInstance()

		config = mockcore.NewMockAppConfig(GinkgoT())
		config.EXPECT().GetAppVersion().Return("1.0.0")
		config.EXPECT().GetRevision().Return("testing")
		config.EXPECT().GetMode().Return(core.TestingMode)
		config.EXPECT().IsTesting().Return(true)
		config.EXPECT().GetCorsConfig().Return(&core.CorsConfig{})

		router = testutils.CreateRouter(func(rp *httpsrv.RouteParams) {
			rp.Config = config
			rp.Routes = []core.HTTPRoute{
				usermgt.NewProfileHandler(usermgt.ProfileHandlerParams{
					Logger: core.NewNoopLogger(),
					DB:     db,
				}),
			}
		})
	})

	It("returns an error if there is no access token", func(ctx SpecContext) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		recorder := httptest.NewRecorder()

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusUnauthorized))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID":  Equal("E-0002"),
			"Message":    Equal("You must be authenticated to use this feature"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns a user profile if there is valid access token", func(ctx SpecContext) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, core.WithTestAuthUser(req, core.AuthenticatedUser{
			ID:          "019135f7-6265-7ef8-8920-57280736f6c0",
			Username:    "testing",
			Email:       "testing@example.com",
			GivenName:   "testing",
			FamilyName:  "testing",
			Locale:      "en",
			Roles:       []string{},
			Permissions: []string{},
		}))

		var response core.Response[usermgt.ProfileResponse]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusOK))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("S-0000"),
			"Message":   Equal("Success"),
			"Data": MatchFields(IgnoreMissing, Fields{
				"ID":        Equal("019135f7-6265-7ef8-8920-57280736f6c0"),
				"Username":  Equal("testing"),
				"Email":     Equal("testing@example.com"),
				"FirstName": Equal("testing"),
				"LastName":  Equal("testing"),
				"Locale":    Equal("en"),
			}),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})
})
