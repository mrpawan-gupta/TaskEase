package api

import (
	"encoding/json"
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

func SuccessResponse(success bool, message string, data interface{}, statusCode int) *Response {
	return &Response{
		Success:    success,
		Message:    message,
		Data:       data,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
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
	Pagination *Pagination `json:"pagination,omitempty"`
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

func WriteResponse(w http.ResponseWriter, r *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return
	}
}

func WritePaginatedResponse(w http.ResponseWriter, r *PaginatedResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		return
	}
}

func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	response := SuccessResponse(true, message, data, statusCode)
	WriteResponse(w, response)
}

func WriteErrorResponse(w http.ResponseWriter, message string, errors []string, statusCode int) {
	response := ErrorResponse(message, errors, statusCode)
	WriteResponse(w, response)
}
