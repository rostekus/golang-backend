package db

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type dbPostgresConfig struct {
	User     string `mapstructure:"POSTGRES_USER" validate:"required"`
	Password string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
	Host     string `mapstructure:"POSTGRES_HOST" validate:"required"`
	Name     string `mapstructure:"POSTGRES_DB" validate:"required"`
	Port     string `mapstructure:"POSTGRES_PORT"`
}

func (config *dbPostgresConfig) PostgresDSNFromConfig() string {

	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	return dns
}

func loadPostgresDBConfig() (config dbPostgresConfig, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("POSTGRES_DB_PORT", "5431")
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

func NewPostgresDatabase() (*PostgresDatabase, error) {
	configDB, err := loadPostgresDBConfig()
	if err != nil {
		panic("Cannot load postgres database config")
	}

	dsn := configDB.PostgresDSNFromConfig()
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("Cannot connect to database")
	}
	return &PostgresDatabase{db: db}, nil
}

func (d *PostgresDatabase) Close() {
	d.db.Close()
}

func (d *PostgresDatabase) GetDB() *sql.DB {
	return d.db
}
