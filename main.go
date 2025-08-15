package main

import (
	"julia/config"
	"julia/internal/di"
	"julia/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	container := di.NewContainer(db)

	r := gin.Default()
	router.SetupRouter(r, container)
	r.Run(":8080")
}
