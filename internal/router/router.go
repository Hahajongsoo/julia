package router

import (
	"julia/internal/di"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, c *di.Container) {
	user := router.Group("/users")
	{
		user.GET("/:id", c.UserHandler.GetByID)
		user.POST("/", c.UserHandler.Create)
		user.PUT("/:id", c.UserHandler.Update)
		user.DELETE("/:id", c.UserHandler.Delete)
	}
}
