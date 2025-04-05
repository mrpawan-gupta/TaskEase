package repository

import (
	"database/sql"

	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"taskease/domain"
	"time"
)

type TaskRepositoryInt interface {
	Create(task domain.Task) (domain.Task, error)
	GetByID(id string) (domain.Task, error)
	Update(task domain.Task) (domain.Task, error)
	Delete(id string) error
	List(filter domain.TaskFilter) ([]domain.Task, int, error)
	Close() error
}

type TaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(connectionString string) (*TaskRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &TaskRepository{
		db: db,
	}, nil
}

func (r *TaskRepository) Create(task domain.Task) (domain.Task, error) {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Insert the task
	query := `
		INSERT INTO tasks (id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, status, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) GetByID(id string) (domain.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	var task domain.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Task{}, errors.New("task not found")
	}

	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Update(task domain.Task) (domain.Task, error) {
	// First check if the task exists
	_, err := r.GetByID(task.ID)
	if err != nil {
		return domain.Task{}, err
	}

	// Update timestamp
	task.UpdatedAt = time.Now()

	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, description, status, created_at, updated_at
	`
	err = r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.ID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *TaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *TaskRepository) List(filter domain.TaskFilter) ([]domain.Task, int, error) {
	baseQuery := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	countQuery := `SELECT COUNT(*) FROM tasks`

	whereClause := ""
	args := []interface{}{}
	argPos := 1

	if filter.Status != nil {
		whereClause = fmt.Sprintf(" WHERE status = $%d", argPos)
		args = append(args, string(*filter.Status))
		argPos++
	}

	// Add pagination
	limit := filter.Limit
	if limit <= 0 {
		limit = 10 // Default limit
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	paginationClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	// Get total count
	var totalCount int
	err := r.db.QueryRow(countQuery+whereClause, args[:argPos-1]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get tasks
	rows, err := r.db.Query(baseQuery+whereClause+" ORDER BY created_at DESC"+paginationClause, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return tasks, totalCount, nil
}

func (r *TaskRepository) Close() error {
	return r.db.Close()
}
