package handlers

import (
	"fmt"
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

func (h *MakeupHandler) GetMakeupsByDate(c *gin.Context) {
	date := c.Param("date")
	makeups, err := h.makeupService.GetMakeupsByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, makeups)
}

func (h *MakeupHandler) GetMakeupsByIDandDate(c *gin.Context) {
	date := c.Param("date")
	userID := c.Param("userID")
	makeups, err := h.makeupService.GetMakeupsByIDandDate(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, makeups)
}

func (h *MakeupHandler) CreateMakeup(c *gin.Context) {
	var input *models.Makeup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(input)
	err := h.makeupService.CreateMakeup(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *MakeupHandler) UpdateMakeup(c *gin.Context) {
	date := c.Param("date")
	userID := c.Param("userID")
	time := c.Param("time")
	var input *models.Makeup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.makeupService.UpdateMakeup(userID, date, time, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *MakeupHandler) DeleteMakeup(c *gin.Context) {
	date := c.Param("date")
	userID := c.Param("userID")
	time := c.Param("time")
	err := h.makeupService.DeleteMakeup(userID, date, time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Makeup deleted successfully"})
}
