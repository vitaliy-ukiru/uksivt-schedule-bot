package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	DefaultIP   = "localhost"
	DefaultPort = 5432
)

func ConnString(user, pass, db, ip string, port int) string {
	if port == 0 {
		port = DefaultPort
	}
	if ip == "" {
		ip = DefaultIP
	}
	return fmt.Sprintf("postgres://%s:%s@%v:%v/%v?sslmode=disable",
		user,
		pass,
		ip,
		port,
		db)
}

type OptionFunc func(config *pgxpool.Config)

func ParseConfig(connString string, options ...OptionFunc) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if len(options) != 0 {
		for _, option := range options {
			option(config)
		}
	}

	return config, nil
}

func New(ctx context.Context, connString string, options ...OptionFunc) (*pgxpool.Pool, error) {
	config, err := ParseConfig(connString, options...)
	if err != nil {
		return nil, err
	}

	return pgxpool.ConnectConfig(ctx, config)
}
