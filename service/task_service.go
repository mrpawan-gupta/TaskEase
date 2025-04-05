package service

import (
	"errors"
	"strings"
	"taskease/domain"
	"taskease/repository"
)

var (
	ErrTaskTitleRequired = errors.New("task title is required")
	ErrTaskIDRequired    = errors.New("task ID is required")
	ErrInvalidTaskStatus = errors.New("invalid task status")

	validStatusMap = map[domain.TaskStatus]bool{
		domain.StatusPending:    true,
		domain.StatusInProgress: true,
		domain.StatusCompleted:  true,
	}
)

type TaskService struct {
	repo repository.TaskRepositoryInt
}

func NewTaskService(repo repository.TaskRepositoryInt) *TaskService {
	return &TaskService{repo: repo}
}

func validateStatus(status domain.TaskStatus) bool {
	return validStatusMap[status]
}

func (s *TaskService) CreateTask(task domain.Task) (domain.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		return domain.Task{}, ErrTaskTitleRequired
	}

	if task.Status == "" {
		task.Status = domain.StatusPending
	} else if !validateStatus(task.Status) {
		return domain.Task{}, ErrInvalidTaskStatus
	}

	return s.repo.Create(task)
}

func (s *TaskService) GetTaskByID(id string) (domain.Task, error) {
	if id == "" {
		return domain.Task{}, ErrTaskIDRequired
	}
	return s.repo.GetByID(id)
}

func (s *TaskService) UpdateTask(task domain.Task) (domain.Task, error) {
	if task.ID == "" {
		return domain.Task{}, ErrTaskIDRequired
	}

	existingTask, err := s.repo.GetByID(task.ID)
	if err != nil {
		return domain.Task{}, err
	}

	// Only update fields that are provided
	if title := strings.TrimSpace(task.Title); title != "" {
		existingTask.Title = title
	}

	if task.Description != "" {
		existingTask.Description = task.Description
	}

	if task.Status != "" {
		if !validateStatus(task.Status) {
			return domain.Task{}, ErrInvalidTaskStatus
		}
		existingTask.Status = task.Status
	}

	return s.repo.Update(existingTask)
}

func (s *TaskService) DeleteTask(id string) error {
	if id == "" {
		return ErrTaskIDRequired
	}
	return s.repo.Delete(id)
}

func (s *TaskService) ListTasks(filter domain.TaskFilter) ([]domain.Task, int, error) {
	return s.repo.List(filter)
}
