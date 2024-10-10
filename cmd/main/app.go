package main

import (
	"context"
	"github.com/RomanV1/go-rest-api/internal/users"
	"github.com/RomanV1/go-rest-api/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	loadEnv()

	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	router := gin.Default()

	api := router.Group("/api")

	client, err := postgres.NewClient(context.Background(), os.Getenv("PG_USERNAME"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_DATABASE"))
	if err != nil {
		logger.WithError(err).Error("Failed to create PostgreSQL client")
		panic(err)
	}

	repository := users.NewRepository(*client, logger)

	service := users.NewService(repository, logger)

	validate := validator.New()
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		u := sl.Current().Interface().(users.UpdateUserDto)
		if u.Username == "" && u.Email == "" && u.Password == "" {
			sl.ReportError("", "atLeastOneFieldRequired", "at_least_one_field_required", "at least one field must be present", "")
		}
	}, users.UpdateUserDto{})

	handler := users.NewHandler(service, *validate, logger)
	handler.Register(api)

	err = router.Run("localhost:3000")
	if err != nil {
		logger.WithError(err).Error("Failed to start the server")
		panic(err)
	}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
