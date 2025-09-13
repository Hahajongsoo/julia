package di

import (
	"database/sql"
	"julia/internal/handlers"
	"julia/internal/repositories"
	"julia/internal/services"
	"net/http"
	"os"
	"strconv"
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

// Helper functions for environment variable parsing
func getDurationFromEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringFromEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolFromEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getSameSiteFromEnv(key string, defaultValue http.SameSite) http.SameSite {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "strict":
			return http.SameSiteStrictMode
		case "lax":
			return http.SameSiteLaxMode
		case "none":
			return http.SameSiteNoneMode
		}
	}
	return defaultValue
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
	if os.Getenv("HMAC_SECRET") == "" {
		panic("HMAC_SECRET is not set")
	}
	authSvc := services.NewAuthService(userRepo, services.Config{
		SessionTTL: getDurationFromEnv("SESSION_TTL", 30*time.Minute),
		HMACSecret: []byte(os.Getenv("HMAC_SECRET")),
		CookieName: getStringFromEnv("COOKIE_NAME", "my-session"),
		CookiePath: getStringFromEnv("COOKIE_PATH", "/"),
		Secure:     getBoolFromEnv("COOKIE_SECURE", true),
		HttpOnly:   getBoolFromEnv("COOKIE_HTTP_ONLY", true),
		SameSite:   getSameSiteFromEnv("COOKIE_SAME_SITE", http.SameSiteLaxMode),
		Domain:     os.Getenv("COOKIE_DOMAIN"),
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
