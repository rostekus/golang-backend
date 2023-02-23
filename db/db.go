package db

import (
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
