package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/RomanV1/go-rest-api/internal/utils"
	"github.com/RomanV1/go-rest-api/pkg/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type Handler interface {
	Register(router *gin.RouterGroup)
}

type handler struct {
	service   Service
	validator validator.Validate
	logger    *logrus.Logger
}

func NewHandler(service Service, valid validator.Validate, logger *logrus.Logger) Handler {
	return &handler{
		service:   service,
		validator: valid,
		logger:    logger,
	}
}

func (h *handler) Register(router *gin.RouterGroup) {
	router.GET(usersURL, h.GetUsers)
	router.GET(userURL, h.GetUserByID)
	router.POST(usersURL, h.CreateUser)
	router.PUT(userURL, h.UpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetUserByID(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		h.logger.WithError(err).WithField("userUUID", c.Param("uuid")).Warn("Invalid UUID format")
		c.JSON(http.StatusBadRequest, responses.InvalidIDParam)
		return
	}

	user, err := h.service.GetUserByID(context.Background(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(404, responses.UserNotFound)
			return
		}

		c.JSON(http.StatusInternalServerError, responses.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) GetUsers(c *gin.Context) {
	limit, err := utils.ParseQueryParam(c.Query("limit"), 1000000)
	if err != nil {
		h.logger.WithError(err).WithField("limit", c.Query("limit")).Warn("Invalid limit parameter")
		c.JSON(http.StatusBadRequest, responses.InvalidLimitParam)
		return
	}

	offset, err := utils.ParseQueryParam(c.Query("offset"), 0)
	if err != nil {
		h.logger.WithError(err).WithField("offset", c.Query("offset")).Warn("Invalid offset parameter")
		c.JSON(http.StatusBadRequest, responses.InvalidOffsetParam)
		return
	}

	users, err := h.service.GetAllUsers(context.Background(), limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, responses.UserNotFound)
			return
		}

		c.JSON(http.StatusInternalServerError, responses.ErrorRetrievingUsers)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *handler) CreateUser(c *gin.Context) {
	var userDto CreateUserDto
	if err := c.BindJSON(&userDto); err != nil {
		h.logger.WithError(err).Warn("Failed to parse request body")
		c.JSON(http.StatusBadRequest, responses.ErrorParsingJSON)
		return
	}

	if err := h.validator.Struct(userDto); err != nil {
		customErrMsg := utils.CustomErrorMessages(err)
		h.logger.WithField("validationErrors", customErrMsg).Warn("User validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": customErrMsg,
		})
		return
	}

	user, err := h.service.CreateUser(context.Background(), CreateUserDto{
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: userDto.Password,
	})
	if err != nil {
		if utils.HandleUniqueViolation(err, c) {
			h.logger.WithError(err).Warn("User creation failed due to unique constraint violation")
			return
		}

		c.JSON(http.StatusBadRequest, responses.UserCreationError)
		return
	}

	c.JSON(201, user)
}

func (h *handler) UpdateUser(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		h.logger.WithError(err).WithField("userUUID", c.Param("uuid")).Warn("Invalid UUID format")
		c.JSON(http.StatusBadRequest, responses.InvalidIDParam)
		return
	}

	var userDto UpdateUserDto
	if err := c.BindJSON(&userDto); err != nil {
		h.logger.WithError(err).WithField("userUUID", userUUID).Warn("Failed to parse request body")
		c.JSON(http.StatusBadRequest, responses.ErrorParsingJSON)
		return
	}

	if err := h.validator.Struct(userDto); err != nil {
		customErrMsg := utils.CustomErrorMessages(err)
		h.logger.WithField("validationErrors", customErrMsg).Warn("User validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": customErrMsg,
		})
		return
	}

	user, err := h.service.UpdateUser(context.Background(), userUUID, userDto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, responses.UserNotFound)
			return
		}

		c.JSON(http.StatusBadRequest, responses.UserUpdateError)
		return
	}

	c.JSON(200, user)
}

func (h *handler) DeleteUser(c *gin.Context) {
	userUUID, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		h.logger.WithError(err).WithField("userUUID", c.Param("uuid")).Warn("Invalid UUID format")
		c.JSON(http.StatusBadRequest, responses.InvalidIDParam)
		return
	}

	err = h.service.DeleteUser(context.Background(), userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, responses.UserNotFound)
			return
		}

		c.JSON(http.StatusInternalServerError, responses.UserDeleteError)
		return
	}

	c.JSON(200, responses.UserDeletionSuccess)
}
