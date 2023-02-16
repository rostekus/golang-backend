package db

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	godotenv.Load(".env")
}

type Database struct {
	db *sql.DB
}

type dbConfig struct {
	User     string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Host     string `mapstructure:"HOST" validate:"required"`
	Name     string `mapstructure:"POSTGRES_DB" validate:"required"`
	Port     string `mapstructure:"DB_PORT"`
}

func (config *dbConfig) DSNFromConfig() string {

	dns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)
	return dns
}

func LoadDBConfig() (config dbConfig, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("DB_PORT", "3306")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}
	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		return
	}
	return
}

func NewDatabase() (*Database, error) {
	configDB, err := LoadDBConfig()
	if err != nil {
		panic("Cannot load database config")
	}

	dsn := configDB.DSNFromConfig()
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("Cannot connect to database")
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
