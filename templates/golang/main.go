package main

import (
	"fmt"

	"github.com/Autodoc-Technology/interview-templates/template/golang/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("error getting config", zap.Error(err))
	}

	fmt.Println("application entry point")
}
