package di

import (
	"database/sql"
	"julia/internal/handlers"
	"julia/internal/repositories"
	"julia/internal/services"
)

type Container struct {
	UserHandler *handlers.UserHandler
}

func NewContainer(db *sql.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	userSvc := services.NewUserService(userRepo)
	userHdl := handlers.NewUserHandler(userSvc)

	return &Container{
		UserHandler: userHdl,
	}
}
