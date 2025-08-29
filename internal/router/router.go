package router

import (
	"julia/internal/di"
	"julia/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, c *di.Container) {
	router.Use(middlewares.CORSMiddleware())
	user := router.Group("/users")
	// user.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		user.GET("/:id", c.UserHandler.GetByID)
		user.POST("", c.UserHandler.Create)
		user.PUT("/:id", c.UserHandler.Update)
		user.DELETE("/:id", c.UserHandler.Delete)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.LoginHandler.Login)
		auth.POST("/logout", middlewares.AuthMiddleware(c.AuthService), c.LoginHandler.Logout)
	}
	makeup := router.Group("/makeups")
	{
		makeup.GET("", c.MakeupHandler.GetAllMakeups)                                      
		makeup.GET("/month/:yearMonth", c.MakeupHandler.GetMakeupsByMonth)                  
		makeup.GET("/date/:date", c.MakeupHandler.GetMakeupsByDate)                         
		makeup.GET("/user/:userID", c.MakeupHandler.GetMakeupsByUser)                      
		makeup.GET("/user/:userID/date/:date", c.MakeupHandler.GetMakeupsByUserAndDate)    
		makeup.POST("", c.MakeupHandler.CreateMakeup)                                      
		makeup.PUT("/user/:userID/date/:date/time/:time", c.MakeupHandler.UpdateMakeup)    
		makeup.DELETE("/user/:userID/date/:date/time/:time", c.MakeupHandler.DeleteMakeup) 
	}
}
