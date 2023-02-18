package rabbitmq

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue
}

type rabbitMQConfig struct {
	User     string `mapstructure:"RABBITMQ_DEFAULT_USER" validate:"required"`
	Password string `mapstructure:"RABBITMQ_DEFAULT_PASS" validate:"required"`
	Host     string `mapstructure:"RABBITMQ_HOST"`
	Port     string `mapstructure:"RABBITMQ_PORT"`
}

func (config *rabbitMQConfig) DSNFromConfig() string {

	dns := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
	return dns
}
func LoadDBConfig() (config rabbitMQConfig, err error) {

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.SetDefault("RABBITMQ_PORT", "5672")
	viper.SetDefault("RABBITMQ_HOST", "rabbitmq")
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

func NewRabbitMQ() *RabbitMQ {
	configDB, err := LoadDBConfig()
	if err != nil {
		log.Fatalf("Cannot read config for RabbitMQ: %v", err)
	}
	dns := configDB.DSNFromConfig()
	conn, err := amqp.Dial(dns)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ server: %v", err)
	}
	return &RabbitMQ{conn: conn, ch: nil, queue: nil}
}
func (r *RabbitMQ) CreateChannel() {
	ch, err := r.conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	r.ch = ch
}

func (r *RabbitMQ) QueueDeclare(queueName string) {
	queue, err := r.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	r.queue = &queue
}

func (r *RabbitMQ) Publish(message string) error {
	err := r.ch.Publish(
		"",           // exchange
		r.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	return err
}
func (r *RabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}
