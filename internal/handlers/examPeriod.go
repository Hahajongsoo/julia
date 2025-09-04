package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExamPeriodHandler struct {
	examPeriodService services.ExamPeriodService
}

func NewExamPeriodHandler(examPeriodService services.ExamPeriodService) *ExamPeriodHandler {
	return &ExamPeriodHandler{examPeriodService: examPeriodService}
}

func (h *ExamPeriodHandler) GetExamPeriodByClassID(c *gin.Context) {
	classID := c.Param("classID")
	classIDInt, err := strconv.ParseInt(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	examPeriods, err := h.examPeriodService.GetExamPeriodByClassID(classIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, examPeriods)
}

func (h *ExamPeriodHandler) UpsertExamPeriod(c *gin.Context) {
	var examPeriod *models.ExamPeriod
	if err := c.ShouldBindJSON(&examPeriod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.examPeriodService.UpsertExamPeriod(examPeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, examPeriod)
}

func (h *ExamPeriodHandler) DeleteExamPeriod(c *gin.Context) {
	examPeriodID := c.Param("examPeriodID")
	examPeriodIDInt, err := strconv.ParseInt(examPeriodID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.examPeriodService.DeleteExamPeriod(examPeriodIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Exam period deleted successfully"})
}
