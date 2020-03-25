package pg

import (
	"context"
	"github.com/aibotsoft/gproxy/internal/config"
	"github.com/jackc/pgx/v4"
	"time"
)

const connTimeout = 1 * time.Second

func init() {
	config.New()
}

func New() (*pgx.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	connConfig, err := pgx.ParseConfig("")
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
