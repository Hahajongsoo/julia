package middlewares

import (
	"julia/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"julia/internal/models"
)

type SessionProvider interface {
	SessionName() string
	ParseAndVerifySID(raw string) (sid string, err error)
	GetSession(sid string) (*services.Session, bool)
}

type UserProvider interface {
	GetUserByID(userID string) (*models.User, error)
}

func AuthMiddleware(sp SessionProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(sp.SessionName())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		sid, err := sp.ParseAndVerifySID(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		session, ok := sp.GetSession(sid)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or not found"})
			return
		}
		c.Set("sid", sid)
		c.Set("session", session)
		c.Next()
	}
}

func AdminAuthMiddleware(sp SessionProvider, up UserProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, err := c.Cookie(sp.SessionName())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		sid, err := sp.ParseAndVerifySID(raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		session, ok := sp.GetSession(sid)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "session expired or not found"})
			return
		}

		user, err := up.GetUserByID(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get user information"})
			return
		}

		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		c.Set("sid", sid)
		c.Set("session", session)
		c.Set("user", user)
		c.Next()
	}
}
