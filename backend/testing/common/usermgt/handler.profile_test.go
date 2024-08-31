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
	mockusermgt "wano-island/testing/mocks/common/usermgt"
	"wano-island/testing/testutils"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

var _ = Describe("handler.profile.go", func() {
	var (
		db             *gorm.DB
		appLifeCycle   *fxtest.Lifecycle
		config         *mockcore.MockAppConfig
		userRepository *mockusermgt.MockUserRepository
		handler        http.Handler
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		db, _ = testutils.CreateTestDBInstance()
		appLifeCycle = fxtest.NewLifecycle(GinkgoT())
		config = mockcore.NewMockAppConfig(GinkgoT())
		config.EXPECT().GetAppVersion().Return("1.0.0-testing")
		config.EXPECT().GetRevision().Return("testing")
		config.EXPECT().GetMode().Return(core.TestingMode)
		config.EXPECT().IsTesting().Return(true)
		// userRepository = usermgt.NewUserRepository(usermgt.UserRepositoryParams{})
		userRepository = mockusermgt.NewMockUserRepository(GinkgoT())
		// userRepository.EXPECT().FindUserByID(mock.Anything, mock.Anything, mock.Anything)
		logger := core.NewStdoutLogger(config)
		handler = httpsrv.NewRouter(httpsrv.RouteParams{
			Config: config,
			Logger: logger,
			I18nBundle: core.NewI18nBundle(core.I18nBundleParams{
				AppLifeCycle: appLifeCycle,
				LocaleFS:     testutils.GetResourceFS(),
			}),
			Routes: []core.HTTPRoute{
				usermgt.NewProfileHandler(usermgt.ProfileHandlerParams{
					Logger:         logger,
					DB:             db,
					UserRepository: userRepository,
				}),
			},
		})

		appLifeCycle.RequireStart()
	})

	AfterEach(func() {
		appLifeCycle.RequireStop()
	})

	It("returns an error if there is no access token", func(ctx SpecContext) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusUnauthorized))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID":  Equal("E-0002"),
			"Message":    Equal("You must be authenticated to use this feature"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})

	It("returns a user profile if there is valid access token", func(ctx SpecContext) {
		userRepository.EXPECT().
			FindUserByID(mock.Anything, mock.Anything, "019135f7-6265-7ef8-8920-57280736f6c0").
			Return(&usermgt.UserModel{
				Model: core.Model{
					ID: uuid.MustParse("019135f7-6265-7ef8-8920-57280736f6c0"),
				},
				HasCreatedAtColumn: core.HasCreatedAtColumn{
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				HasUpdatedAtColumn: core.HasUpdatedAtColumn{
					UpdatedAt: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
				},
				Username:  "testing",
				Email:     "testing@example.com",
				FirstName: "testing",
				LastName:  "testing",
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		recorder := httptest.NewRecorder()
		//nolint:lll // No need to fix
		req.Header.Add(core.AuthorizationHeader, "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMTkxMzVmNy02MjY1LTdlZjgtODkyMC01NzI4MDczNmY2YzAiLCJleHAiOjE3MjQ4NDc3MTksIm5iZiI6MTcyNDg0NzcxOSwiaWF0IjoxNzI0ODQ3NzE5LCJyb2xlcyI6W10sInBlcm1pc3Npb25zIjpbXX0.I9-Kr2ArmW3V-eUN9KKxKShmV9oDWefKBzaXo5BJCqV6fqVtddNFSxnmGzj72WMykCXSTrz92NDGtH8M-lZWwBsNOJY7XCZFoDdYKHk_OyGR9Nk-lRvburgMgaNChw6lD-zjZTb2xfJhmdj4IMbZOcDMB6bdo5bAz_M_3iiPw1gMX9Jkd5yXIwchjOWwVasVO0ycZZ3qFz-mBrSn1FyG8T_ox6avcEHFdiDiBUR6YBaXZwIpiFqhy0aDdvGz8MCvT95b5keTO6jcNLwHZrm1YnZD-lPz5xJQL14n-FnKOvi0UVpEbmkkmyfQz4IH5kdzaRaEdHEYsSyjpNJ1Xaq5lA")
		handler.ServeHTTP(recorder, req)

		var response core.Response[usermgt.ProfileResponse]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusOK))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("S-0000"),
			"Message":   Equal("Success"),
			"Data": MatchFields(IgnoreExtras, Fields{
				"ID":        Equal("019135f7-6265-7ef8-8920-57280736f6c0"),
				"Username":  Equal("testing"),
				"Email":     Equal("testing@example.com"),
				"FirstName": Equal("testing"),
				"LastName":  Equal("testing"),
				"CreatedAt": Equal(time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)),
				"UpdatedAt": Equal(time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)),
			}),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})
})
