package showmgt_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"
	"wano-island/common/core"
	"wano-island/common/showmgt"
	"wano-island/console/modules/httpsrv"
	mockcore "wano-island/testing/mocks/common/core"
	"wano-island/testing/testutils"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"gorm.io/gorm"
)

var _ = Describe("[handler.create-movie.go]", func() {
	var (
		db       *gorm.DB
		mockedDB sqlmock.Sqlmock
		router   http.Handler
		config   *mockcore.MockAppConfig
	)

	BeforeEach(func() {
		testutils.DetectLeakyGoroutines()
		db, mockedDB = testutils.CreateTestDBInstance()

		config = mockcore.NewMockAppConfig(GinkgoT())
		config.EXPECT().GetAppVersion().Return("1.0.0")
		config.EXPECT().GetRevision().Return("testing")
		config.EXPECT().GetMode().Return(core.TestingMode)
		config.EXPECT().IsTesting().Return(true)
		config.EXPECT().GetJWTConfig().Return(testutils.GetJWTConfig())
		config.EXPECT().GetCorsConfig().Return(&core.CorsConfig{})

		router = testutils.CreateRouter(func(rp *httpsrv.RouteParams) {
			rp.Config = config
			rp.Routes = []core.HTTPRoute{
				showmgt.NewCreateMovieHandler(showmgt.CreateShowHandlerParams{
					Logger: core.NewNoopLogger(),
					DB:     db,
				}),
			}
		})
	})

	It("should return an error if cannot decode the request body", func() {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/shows", bytes.NewReader([]byte(`
        {
            invalid-data
        }`)))
		router.ServeHTTP(recorder, testutils.WithFakeJWT(request))

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

	It("should return an error if cannot create a show", func() {
		mockedDB.ExpectBegin()
		mockedDB.ExpectExec(`INSERT INTO "public"."shows"`).
			WithArgs(
				// id
				testutils.AnyUUIDArg{},
				// created_at
				testutils.AnyTimeArg{},
				// updated_At
				testutils.AnyTimeArg{},
				// kind
				"movie",
				// original_language
				"ja",
				// original_title
				"Naruto - Title",
				// original_overview
				"Naruto - Overview",
				// keywords
				`{"naruto"}`,
				// is_released
				true,
			).
			WillReturnError(errors.New("something went wrong"))
		mockedDB.ExpectRollback()

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/shows", bytes.NewReader([]byte(`
        {
            "kind": "movie",
            "originalLanguage": "ja",
            "originalTitle": "Naruto - Title",
            "originalOverview": "Naruto - Overview",
            "keywords": [
                "naruto"
            ],
            "isReleased": true
        }`)))

		router.ServeHTTP(recorder, testutils.WithFakeJWT(request))

		var response core.Response[any]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusBadRequest))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("E-0007"),
			//nolint:lll // No need to fix
			"Message":    Equal("Cannot create the show. Please try again. If the problem continues, kindly reach out to the system administrator for support"),
			"Data":       BeNil(),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})

	It("should create a show successfully", func() {
		mockedDB.ExpectBegin()
		mockedDB.ExpectExec(`INSERT INTO "public"."shows"`).
			WithArgs(
				// id
				testutils.AnyUUIDArg{},
				// created_at
				testutils.AnyTimeArg{},
				// updated_At
				testutils.AnyTimeArg{},
				// kind
				"movie",
				// original_language
				"ja",
				// original_title
				"Naruto - Title",
				// original_overview
				"Naruto - Overview",
				// keywords
				`{"naruto"}`,
				// is_released
				true,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockedDB.ExpectCommit()

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/shows", bytes.NewReader([]byte(`
        {
            "kind": "movie",
            "originalLanguage": "ja",
            "originalTitle": "Naruto - Title",
            "originalOverview": "Naruto - Overview",
            "keywords": [
                "naruto"
            ],
            "isReleased": true
        }`)))

		router.ServeHTTP(recorder, testutils.WithFakeJWT(request))

		var response core.Response[showmgt.ShowDTO]
		_ = json.Unmarshal(recorder.Body.Bytes(), &response)

		Expect(recorder).To(HaveHTTPStatus(http.StatusCreated))
		Expect(response).To(MatchFields(IgnoreMissing, Fields{
			"MessageID": Equal("S-0000"),
			"Message":   Equal("Success"),
			"Data": MatchFields(IgnoreMissing, Fields{
				"ID":               Not(BeEmpty()),
				"Kind":             Equal("movie"),
				"OriginalLanguage": Equal("ja"),
				"OriginalTitle":    Equal("Naruto - Title"),
				"OriginalOverview": PointTo(Equal("Naruto - Overview")),
				"Keywords":         Equal([]string{"naruto"}),
				"IsReleased":       Equal(true),
				"CreatedAt":        BeTemporally("~", time.Now(), time.Minute),
				"UpdatedAt":        BeTemporally("~", time.Now(), time.Minute),
			}),
			"Pagination": BeNil(),
			"Timestamp":  BeTemporally("~", time.Now(), time.Minute),
			"RequestID":  Not(BeEmpty()),
		}))
	})
})
