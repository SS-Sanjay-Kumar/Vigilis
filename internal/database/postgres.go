package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// todo: Write a function that returns a *pgxpool.Pool
type PostgresTools struct {
	logger *zap.Logger
}

func LoadPostgresTools(l *zap.Logger) *PostgresTools {
	return &PostgresTools{logger: l}
}

func (pt *PostgresTools) ConnectDB() (*pgxpool.Pool, error) {
	db, nonEmpty := os.LookupEnv("DB_URL")
	if !nonEmpty {
		pt.logger.Error("Missing ENV VARS")
		return nil, errors.New("Missing ENV VARS")
	}

	dbPool, err := pgxpool.New(context.Background(), db)
	if err != nil {
		pt.logger.Error(fmt.Sprintf("Unable to create connection pool: %v\n", err))
		return nil, err
	}
	checkDB := dbPool.Ping(context.Background())
	if checkDB != nil {
		pt.logger.Error("Database Ping Failed!!!")
		return nil, checkDB
	}

	return dbPool, nil
	// defer dbPool.Close()
}
