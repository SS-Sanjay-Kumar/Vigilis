package repository

import "github.com/jackc/pgx/v5/pgxpool"

// todo: create a struct that holds the pgxpool.Pool

type LogRepository struct{
	dbPool *pgxpool.Pool
}

func NewLogRepository(pool *pgxpool.Pool) *LogRepository {
	return &LogRepository{dbPool: pool}
}

