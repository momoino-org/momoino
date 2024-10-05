package core

import (
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
	"go.uber.org/fx"
)

type requestIDCtxKey string

const RequestIDHeader = "X-Request-Id"
const RequestIDKey requestIDCtxKey = "RequestIDKey"
const AuthorizationHeader = "X-Auth-Access-Token"
const DefaultPage = 1
const DefaultPageSize = 10
const MaxPageSize = 100

func GetRequestID(r *http.Request) string {
	return r.Header.Get(RequestIDHeader)
}

func NewRequestModule() fx.Option {
	return fx.Module(
		"Request Module",
		fx.Provide(func() *schema.Decoder {
			schema := schema.NewDecoder()
			schema.IgnoreUnknownKeys(true)

			return schema
		}),
	)
}

// GetPage retrieves the page number from the URL query parameters of the given HTTP request.
// If the "page" parameter is not provided or cannot be converted to an integer,
// it returns the default page number (DefaultPage) which is 1.
// If the retrieved page number is less than the default page number,
// it also returns the default page number.
func GetPage(r *http.Request) int {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		return DefaultPage
	}

	if page < DefaultPage {
		return DefaultPage
	}

	return page
}

// GetPageSize retrieves the page size from the URL query parameters of the given HTTP request.
// If the "pageSize" parameter is not provided or cannot be converted to an integer,
// it returns the default page size which is 10.
// If the retrieved page size is less than 1, it returns the default page size.
// If the retrieved page size is greater than the maximum page size (100),
// it returns the maximum page size.
func GetPageSize(r *http.Request) int {
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if err != nil {
		return DefaultPageSize
	}

	if pageSize < 1 {
		return DefaultPageSize
	}

	if pageSize > MaxPageSize {
		return MaxPageSize
	}

	return pageSize
}

// GetOffset calculates the offset for pagination based on the provided HTTP request.
func GetOffset(r *http.Request) int {
	page := GetPage(r)
	pageSize := GetPageSize(r)

	return (page - 1) * pageSize
}
