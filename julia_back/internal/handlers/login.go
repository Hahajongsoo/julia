package handlers

import (
	"julia/internal/models"
	"julia/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	authService services.AuthService
}

func NewLoginHandler(authService services.AuthService) *LoginHandler {
	return &LoginHandler{authService: authService}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var input *models.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, err := h.authService.VerifyCredential(input.ID, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	sid, session, err := h.authService.CreateSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.authService.IssueSessionCookie(c.Writer, sid)
	c.JSON(http.StatusOK, gin.H{
		"message":    "login successful",
		"expires_at": session.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *LoginHandler) Logout(c *gin.Context) {
	if sid, ok := c.Get("sid"); ok {
		h.authService.DeleteSession(sid.(string))
	}
	h.authService.ClearSessionCookie(c.Writer)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// GetCurrentUser returns the current user information
func (h *LoginHandler) GetCurrentUser(c *gin.Context) {
	sessionValue, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session not found"})
		return
	}
	
	session, ok := sessionValue.(*services.Session)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session type"})
		return
	}
	
	user, err := h.authService.GetUserByID(session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, user.ToResponseUser())
}
