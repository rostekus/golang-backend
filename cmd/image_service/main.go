package main

import (
	"fmt"
	"log"
	"net/http"
	"rostekus/golang-backend/rabbitmq"
	"rostekus/golang-backend/util"
	"time"
)

func main() {
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

	log.Println("Listening on :8080...")
	if err := http.ListenAndServe(":23451", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
