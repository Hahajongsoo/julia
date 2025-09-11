package di

import (
	"database/sql"
	"julia/internal/handlers"
	"julia/internal/repositories"
	"julia/internal/services"
	"net/http"
	"os"
	"time"
)

type Container struct {
	UserHandler         *handlers.UserHandler
	LoginHandler        *handlers.LoginHandler
	AuthService         services.AuthService
	UserService         services.UserService
	MakeupHandler       *handlers.MakeupHandler
	PushHandler         *handlers.PushHandler
	NotificationService *services.NotificationService
	ClassHandler        *handlers.ClassHandler
	ExamPeriodHandler   *handlers.ExamPeriodHandler
	TodoHandler         *handlers.TodoHandler
	AssignmentHandler   *handlers.AssignmentHandler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	makeupRepo := repositories.NewMakeupRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	classRepo := repositories.NewClassRepository(db)
	pushRepo := repositories.NewPushRepository(db)
	todoRepo := repositories.NewTodoRepository(db)
	examPeriodRepo := repositories.NewExamPeriodRepository(db)
	assignmentRepo := repositories.NewAssignmentRepository(db)

	userSvc := services.NewUserService(userRepo, todoRepo)
	userHdl := handlers.NewUserHandler(userSvc)

	makeupSvc := services.NewMakeupService(makeupRepo)
	notificationSvc := services.NewNotificationService(notificationRepo)
	makeupHdl := handlers.NewMakeupHandler(makeupSvc, notificationSvc)

	pushSvc := services.NewPushService(pushRepo)
	pushHdl := handlers.NewPushHandler(pushSvc)

	classSvc := services.NewClassService(classRepo, userRepo)
	classHdl := handlers.NewClassHandler(classSvc)

	examPeriodSvc := services.NewExamPeriodService(examPeriodRepo)
	examPeriodHdl := handlers.NewExamPeriodHandler(examPeriodSvc)

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	todoSvc := services.NewTodoService(todoRepo)
	todoHdl := handlers.NewTodoHandler(todoSvc)

	authSvc := services.NewAuthService(userRepo, services.Config{
		SessionTTL: 30 * time.Minute,
		HMACSecret: []byte(os.Getenv("HMAC_SECRET")),
		CookieName: "my-session",
		CookiePath: "/",
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
	})

	loginHdl := handlers.NewLoginHandler(authSvc)

	assignmentSvc := services.NewAssignmentService(assignmentRepo)
	assignmentHdl := handlers.NewAssignmentHandler(assignmentSvc)

	return &Container{
		UserHandler:       userHdl,
		LoginHandler:      loginHdl,
		AuthService:       authSvc,
		UserService:       userSvc,
		MakeupHandler:     makeupHdl,
		PushHandler:       pushHdl,
		ClassHandler:      classHdl,
		ExamPeriodHandler: examPeriodHdl,
		TodoHandler:       todoHdl,
		AssignmentHandler: assignmentHdl,
	}
}
