package db

import (
	"context"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitPool() error {
	config, err := pgxpool.ParseConfig(os.Getenv("DB_CONNECTION"))
	if err != nil {
		return err
	}

	maxc, err := strconv.Atoi(os.Getenv("DB_MAX_CONN"))
	if err != nil {
		return err
	}
	minc, err := strconv.Atoi(os.Getenv("DB_MIN_CONN"))
	if err != nil {
		return err
	}

	config.MaxConns = int32(maxc)
	config.MinConns = int32(minc)

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	db = pool
	return nil
}

func GetConnection() (*pgxpool.Conn, error) {
	if db == nil {
		InitPool()
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
