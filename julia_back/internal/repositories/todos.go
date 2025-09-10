package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type TodoRepository interface {
	GetAllTodos() ([]*models.Todo, error)
	GetTodosByUserID(userID string) ([]*models.Todo, error)
	GetTodoByID(id int64) (*models.Todo, error)
	CreateTodo(todo *models.Todo) error
	UpdateTodo(id int64, todo *models.Todo) error
	DeleteTodo(id int64) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) GetAllTodos() ([]*models.Todo, error) {
	query := `
		SELECT id, title, user_id, description, completed, created_at, updated_at
		FROM todos
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := make([]*models.Todo, 0)
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.UserID, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *todoRepository) GetTodosByUserID(userID string) ([]*models.Todo, error) {
	query := `
		SELECT id, title, user_id, description, completed, created_at, updated_at
		FROM todos
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := make([]*models.Todo, 0)
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.UserID, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *todoRepository) GetTodoByID(id int64) (*models.Todo, error) {
	query := `
		SELECT id, title, user_id, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.UserID, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) CreateTodo(todo *models.Todo) error {
	query := `
		INSERT INTO todos (title, user_id, description, completed)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, todo.Title, todo.UserID, todo.Description, todo.Completed)
	if err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) UpdateTodo(id int64, todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, user_id = $2, description = $3, completed = $4, updated_at = now()
		WHERE id = $5
	`
	_, err := r.db.Exec(query, todo.Title, todo.UserID, todo.Description, todo.Completed, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) DeleteTodo(id int64) error {
	query := `
		DELETE FROM todos
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
