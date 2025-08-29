package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MakeupHandler struct {
	makeupService services.MakeupService
}

func NewMakeupHandler(makeupService services.MakeupService) *MakeupHandler {
	return &MakeupHandler{makeupService: makeupService}
}

func (h *MakeupHandler) GetAllMakeups(c *gin.Context) {
	makeups, err := h.makeupService.GetAllMakeups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseMakeups := make([]*models.MakeupDTO, len(makeups))
	for i, makeup := range makeups {
		responseMakeups[i] = makeup.ToMakeupDTO()
	}

	c.JSON(http.StatusOK, responseMakeups)
}

func (h *MakeupHandler) GetMakeupsByDate(c *gin.Context) {
	date := c.Param("date")
	makeups, err := h.makeupService.GetMakeupsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseMakeups := make([]*models.MakeupDTO, len(makeups))
	for i, makeup := range makeups {
		responseMakeups[i] = makeup.ToMakeupDTO()
	}

	c.JSON(http.StatusOK, responseMakeups)
}

func (h *MakeupHandler) GetMakeupsByMonth(c *gin.Context) {
	yearMonth := c.Param("yearMonth")
	makeups, err := h.makeupService.GetMakeupsByMonth(yearMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseMakeups := make([]*models.MakeupDTO, len(makeups))
	for i, makeup := range makeups {
		responseMakeups[i] = makeup.ToMakeupDTO()
	}

	c.JSON(http.StatusOK, responseMakeups)
}

func (h *MakeupHandler) GetMakeupsByUser(c *gin.Context) {
	userID := c.Param("userID")
	makeups, err := h.makeupService.GetMakeupsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseMakeups := make([]*models.MakeupDTO, len(makeups))
	for i, makeup := range makeups {
		responseMakeups[i] = makeup.ToMakeupDTO()
	}

	c.JSON(http.StatusOK, responseMakeups)
}

func (h *MakeupHandler) GetMakeupsByUserAndDate(c *gin.Context) {
	userID := c.Param("userID")
	date := c.Param("date")
	makeups, err := h.makeupService.GetMakeupsByUserAndDate(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseMakeups := make([]*models.MakeupDTO, len(makeups))
	for i, makeup := range makeups {
		responseMakeups[i] = makeup.ToMakeupDTO()
	}

	c.JSON(http.StatusOK, responseMakeups)
}

func (h *MakeupHandler) CreateMakeup(c *gin.Context) {
	var input *models.MakeupDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.makeupService.CreateMakeup(input.ToMakeup())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}

func (h *MakeupHandler) UpdateMakeup(c *gin.Context) {
	userID := c.Param("userID")
	date := c.Param("date")
	time := c.Param("time")
	var input *models.MakeupDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.makeupService.UpdateMakeup(userID, date, time, input.ToMakeup())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, input)
}

func (h *MakeupHandler) DeleteMakeup(c *gin.Context) {
	userID := c.Param("userID")
	date := c.Param("date")
	time := c.Param("time")
	err := h.makeupService.DeleteMakeup(userID, date, time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Makeup deleted successfully"})
}
