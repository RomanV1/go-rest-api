package users

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"

	"time"
)

type Repository interface {
	GetOne(ctx context.Context, uuid uuid.UUID) (User, error)
	GetAll(ctx context.Context, limit, offset int) ([]User, error)
	Create(ctx context.Context, dto CreateUserDto) (User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	Update(ctx context.Context, uuid uuid.UUID, dto UpdateUserDto) (User, error)
}

type db struct {
	connect pgx.Conn
	logger  *logrus.Logger
}

func NewRepository(connect pgx.Conn, logger *logrus.Logger) Repository {
	return &db{
		connect: connect,
		logger:  logger,
	}
}

func (d *db) GetOne(ctx context.Context, uuid uuid.UUID) (User, error) {
	query := `SELECT uuid, username, email, password_hash, created_at, updated_at FROM users WHERE uuid = $1`

	var user User

	row := d.connect.QueryRow(ctx, query, uuid)
	err := row.Scan(&user.UUID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			d.logger.WithField("userID", uuid).Warn("User not found")
			return User{}, sql.ErrNoRows
		}

		d.logger.WithError(err).WithField("userID", uuid).Error("Failed to get user by ID")
		return User{}, err
	}

	return user, nil
}

func (d *db) GetAll(ctx context.Context, limit, offset int) ([]User, error) {
	if limit <= 0 {
		limit = 1000000
	}

	if offset <= 0 {
		offset = 0
	}

	query := `SELECT uuid, username, email, password_hash, created_at, updated_at FROM users ORDER BY uuid LIMIT $1 OFFSET $2`

	rows, err := d.connect.Query(ctx, query, limit, offset)
	if err != nil {
		d.logger.WithError(err).Error("Failed to execute query to get users")
		return nil, err
	}

	defer rows.Close()

	var allUsers []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
			d.logger.WithError(err).Error("Failed to scan user row")
			return nil, err
		}

		allUsers = append(allUsers, user)
	}
	if err := rows.Err(); err != nil {
		d.logger.WithError(err).Error("Error encountered while iterating through rows")
		return nil, err
	}

	return allUsers, nil
}

func (d *db) Create(ctx context.Context, dto CreateUserDto) (User, error) {
	query := `INSERT INTO users (username, email, password_hash) 
			  VALUES ($1, $2, $3)
			  RETURNING uuid, username, email, password_hash, created_at, updated_at`

	var user User

	row := d.connect.QueryRow(ctx, query, dto.Username, dto.Email, dto.Password)
	err := row.Scan(&user.UUID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		d.logger.WithError(err).Error("Failed to create new user")
		return User{}, err
	}

	return user, nil
}

func (d *db) Update(ctx context.Context, uuid uuid.UUID, dto UpdateUserDto) (User, error) {
	var queryParts []string
	args := []interface{}{uuid}

	if dto.Username != "" {
		queryParts = append(queryParts, "username = $"+strconv.Itoa(len(args)+1))
		args = append(args, dto.Username)
	}
	if dto.Email != "" {
		queryParts = append(queryParts, "email = $"+strconv.Itoa(len(args)+1))
		args = append(args, dto.Email)
	}
	if dto.Password != "" {
		queryParts = append(queryParts, "password_hash = $"+strconv.Itoa(len(args)+1))
		args = append(args, dto.Password)
	}

	queryParts = append(queryParts, "updated_at = $"+strconv.Itoa(len(args)+1))
	args = append(args, time.Now())

	query := `UPDATE users SET ` + strings.Join(queryParts, ", ") + `
              WHERE uuid = $1
			  RETURNING uuid, username, email, password_hash, created_at, updated_at`

	var user User

	row := d.connect.QueryRow(ctx, query, args...)
	err := row.Scan(&user.UUID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			d.logger.WithField("userID", uuid).Warn("User not found for update")
			return User{}, sql.ErrNoRows
		}
		d.logger.WithError(err).Error("Failed to update user")
		return User{}, err
	}

	return user, nil
}

func (d *db) Delete(ctx context.Context, uuid uuid.UUID) error {
	query := `DELETE FROM users WHERE uuid = $1`

	res, err := d.connect.Exec(ctx, query, uuid)
	if err != nil {
		d.logger.WithError(err).Error("Failed to execute delete query")
		return err
	}

	if res.RowsAffected() == 0 {
		d.logger.WithField("userID", uuid).Warn("No user found to delete")
		return sql.ErrNoRows
	}

	return nil
}
