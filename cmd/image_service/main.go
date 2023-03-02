package main

import (
	"fmt"
	"log"
	"rostekus/golang-backend/db"
	"rostekus/golang-backend/internal/image"
	"rostekus/golang-backend/middleware"
	"rostekus/golang-backend/rabbitmq"
	"rostekus/golang-backend/util"

	"github.com/gin-gonic/gin"
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
	log.Println("Connected to RabbitMQ")
	postgres, err := db.NewPostgresDatabase()
	if err != nil {
		log.Fatalf("Could not initialize database connection")
	}
	log.Println("Connected to Postgres")
	defer postgres.Close()

	rep := image.NewRepository(mongo, "golang")
	srv := image.NewService(rep, rabbitmq, postgres.GetDB())
	handler := image.NewHandler(srv)
	router := gin.Default()
	router.POST("/upload", middleware.JWTAuth, handler.UploadFile)
	router.GET("/download", middleware.JWTAuth, handler.DownloadFile)
	router.GET("/images", middleware.JWTAuth, handler.GetImagesForUser)
	router.Run(":23451")

	defer rabbitmq.Close()

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	message := time.Now().String()
	// 	rabbitmq.Publish(message)
	// 	fmt.Fprint(w, "Message published")
	// })
	// http.ListenAndServe(":23451", nil)

}
