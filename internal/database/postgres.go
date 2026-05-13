package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Write a function that returns a *pgxpool.Pool
type PostgresSetup struct {
	logger *zap.Logger
}

func NewPostgresSetup(l *zap.Logger) *PostgresSetup {
	return &PostgresSetup{logger: l}
}

func (pt *PostgresSetup) ConnectDB() (*pgxpool.Pool, error) {
	db, nonEmpty := os.LookupEnv("DB_URL")
	if !nonEmpty {
		// it is advisable to use Error instead of Panic and let main.go(or the caller) to handle it
		// this make this code more reusable
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
	// defer dbPool.Close() in main function
}
