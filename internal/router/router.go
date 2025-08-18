package router

import (
	"julia/internal/di"
	"julia/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, c *di.Container) {
	user := router.Group("/users")
	user.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		user.GET("/:id", c.UserHandler.GetByID)
		user.POST("/", c.UserHandler.Create)
		user.PUT("/:id", c.UserHandler.Update)
		user.DELETE("/:id", c.UserHandler.Delete)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.LoginHandler.Login)
		auth.POST("/logout", middlewares.AuthMiddleware(c.AuthService), c.LoginHandler.Logout)
	}
}
