package main

import (
	"context"
	"julia/config"
	"julia/internal/di"
	"julia/internal/router"
	"julia/workers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	worker := workers.Worker{DB: db}
	go worker.Start(context.Background())

	container := di.NewContainer(db)
	r := gin.Default()
	router.SetupRouter(r, container)
	r.Run(":8080")
}