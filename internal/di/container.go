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
	UserHandler  *handlers.UserHandler
	LoginHandler *handlers.LoginHandler
	AuthService  services.AuthService
}

func NewContainer(db *sql.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userHdl := handlers.NewUserHandler(userSvc)

	authSvc := services.NewAuthService(userRepo, services.Config{
		SessionTTL: 30 * time.Minute,
		HMACSecret: []byte(os.Getenv("HMAC_SECRET")),
		CookieName: "my-session",
		CookiePath: "/",
		Secure:     false,
		SameSite:   http.SameSiteLaxMode,
	})
	return &Container{
		UserHandler:  userHdl,
		LoginHandler: handlers.NewLoginHandler(authSvc),
		AuthService:  authSvc,
	}
}
