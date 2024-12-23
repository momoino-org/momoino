package core

import (
	"math"
	"net/http"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalRows  int64 `json:"totalRows"`
	TotalPages int   `json:"totalPages"`
}

type Response[T any] struct {
	MessageID  string      `json:"messageId"`
	Message    string      `json:"message"`
	Data       T           `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"requestId"`
}

type ResponseBuilder struct {
	request  *http.Request
	response *Response[any]
}

func (r *ResponseBuilder) MessageID(id string) *ResponseBuilder {
	r.response.MessageID = id
	return r
}

func (r *ResponseBuilder) Data(data any) *ResponseBuilder {
	r.response.Data = data

	return r
}

func (r *ResponseBuilder) Pagination(totalRows int64) *ResponseBuilder {
	page := GetPage(r.request)
	pageSize := GetPageSize(r.request)

	r.response.Pagination = &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(pageSize))),
	}

	return r
}

func (r *ResponseBuilder) Build() *Response[any] {
	localizer := GetLocalizer(r.request)

	if r.response.MessageID == "" {
		r.response.MessageID = MsgSuccess
	}

	r.response.Message = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: r.response.MessageID,
		},
	})

	r.response.RequestID = GetRequestID(r.request)
	r.response.Timestamp = time.Now()

	return r.response
}

func NewResponseBuilder(request *http.Request) *ResponseBuilder {
	return &ResponseBuilder{
		request:  request,
		response: &Response[any]{},
	}
}
