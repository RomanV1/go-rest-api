package main

import (
	"context"
	"github.com/RomanV1/go-rest-api/internal/users"
	"github.com/RomanV1/go-rest-api/pkg/postgres"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	router := gin.Default()

	api := router.Group("/api")

	client, err := postgres.NewClient(context.Background(), "postgres", "postgres", "localhost", "5432", "postgres")
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
