package db

import (
	"database/sql"

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
	client *mongo.Client
}

type dbPostgresConfig struct {
	User     string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Host     string `mapstructure:"POSTGRES_HOST" validate:"required"`
	Name     string `mapstructure:"POSTGRES_DB" validate:"required"`
	Port     string `mapstructure:"POSTGRES_PORT"`
}

type dbMongoConfig struct {
	User     string `mapstructure:"MONGO_INITDB_ROOT_USERNAME" validate:"required"`
	Password string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD" validate:"required"`
	Host     string `mapstructure:"MONGO_HOST" validate:"required"`
	Port     string `mapstructure:"MONGO_PORT"`
	Name     string `mapstructure:"MONGO_INITDB_DATABASE" validate:"required"`
}
