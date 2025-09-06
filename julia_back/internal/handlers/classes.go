package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClassHandler struct {
	ClassService services.ClassService
}

func NewClassHandler(classService services.ClassService) *ClassHandler {
	return &ClassHandler{ClassService: classService}
}

func (h *ClassHandler) GetAllClasses(c *gin.Context) {
	classes, err := h.ClassService.GetAllClasses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *ClassHandler) GetClassByID(c *gin.Context) {
	classID := c.Param("classID")
	classIDInt, err := strconv.ParseInt(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	class, err := h.ClassService.GetClassByID(classIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *ClassHandler) CreateClass(c *gin.Context) {
	var class *models.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.ClassService.CreateClass(class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *ClassHandler) UpdateClass(c *gin.Context) {
	var class *models.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.ClassService.UpdateClass(class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

func (h *ClassHandler) DeleteClass(c *gin.Context) {
	classID := c.Param("classID")
	classIDInt, err := strconv.ParseInt(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.ClassService.DeleteClass(classIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Class deleted successfully"})
}
