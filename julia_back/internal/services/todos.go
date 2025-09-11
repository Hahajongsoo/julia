package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
)

type TodoService interface {
	GetAllTodos() ([]*models.Todo, error)
	GetTodoByID(id int64) (*models.Todo, error)
	CreateTodo(todo *models.Todo) error
	UpdateTodo(id int64, todo *models.Todo) error
	DeleteTodo(id int64) error
}

type todoService struct {
	todoRepo repositories.TodoRepository
}

func NewTodoService(todoRepo repositories.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (s *todoService) GetAllTodos() ([]*models.Todo, error) {
	return s.todoRepo.GetAllTodos()
}

func (s *todoService) GetTodoByID(id int64) (*models.Todo, error) {
	return s.todoRepo.GetTodoByID(id)
}

func (s *todoService) CreateTodo(todo *models.Todo) error {
	return s.todoRepo.CreateTodo(todo)
}

func (s *todoService) UpdateTodo(id int64, todo *models.Todo) error {
	return s.todoRepo.UpdateTodo(id, todo)
}

func (s *todoService) DeleteTodo(id int64) error {
	return s.todoRepo.DeleteTodo(id)
}
