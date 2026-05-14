package repository

import (
	"context"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	// "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// create a struct that holds the pgxpool.Pool

type LogRepository struct {
	dbPool *pgxpool.Pool
}

func NewLogRepository(pool *pgxpool.Pool) *LogRepository {
	return &LogRepository{dbPool: pool}
}

// Task: Try to implement a method InsertLogBatch in your repository
func (lr *LogRepository) InsertLogBatch(ctx context.Context, logs []models.LogEntry) error {
	/*
		* context.Context parameter lets the caller:
			cancel the DB operation
			enforce timeouts
			propagate request-scoped values
			stop work if the HTTP request/client disconnects
		Without context, your DB query could continue running even when nobody cares about the result anymore.
	*/
	//* Return an error when the function(even void func like this) can realistically fail in a way the caller should handle.

	// step 1: Specify the columns needed
	columns := []string{"level", "ts", "caller", "message"}

	// step 2: Use copyfrom to stream data into the database in very high speed through a special pipeline

	_, err := lr.dbPool.CopyFrom(
		ctx,                    // context
		pgx.Identifier{"logs"}, // table name
		columns,                // columns
		pgx.CopyFromSlice(len(logs), func(i int) ([]any, error) { // row source
			// This function tells Go how to map your LogEntry to the columns we defined above.
			return []any{logs[i].Level, logs[i].Timestamp, logs[i].Caller, logs[i].Message}, nil
		}),
	)
	return err

}
