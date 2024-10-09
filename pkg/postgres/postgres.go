package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func NewClient(ctx context.Context, username, password, host, port, database string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("PostgreSQL ping failed: %w", err)
	}

	return conn, nil
}
