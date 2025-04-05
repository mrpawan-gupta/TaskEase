package api

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"taskease/domain"
	"taskease/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteErrorResponse(
			w,
			"Invalid request payload",
			[]string{"The request body could not be parsed as JSON"},
			http.StatusBadRequest,
		)
		return
	}

	createdTask, err := h.service.CreateTask(task)
	if err != nil {
		WriteErrorResponse(
			w,
			"Failed to create task",
			[]string{err.Error()},
			http.StatusBadRequest,
		)
		return
	}

	WriteSuccessResponse(
		w,
		"Task created successfully",
		createdTask,
		http.StatusCreated,
	)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		WriteErrorResponse(
			w,
			"Task not found",
			[]string{err.Error()},
			http.StatusNotFound,
		)
		return
	}

	WriteSuccessResponse(
		w,
		"Task retrieved successfully",
		task,
		http.StatusOK,
	)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteErrorResponse(
			w,
			"Invalid request payload",
			[]string{"The request body could not be parsed as JSON"},
			http.StatusBadRequest,
		)
		return
	}

	// Set ID from path parameter
	task.ID = id

	updatedTask, err := h.service.UpdateTask(task)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, service.ErrTaskIDRequired) || err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		}

		WriteErrorResponse(
			w,
			"Failed to update task",
			[]string{err.Error()},
			statusCode,
		)
		return
	}

	WriteSuccessResponse(
		w,
		"Task updated successfully",
		updatedTask,
		http.StatusOK,
	)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteTask(id); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, service.ErrTaskIDRequired) || err.Error() == "task not found" {
			statusCode = http.StatusNotFound
		}

		WriteErrorResponse(
			w,
			"Failed to delete task",
			[]string{err.Error()},
			statusCode,
		)
		return
	}

	WriteSuccessResponse(
		w,
		"Task deleted successfully",
		nil,
		http.StatusOK,
	)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	statusStr := r.URL.Query().Get("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var statusFilter *domain.TaskStatus
	if statusStr != "" {
		s := domain.TaskStatus(statusStr)
		statusFilter = &s
	}

	filter := domain.TaskFilter{
		Status: statusFilter,
		Limit:  limit,
		Offset: offset,
	}

	tasks, totalCount, err := h.service.ListTasks(filter)
	if err != nil {
		WriteErrorResponse(
			w,
			"Failed to retrieve tasks",
			[]string{err.Error()},
			http.StatusInternalServerError,
		)
		return
	}

	totalPages := (totalCount + limit - 1) / limit
	pagination := &Pagination{
		CurrentPage: (offset / limit) + 1,
		TotalPages:  totalPages,
		PerPage:     limit,
		TotalItems:  totalCount,
	}

	paginatedResponse := NewPaginatedResponse(
		true,
		"Tasks retrieved successfully",
		tasks,
		http.StatusOK,
		pagination,
	)

	WritePaginatedResponse(w, paginatedResponse)
}
