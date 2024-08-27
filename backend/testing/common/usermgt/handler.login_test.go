package usermgt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"
	"wano-island/common/core"
	"wano-island/common/usermgt"
	"wano-island/testing/internal"
	mockcore "wano-island/testing/mocks/common/core"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
	. "github.com/onsi/gomega/gstruct"
	"go.uber.org/fx/fxtest"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var _ = Describe("Login Handler", func() {
	var (
		db             *gorm.DB
		mockedDB       sqlmock.Sqlmock
		appLifeCycle   *fxtest.Lifecycle
		config         *mockcore.MockAppConfig
		userRepository usermgt.UserRepository
		localizer      *i18n.Localizer
		handler        core.HTTPRoute
	)

	BeforeEach(func() {
		db, mockedDB = internal.CreateTestDBInstance()
		appLifeCycle = fxtest.NewLifecycle(GinkgoT())
		config = mockcore.NewMockAppConfig(GinkgoT())
		config.EXPECT().GetJWTConfig().Return(generateJWTConfig())
		userRepository = usermgt.NewUserRepository(usermgt.UserRepositoryParams{})
		userService := usermgt.NewUserService(usermgt.UserServiceParams{})
		localizer = i18n.NewLocalizer(core.NewI18nBundle(core.I18nBundleParams{
			AppLifeCycle: appLifeCycle,
			LocaleFS:     internal.GetResourceFS(),
		}), language.English.String())
		handler = usermgt.NewLoginHandler(usermgt.LoginHandlerParams{
			AppLifeCycle:   appLifeCycle,
			Logger:         core.NewNoopLogger(),
			Config:         config,
			DB:             db,
			UserService:    userService,
			UserRepository: userRepository,
		})

		appLifeCycle.RequireStart()
	})

	AfterEach(func() {
		appLifeCycle.RequireStop()
		mockedDB.ExpectClose()
		internal.CloseGormDB(db)
		Expect(mockedDB.ExpectationsWereMet()).NotTo(HaveOccurred())
		Eventually(Goroutines).ShouldNot(HaveLeaked())
	})

	It("ensures API login is not changed", func() {
		Expect(handler.Pattern()).To(Equal(fmt.Sprintf("%s %s", http.MethodPost, "/api/v1/login")))
	})

	It("returns an error if cannot decode request body", func(ctx SpecContext) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{,}")))
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, core.WithLocalizer(req, localizer))

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

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": ""
        }`)))
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, core.WithLocalizer(req, localizer))

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

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": "incorrect-password"
        }`)))
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, core.WithLocalizer(req, localizer))

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

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "email": "testing@internal.com",
            "password": "Keep!t5ecret"
        }`)))
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, core.WithLocalizer(req, localizer))

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
