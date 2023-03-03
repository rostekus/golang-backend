package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	godotenv.Load(".env")
}

type PostgresDatabase struct {
	db *sql.DB
}

type MongoDatabase struct {
	client  *mongo.Client
	timeout time.Duration
}

type SQLDatabase interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
