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
	{
		user.GET("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.UserHandler.GetAll)
		user.GET("/:id", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.UserHandler.GetByID)
		user.GET("/:id/todos", c.UserHandler.GetTodosByUserID)
		user.POST("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.UserHandler.Create)
		user.PUT("/:id", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.UserHandler.Update)
		user.DELETE("/:id", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.UserHandler.Delete)
	}
	classes := router.Group("/classes")
	classes.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		classes.GET("", c.ClassHandler.GetAllClasses)
		classes.GET("/:classID", c.ClassHandler.GetClassByID)
		classes.GET("/:classID/users", c.ClassHandler.GetUsersByClassID)
		classes.POST("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ClassHandler.CreateClass)
		classes.PUT("/:classID", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ClassHandler.UpdateClass)
		classes.DELETE("/:classID", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ClassHandler.DeleteClass)
		classes.GET("/assignments", c.AssignmentHandler.GetAssignmentWithClassID)
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
	examPeriod := router.Group("/exam-periods")
	examPeriod.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		examPeriod.GET("", c.ExamPeriodHandler.GetAllExamPeriods)
		examPeriod.GET("/class/:classID", c.ExamPeriodHandler.GetExamPeriodByClassID)
		examPeriod.POST("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ExamPeriodHandler.CreateExamPeriod)
		examPeriod.PUT("/:examPeriodID", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ExamPeriodHandler.UpdateExamPeriod)
		examPeriod.DELETE("/:examPeriodID", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.ExamPeriodHandler.DeleteExamPeriod)
	}

	todos := router.Group("/todos")
	todos.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		todos.GET("", c.TodoHandler.GetAllTodos)
		todos.POST("", c.TodoHandler.CreateTodo)
		todos.PUT("/:id", c.TodoHandler.UpdateTodo)
		todos.DELETE("/:id", c.TodoHandler.DeleteTodo)
	}
	assignment := router.Group("/assignments")
	assignment.Use(middlewares.AuthMiddleware(c.AuthService))
	{
		assignment.GET("", c.AssignmentHandler.GetAllAssignments)
		assignment.GET("/user/:userID", c.AssignmentHandler.GetAssignmentsByUserID)
		assignment.GET("/makeup/:makeupID", c.AssignmentHandler.GetAssignmentsByMakeupID)
		assignment.POST("", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.AssignmentHandler.UpsertAssignment)
		assignment.DELETE("/:assignmentID", middlewares.AdminAuthMiddleware(c.AuthService, c.UserService), c.AssignmentHandler.DeleteAssignment)
	}
}
