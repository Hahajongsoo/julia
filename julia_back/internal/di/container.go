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
	AdminMemoHandler    *handlers.AdminMemoHandler
	AssignmentHandler   *handlers.AssignmentHandler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userHdl := handlers.NewUserHandler(userSvc)

	makeupRepo := repositories.NewMakeupRepository(db)
	makeupSvc := services.NewMakeupService(makeupRepo)
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationSvc := services.NewNotificationService(notificationRepo)
	makeupHdl := handlers.NewMakeupHandler(makeupSvc, notificationSvc)

	pushRepo := repositories.NewPushRepository(db)
	pushSvc := services.NewPushService(pushRepo)
	pushHdl := handlers.NewPushHandler(pushSvc)

	classRepo := repositories.NewClassRepository(db)
	classSvc := services.NewClassService(classRepo)
	classHdl := handlers.NewClassHandler(classSvc)

	examPeriodRepo := repositories.NewExamPeriodRepository(db)
	examPeriodSvc := services.NewExamPeriodService(examPeriodRepo)
	examPeriodHdl := handlers.NewExamPeriodHandler(examPeriodSvc)

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	adminMemoRepo := repositories.NewAdminMemoRepository(dataDir)
	adminMemoSvc := services.NewAdminMemoService(adminMemoRepo)
	adminMemoHdl := handlers.NewAdminMemoHandler(adminMemoSvc)

	authSvc := services.NewAuthService(userRepo, services.Config{
		SessionTTL: 30 * time.Minute,
		HMACSecret: []byte(os.Getenv("HMAC_SECRET")),
		CookieName: "my-session",
		CookiePath: "/",
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
	})

	loginHdl := handlers.NewLoginHandler(authSvc)

	assignmentRepo := repositories.NewAssignmentRepository(db)
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
		AdminMemoHandler:  adminMemoHdl,
		AssignmentHandler: assignmentHdl,
	}
}
