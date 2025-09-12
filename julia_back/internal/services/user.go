package services

import (
	"julia/internal/models"
	"julia/internal/repositories"
	"julia/utils"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	GetTodosByUserID(userID string) ([]*models.ClassTodos, error)
	CreateUser(user *models.User) error
	UpdateUser(id string, user *models.User) error
	DeleteUser(id string) error
}

type userService struct {
	userRepo repositories.UserRepository
	todoRepo repositories.TodoRepository
}

func NewUserService(userRepo repositories.UserRepository, todoRepo repositories.TodoRepository) UserService {
	return &userService{userRepo: userRepo, todoRepo: todoRepo}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}

func (s *userService) GetTodosByUserID(userID string) ([]*models.ClassTodos, error) {
	rows, err := s.todoRepo.GetTodosByUserID(userID)
	if err != nil {
		return nil, err
	}
	return groupByClass(rows), nil
}

func groupByClass(rows []*models.TodoRow) []*models.ClassTodos {
	if len(rows) == 0 {
		return []*models.ClassTodos{}
	}

	var result []*models.ClassTodos
	var curClassTodos *models.ClassTodos
	var curClassID int64 = -2

	for _, row := range rows {
		if row.ClassID != curClassID {
			if curClassTodos != nil {
				result = append(result, curClassTodos)
			}
			curClassTodos = &models.ClassTodos{
				ClassID:   row.ClassID,
				ClassName: row.ClassName,
				Todos:     make([]*models.Todo, 0),
			}
			curClassID = row.ClassID
		}
		curClassTodos.Todos = append(curClassTodos.Todos, &models.Todo{
			ID:          row.TodoID,
			Title:       row.Title,
			UserID:      row.UserID,
			ClassID:     &row.ClassID,
			Description: row.Description,
			Completed:   row.Completed,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		})
	}

	if curClassTodos != nil {
		result = append(result, curClassTodos)
	}
	return result
}

func (s *userService) CreateUser(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.CreateUser(user)
}

func (s *userService) UpdateUser(id string, user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.userRepo.UpdateUser(id, user)
}

func (s *userService) DeleteUser(id string) error {
	return s.userRepo.DeleteUser(id)
}
