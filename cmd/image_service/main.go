package main

import (
	"fmt"
	"net/http"
	"rostekus/golang-backend/internal/image"
	"rostekus/golang-backend/rabbitmq"
	"rostekus/golang-backend/util"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// configPath := "config.json"
	// config := util.ReadConfig(configPath)
	// fmt.Println(config.QueueName)
	// rabbitmq := rabbitmq.NewRabbitMQ()
	// rabbitmq.CreateChannel()
	// rabbitmq.QueueDeclare(config.QueueName)
	// defer rabbitmq.Close()
	router := gin.Default()
	router.POST("/upload", image.UploadFile)
	router.Run(":8080")
	config := util.ReadConfig()
	fmt.Println(config.QueueName)
	rabbitmq := rabbitmq.NewRabbitMQ()
	rabbitmq.CreateChannel()
	rabbitmq.QueueDeclare(config.QueueName)
	defer rabbitmq.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		message := time.Now().String()
		rabbitmq.Publish(message)
		fmt.Fprint(w, "Message published")
	})

}
