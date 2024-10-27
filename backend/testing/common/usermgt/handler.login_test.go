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
	"wano-island/console/modules/httpsrv"
	"wano-island/testing/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gorm.io/gorm"
)

var _ = Describe("Login Handler", Ordered, func() {
	var (
		db                  *gorm.DB
		mockedDB            sqlmock.Sqlmock
		userRepository      usermgt.UserRepository
		router              http.Handler
		sessionManager      *scs.SessionManager
		loginSessionManager *scs.SessionManager
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		testutils.ConfigureMinimumEnvVariables()
		db, mockedDB = testutils.CreateTestDBInstance()

		config, err := core.NewAppConfig()
		Expect(err).NotTo(HaveOccurred())

		userRepository = usermgt.NewUserRepository(usermgt.UserRepositoryParams{})
		userService := usermgt.NewUserService(usermgt.UserServiceParams{
			Logger: core.NewNoopLogger(),
			Config: config,
		})

		DeferCleanup(func() {
			sessionManager, err = core.NewSessionManager(config)
			Expect(err).NotTo(HaveOccurred())

			loginSessionManager, err = core.NewLoginSessionManager(config)
			Expect(err).NotTo(HaveOccurred())

			if gstore, ok := sessionManager.Store.(*memstore.MemStore); ok {
				gstore.StopCleanup()
			}

			if memStore, ok := loginSessionManager.Store.(*memstore.MemStore); ok {
				memStore.StopCleanup()
			}
		})

		router = testutils.CreateRouter(func(rp *httpsrv.RouteParams) {
			rp.Config = config
			rp.Routes = []core.HTTPRoute{
				usermgt.NewLoginHandler(usermgt.LoginHandlerParams{
					Logger:              core.NewNoopLogger(),
					Config:              config,
					DB:                  db,
					SessionManager:      sessionManager,
					LoginSessionManager: loginSessionManager,
					UserService:         userService,
					UserRepository:      userRepository,
				}),
			}
		})
	})

	It("returns an error if missing the CSRF token & CSRF cookie", func() {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{,}")))

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusForbidden))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("E_INVALID_CSRF"),
			//nolint:lll // This is for testing purpose
			"Message":    Equal("Your session has expired or there was an issue with your request. Please refresh the page and try again."),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns an error if login session cookie is invalid", func() {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", nil)
		req.Header.Set("X-Csrf-Token", "CxffrOLOHQ34nokOSwCzNTcZQ+9AFci4YzmHuaio+G0h0+Oz3XGJRukXwoZT2BjsIyzPF4VlMlYmemY/D5yjjw==")
		req.Header.Set("Cookie", fmt.Sprintf("%v=%v; %v=%v;",
			core.CsrfCookie,
			"MTczMDAyMTExOXxJa3R6VVRoSWVpc3ZiRVZ6VW1sVmRVbEhUbWx5TWxKUk1XcFFha1pqVUhKMVVsVlFhR2h4WXpCWEswazlJZ289fGTZUOz5rj3AJUXgKJrDwyuAUdVq-Bq3b1L_OdxHLPC-",
			core.LoginSessionCookie,
			"NGGdyYOThhTRopdAz2ZWCKNVJKolyDWIgxXVYfYO8cY",
		))

		router.ServeHTTP(recorder, req)

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusForbidden))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("E_INVALID_SESSION"),

			"Message":    Equal("Uh-oh! It looks like your session is no longer valid."),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns an error if cannot decode request body", func(ctx SpecContext) {
		recorder := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte("{,}")))
		sessionCtx, err := loginSessionManager.Load(req.Context(), "")
		Expect(err).NotTo(HaveOccurred())
		sessionToken, _, err := loginSessionManager.Commit(sessionCtx)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("X-Csrf-Token", "CxffrOLOHQ34nokOSwCzNTcZQ+9AFci4YzmHuaio+G0h0+Oz3XGJRukXwoZT2BjsIyzPF4VlMlYmemY/D5yjjw==")
		req.Header.Set("Cookie", fmt.Sprintf("%v=%v; %v=%v;",
			core.CsrfCookie,
			"MTczMDAyMTExOXxJa3R6VVRoSWVpc3ZiRVZ6VW1sVmRVbEhUbWx5TWxKUk1XcFFha1pqVUhKMVVsVlFhR2h4WXpCWEswazlJZ289fGTZUOz5rj3AJUXgKJrDwyuAUdVq-Bq3b1L_OdxHLPC-",
			core.LoginSessionCookie,
			sessionToken,
		))

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
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns an error if username is incorrect", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "username": "testing@internal.com",
            "password": ""
        }`)))
		sessionCtx, err := loginSessionManager.Load(req.Context(), "")
		Expect(err).NotTo(HaveOccurred())
		sessionToken, _, err := loginSessionManager.Commit(sessionCtx)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("X-Csrf-Token", "CxffrOLOHQ34nokOSwCzNTcZQ+9AFci4YzmHuaio+G0h0+Oz3XGJRukXwoZT2BjsIyzPF4VlMlYmemY/D5yjjw==")
		req.Header.Set("Cookie", fmt.Sprintf("%v=%v; %v=%v;",
			core.CsrfCookie,
			"MTczMDAyMTExOXxJa3R6VVRoSWVpc3ZiRVZ6VW1sVmRVbEhUbWx5TWxKUk1XcFFha1pqVUhKMVVsVlFhR2h4WXpCWEswazlJZ289fGTZUOz5rj3AJUXgKJrDwyuAUdVq-Bq3b1L_OdxHLPC-",
			core.LoginSessionCookie,
			sessionToken,
		))

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
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns an error if password is incorrect", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).AddRow(
				"testing@internal.com",
				"$2a$10$4LGRfD5yIX02UIe.4mEmfO60OkPVOQ5rsWgVS708v2TkurwnRW51."))

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "username": "testing@internal.com",
            "password": "incorrect-password"
        }`)))
		sessionCtx, err := loginSessionManager.Load(req.Context(), "")
		Expect(err).NotTo(HaveOccurred())
		sessionToken, _, err := loginSessionManager.Commit(sessionCtx)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("X-Csrf-Token", "CxffrOLOHQ34nokOSwCzNTcZQ+9AFci4YzmHuaio+G0h0+Oz3XGJRukXwoZT2BjsIyzPF4VlMlYmemY/D5yjjw==")
		req.Header.Set("Cookie", fmt.Sprintf("%v=%v; %v=%v;",
			core.CsrfCookie,
			"MTczMDAyMTExOXxJa3R6VVRoSWVpc3ZiRVZ6VW1sVmRVbEhUbWx5TWxKUk1XcFFha1pqVUhKMVVsVlFhR2h4WXpCWEswazlJZ289fGTZUOz5rj3AJUXgKJrDwyuAUdVq-Bq3b1L_OdxHLPC-",
			core.LoginSessionCookie,
			sessionToken,
		))

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
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("returns an accessToken if the user inputs the correct credentials.", func(ctx SpecContext) {
		mockedDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WithArgs("testing@internal.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "password"}).AddRow(
				"testing@internal.com",
				"$2a$10$4LGRfD5yIX02UIe.4mEmfO60OkPVOQ5rsWgVS708v2TkurwnRW51."))

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader([]byte(`{
            "username": "testing@internal.com",
            "password": "Keep!t5ecret"
        }`)))
		sessionCtx, err := loginSessionManager.Load(req.Context(), "")
		Expect(err).NotTo(HaveOccurred())
		sessionToken, _, err := loginSessionManager.Commit(sessionCtx)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set("X-Csrf-Token", "CxffrOLOHQ34nokOSwCzNTcZQ+9AFci4YzmHuaio+G0h0+Oz3XGJRukXwoZT2BjsIyzPF4VlMlYmemY/D5yjjw==")
		req.Header.Set("Cookie", fmt.Sprintf("%v=%v; %v=%v;",
			core.CsrfCookie,
			"MTczMDAyMTExOXxJa3R6VVRoSWVpc3ZiRVZ6VW1sVmRVbEhUbWx5TWxKUk1XcFFha1pqVUhKMVVsVlFhR2h4WXpCWEswazlJZ289fGTZUOz5rj3AJUXgKJrDwyuAUdVq-Bq3b1L_OdxHLPC-",
			core.LoginSessionCookie,
			sessionToken,
		))

		router.ServeHTTP(recorder, req)

		var response core.Response[usermgt.LoginResponse]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusOK))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("S-0000"),
			"Message":   Equal("Success"),
			"Data": MatchFields(IgnoreMissing, Fields{
				"AccessToken": Not(BeEmpty()),
				// "RefreshToken": Not(BeEmpty()),
			}),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})
})
