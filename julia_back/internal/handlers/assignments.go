package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssignmentHandler struct {
	assignmentService services.AssignmentService
}

func NewAssignmentHandler(assignmentService services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{assignmentService: assignmentService}
}

func (h *AssignmentHandler) GetAllAssignments(c *gin.Context) {
	assignments, err := h.assignmentService.GetAllAssignments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func (h *AssignmentHandler) GetAssignmentsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	assignments, err := h.assignmentService.GetAssignmentsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func (h *AssignmentHandler) GetAssignmentsByMakeupID(c *gin.Context) {
	makeupID := c.Param("makeupID")
	makeupIDInt, err := strconv.ParseInt(makeupID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	assignments, err := h.assignmentService.GetAssignmentsByMakeupID(makeupIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func (h *AssignmentHandler) GetAssignmentWithClassID(c *gin.Context) {
	assignments, err := h.assignmentService.GetAssignmentWithClassID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignments)
}

func (h *AssignmentHandler) UpsertAssignment(c *gin.Context) {
	var assignment models.Assignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.assignmentService.UpsertAssignment(&assignment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assignment)
}

func (h *AssignmentHandler) DeleteAssignment(c *gin.Context) {
	assignmentID := c.Param("assignmentID")
	assignmentIDInt, err := strconv.ParseInt(assignmentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.assignmentService.DeleteAssignment(assignmentIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
}
