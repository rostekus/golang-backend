package main

import (
	"fmt"
	"log"
	"rostekus/golang-backend/internal/health"
	"rostekus/golang-backend/internal/image"
	"rostekus/golang-backend/internal/router"
	"rostekus/golang-backend/internal/server"
	"rostekus/golang-backend/pkg/db"
	"rostekus/golang-backend/pkg/rabbitmq"
	"rostekus/golang-backend/pkg/util"
)

func main() {
	mongo, err := db.NewMongoClient("golang")
	if err != nil {
		panic(err.Error())
	}
	log.Println("Connected to MongoDB")

	config := util.ReadConfig()
	fmt.Println(config.QueueName)
	rabbitmq := rabbitmq.NewRabbitMQ()
	rabbitmq.CreateChannel()
	rabbitmq.QueueDeclare("video_queue")
	defer rabbitmq.Close()
	log.Println("Connected to RabbitMQ")
	postgres, err := db.NewPostgresDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection")
	}
	log.Println("Connected to Postgres")
	defer postgres.Close()

	rep := image.NewRepository(mongo, "golang")
	imageService := image.NewService(rep, rabbitmq, postgres.GetDB())
	imageHandler := image.NewHandler(imageService)
	healthHandler := health.NewHandler()
	router := router.NewImageServiceRouter(imageHandler, healthHandler)

	srv := server.NewServer("23451", router.Router)
	srv.Run()
}
