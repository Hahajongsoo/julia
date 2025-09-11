package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responseUsers := make([]*models.ResponseUser, len(users))
	for i, user := range users {
		responseUsers[i] = user.ToResponseUser()
	}
	c.JSON(http.StatusOK, responseUsers)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.ToResponseUser())
}

func (h *UserHandler) GetTodosByUserID(c *gin.Context) {
	id := c.Param("id")
	todos, err := h.userService.GetTodosByUserID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *UserHandler) Create(c *gin.Context) {
	var inputUser *models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.CreateUser(inputUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inputUser.ToResponseUser())
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var inputUser *models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userService.UpdateUser(id, inputUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inputUser.ToResponseUser())
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.userService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
