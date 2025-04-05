package api

import (
	"net/http"
	"time"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     []string    `json:"errors,omitempty"`
	StatusCode int         `json:"status_code"`
	Timestamp  time.Time   `json:"timestamp"`
}

func NewResponse(success bool, message string, data interface{}, statusCode int) *Response {
	return &Response{
		Success:    success,
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

func SuccessResponse(message string, data interface{}) *Response {
	return NewResponse(true, message, data, http.StatusOK)
}

func ErrorResponse(message string, errors []string, statusCode int) *Response {
	return &Response{
		Success:    false,
		Message:    message,
		Errors:     errors,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

type PaginatedResponse struct {
	Response
	Pagination *Pagination `json:" ,omitempty"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	PerPage     int `json:"per_page"`
	TotalItems  int `json:"total_items"`
}

func NewPaginatedResponse(success bool, message string, data interface{}, statusCode int, pagination *Pagination) *PaginatedResponse {
	return &PaginatedResponse{
		Response: Response{
			Success:    success,
			Message:    message,
			Data:       data,
			StatusCode: statusCode,
			Timestamp:  time.Now(),
		},
		Pagination: pagination,
	}
}
