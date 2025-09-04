package router

import (
	"julia/internal/di"
	"julia/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, c *di.Container) {
	router.Use(gin.Recovery())
	router.Use(middlewares.CustomLoggerMiddleware())
	router.Use(middlewares.CORSMiddleware())

	user := router.Group("/users")
	user.Use(middlewares.AuthMiddleware(c.AuthService))
	user.Use(middlewares.AdminAuthMiddleware(c.AuthService, c.UserService))
	{
		user.GET("", c.UserHandler.GetAll)
		user.GET("/:id", c.UserHandler.GetByID)
		user.POST("", c.UserHandler.Create)
		user.PUT("/:id", c.UserHandler.Update)
		user.DELETE("/:id", c.UserHandler.Delete)
	}
	classes := router.Group("/classes")
	classes.Use(middlewares.AuthMiddleware(c.AuthService))
	classes.Use(middlewares.AdminAuthMiddleware(c.AuthService, c.UserService))
	{
		classes.GET("", c.ClassHandler.GetAllClasses)
		classes.GET("/:classID", c.ClassHandler.GetClassByID)
		classes.POST("", c.ClassHandler.CreateClass)
		classes.PUT("/:classID", c.ClassHandler.UpdateClass)
		classes.DELETE("/:classID", c.ClassHandler.DeleteClass)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/login", c.LoginHandler.Login)
		auth.POST("/logout", middlewares.AuthMiddleware(c.AuthService), c.LoginHandler.Logout)
		auth.GET("/me", middlewares.AuthMiddleware(c.AuthService), c.LoginHandler.GetCurrentUser)
	}
	makeup := router.Group("/makeups")
	makeup.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		makeup.GET("", c.MakeupHandler.GetAllMakeups)
		makeup.GET("/month/:yearMonth", c.MakeupHandler.GetMakeupsByMonth)
		makeup.GET("/date/:date", c.MakeupHandler.GetMakeupsByDate)
		makeup.GET("/user/:userID", c.MakeupHandler.GetMakeupsByUser)
		makeup.GET("/user/:userID/date/:date", c.MakeupHandler.GetMakeupsByUserAndDate)
		makeup.POST("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.MakeupHandler.CreateMakeup)
		makeup.PUT("/user/:userID/date/:date/time/:time", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.MakeupHandler.UpdateMakeup)
		makeup.DELETE("/user/:userID/date/:date/time/:time", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.MakeupHandler.DeleteMakeup)
	}
	push := router.Group("/push")
	push.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		push.GET("/vapid-public", c.PushHandler.GetVapidPublic)
		push.POST("/subscriptions", c.PushHandler.CreateSubscription)
		push.DELETE("/subscriptions/:userID", c.PushHandler.DeleteSubscription)
		push.GET("/subscriptions/:userID", c.PushHandler.GetSubscriptions)
		push.POST("/notifications/:userID", c.PushHandler.SendNotification)
	}
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
}
