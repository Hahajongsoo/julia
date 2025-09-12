package repositories

import (
	"database/sql"
	"julia/internal/models"
)

type TodoRepository interface {
	GetAllTodos() ([]*models.Todo, error)
	GetTodosByUserID(userID string) ([]*models.TodoRow, error)
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
		ORDER BY created_at DESC
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

func (r *todoRepository) GetTodosByUserID(userID string) ([]*models.TodoRow, error) {
	query := `
		SELECT t.id, t.title, t.user_id, COALESCE(t.class_id, -1) as class_id, COALESCE(c.class_name, '할 일') as class_name, t.description, t.completed, t.created_at, t.updated_at
		FROM todos t
		LEFT JOIN classes c ON t.class_id = c.class_id
		WHERE t.user_id = $1
		ORDER BY class_id, t.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := make([]*models.TodoRow, 0)
	for rows.Next() {
		var todo models.TodoRow
		err := rows.Scan(&todo.TodoID, &todo.Title, &todo.UserID, &todo.ClassID, &todo.ClassName, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

func (r *todoRepository) GetTodoByID(id int64) (*models.Todo, error) {
	query := `
		SELECT id, title, user_id, class_id, description, completed, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)
	var todo models.Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.UserID, &todo.ClassID, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) CreateTodo(todo *models.Todo) error {
	query := `
		INSERT INTO todos (title, user_id, class_id, description, completed)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, todo.Title, todo.UserID, todo.ClassID, todo.Description, todo.Completed).
		Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
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
		RETURNING updated_at
	`
	err := r.db.QueryRow(query, todo.Title, todo.UserID, todo.Description, todo.Completed, id).
		Scan(&todo.UpdatedAt)
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
