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

	authSvc := services.NewAuthService(userRepo, services.Config{
		SessionTTL: 30 * time.Minute,
		HMACSecret: []byte(os.Getenv("HMAC_SECRET")),
		CookieName: "my-session",
		CookiePath: "/",
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
	})
	return &Container{
		UserHandler:   userHdl,
		LoginHandler:  handlers.NewLoginHandler(authSvc),
		AuthService:   authSvc,
		UserService:   userSvc,
		MakeupHandler: makeupHdl,
		PushHandler:   pushHdl,
	}
}
