package utils

import (
	"errors"
	"github.com/RomanV1/go-rest-api/pkg/responses"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

func HandleUniqueViolation(err error, c *gin.Context) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_email_key":
			c.JSON(http.StatusBadRequest, responses.EmailAlreadyExists)
		case "users_username_key":
			c.JSON(http.StatusBadRequest, responses.UsernameAlreadyExists)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Unique constraint violation"})
		}
		return true
	}
	return false
}
