package main

import (
	"context"
	"fmt"

	"github.com/Autodoc-Technology/interview-templates/template/golang/config"
	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/handler"
	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/service"
	"github.com/Autodoc-Technology/interview-templates/template/golang/pkg/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("error getting config", zap.Error(err))
	}

	mongo, err := storage.NewMongoClient(ctx, cfg.MongoUrl, logger)
	if err != nil {
		logger.Fatal("error connecting to mongo", zap.Error(err))
	}

	services := service.NewService(logger, mongo)
	handlers := handler.NewHandler(logger, services)

	r := gin.Default()

	// endpoints

	api := r.Group("/api/v1")

	apiUser := api.Group("/user")
	apiUser.POST("/create", handlers.CreateUser)
	apiUser.POST("/get", handlers.GetUser)
	apiUser.DELETE("/delete", handlers.DeleteUser)
	apiUser.GET("/list", handlers.ListUsers)

	if err := r.Run(fmt.Sprintf("localhost:%d", cfg.Port)); err != nil {
		logger.Fatal("error r.Run", zap.Error(err))
	}
}
