package core

import (
	"math"
	"time"
)

type Pagination struct {
	Page       int   `json:"currentPage"`
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
}

type ResponseBuilder struct {
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

func (r *ResponseBuilder) Pagination(page int, limit int, totalRows int64) *ResponseBuilder {
	r.response.Pagination = &Pagination{
		Page:       page,
		PageSize:   limit,
		TotalRows:  totalRows,
		TotalPages: int(math.Ceil(float64(totalRows) / float64(limit))),
	}

	return r
}

func (r *ResponseBuilder) Build() *Response[any] {
	if r.response.MessageID == "" {
		r.response.MessageID = "S0200"
		r.response.Message = "Success"
	}

	r.response.Timestamp = time.Now()

	return r.response
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{
		response: &Response[any]{},
	}
}
