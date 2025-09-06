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

func (h *ExamPeriodHandler) GetAllExamPeriods(c *gin.Context) {
	examPeriods, err := h.examPeriodService.GetAllExamPeriods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	responseExamPeriods := make([]*models.ExamPeriodDTO, len(examPeriods))
	for i, examPeriod := range examPeriods {
		responseExamPeriods[i] = examPeriod.ToExamPeriodDTO()
	}
	c.JSON(http.StatusOK, responseExamPeriods)
}

func (h *ExamPeriodHandler) GetExamPeriodByClassID(c *gin.Context) {
	classID := c.Param("classID")
	classIDInt, err := strconv.ParseInt(classID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	examPeriods, err := h.examPeriodService.GetExamPeriodByClassID(classIDInt)
	responseExamPeriods := make([]*models.ExamPeriodDTO, len(examPeriods))
	for i, examPeriod := range examPeriods {
		responseExamPeriods[i] = examPeriod.ToExamPeriodDTO()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseExamPeriods)
}

func (h *ExamPeriodHandler) CreateExamPeriod(c *gin.Context) {
	var dto *models.ExamPeriodDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	examPeriod := dto.ToExamPeriod()
	if examPeriod.StartDate.After(examPeriod.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return
	}
	err := h.examPeriodService.CreateExamPeriod(examPeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, examPeriod)
}

func (h *ExamPeriodHandler) UpdateExamPeriod(c *gin.Context) {
	var dto *models.ExamPeriodDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	examPeriod := dto.ToExamPeriod()
	if err := h.examPeriodService.UpdateExamPeriod(examPeriod); err != nil {
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
