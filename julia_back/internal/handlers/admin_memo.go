package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMemoHandler struct {
	memoService services.AdminMemoService
}

func NewAdminMemoHandler(memoService services.AdminMemoService) *AdminMemoHandler {
	return &AdminMemoHandler{
		memoService: memoService,
	}
}

func (h *AdminMemoHandler) GetMemo(c *gin.Context) {
	memo, err := h.memoService.GetAdminMemo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memo)
}

func (h *AdminMemoHandler) SaveMemo(c *gin.Context) {
	var memo models.AdminMemo
	if err := c.ShouldBindJSON(&memo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	savedMemo, err := h.memoService.SaveAdminMemo(&memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, savedMemo)
}

func (h *AdminMemoHandler) DeleteMemo(c *gin.Context) {
	err := h.memoService.DeleteAdminMemo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Memo deleted successfully"})
}
