package usermgt_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"
	"wano-island/common/core"
	"wano-island/common/usermgt"
	"wano-island/console/modules/httpsrv"
	mockcore "wano-island/testing/mocks/common/core"
	"wano-island/testing/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

var _ = Describe("Login Handler", func() {
	var (
		db             *gorm.DB
		mockedDB       sqlmock.Sqlmock
		config         *mockcore.MockAppConfig
		userRepository usermgt.UserRepository
		router         http.Handler
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		db, mockedDB = testutils.CreateTestDBInstance()

		config = mockcore.NewMockAppConfig(GinkgoT())
		config.EXPECT().GetAppVersion().Return("1.0.0")
		config.EXPECT().GetRevision().Return("testing")
		config.EXPECT().GetMode().Return(core.TestingMode)
		config.EXPECT().IsTesting().Return(true)
		config.EXPECT().GetJWTConfig().Return(generateJWTConfig())

		userRepository = usermgt.NewUserRepository(usermgt.UserRepositoryParams{})
		userService := usermgt.NewUserService(usermgt.UserServiceParams{})

		router = testutils.WithFxLifeCycle(func(l *fxtest.Lifecycle) http.Handler {
			return testutils.CreateRouter(func(rp *httpsrv.RouteParams) {
				rp.Config = config
				rp.Routes = []core.HTTPRoute{
					usermgt.NewLoginHandler(usermgt.LoginHandlerParams{
						AppLifeCycle:   l,
						Logger:         core.NewNoopLogger(),
						Config:         config,
						DB:             db,
						UserService:    userService,
						UserRepository: userRepository,
					}),
				}
			})
		})
	})

	It("returns an error if cannot decode request body", func(ctx SpecContext) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{,}")))

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusBadRequest))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID":  Equal("E-0001"),
			"Message":    Equal("Cannot decode the request body, please re-check the request body"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})

	It("returns an error if email is incorrect", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": ""
        }`)))

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusUnauthorized))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID":  Equal("E-0000"),
			"Message":    Equal("Invalid email or password"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})

	It("returns an error if password is incorrect", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"email", "password"}).AddRow(
				"testing@internal.com",
				"$2a$10$4LGRfD5yIX02UIe.4mEmfO60OkPVOQ5rsWgVS708v2TkurwnRW51."))

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": "incorrect-password"
        }`)))

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusUnauthorized))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID":  Equal("E-0000"),
			"Message":    Equal("Invalid email or password"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})

	It("returns an accessToken and refreshToken if the user inputs the correct credentials.", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"email", "password"}).AddRow(
				"testing@internal.com",
				"$2a$10$4LGRfD5yIX02UIe.4mEmfO60OkPVOQ5rsWgVS708v2TkurwnRW51."))

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": "Keep!t5ecret"
        }`)))

		router.ServeHTTP(recorder, req)

		var response core.Response[usermgt.LoginResponse]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusOK))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("S-0000"),
			"Message":   Equal("Success"),
			"Data": MatchFields(IgnoreMissing, Fields{
				"AccessToken":  Not(BeEmpty()),
				"RefreshToken": Not(BeEmpty()),
			}),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
		}))
	})
})
