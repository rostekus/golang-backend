package db

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbMongoConfig struct {
	User     string `mapstructure:"MONGO_INITDB_ROOT_USERNAME" validate:"required"`
	Password string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD" validate:"required"`
	Host     string `mapstructure:"MONGO_HOST" validate:"required"`
	Port     string `mapstructure:"MONGO_PORT"`
	Name     string `mapstructure:"MONGO_INITDB_DATABASE" validate:"required"`
}

func loadMongoConfig() (config dbMongoConfig, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("MONGO_PORT", "27017")
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

func (configDB *dbMongoConfig) MongoURIFromConfig() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s",
		configDB.User,
		configDB.Password,
		configDB.Host,
		configDB.Port)
}

func NewMongoClient(databaseName string) (*MongoDatabase, error) {
	configDB, err := loadMongoConfig()
	if err != nil {
		panic("Cannot load database config")
	}

	uri := configDB.MongoURIFromConfig()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &MongoDatabase{client: client, timeout: time.Duration(10) * time.Second}, nil

}

func (mongoDB *MongoDatabase) UploadFile(fileReader multipart.File, filename string, databaseName string) error {
	bucket, err := gridfs.NewBucket(
		mongoDB.client.Database(databaseName),
	)
	if err != nil {
		return err
	}
	uploadStream, err := bucket.OpenUploadStream(
		filename,
	)
	if err != nil {
		return err
	}
	defer uploadStream.Close()
	_, err = io.Copy(uploadStream, fileReader)
	if err != nil {
		return err
	}

	// Close the upload stream
	err = uploadStream.Close()
	if err != nil {
		return err
	}

	return nil
}

func (mongoDB *MongoDatabase) DownloadFile(filename string, databaseName string) (*bytes.Buffer, error) {
	db := mongoDB.client.Database(databaseName)
	fsFiles := db.Collection("fs.files")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var results bson.M
	err := fsFiles.FindOne(ctx, bson.M{}).Decode(&results)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// you can print out the results
	fmt.Println(results)

	bucket, _ := gridfs.NewBucket(
		db,
	)

	var buf bytes.Buffer
	_, err = bucket.DownloadToStreamByName("test.png", &buf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &buf, nil
}
