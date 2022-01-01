package rabbitmq

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"log"
	"neotype-backend/pkg/config"
	"os"
)

const (
	cloudamqpEnv     = "CLOUDAMQP_URL"
	localUserEnv     = "RABBITMQ_USER"
	localPasswordEnv = "RABBITMQ_PASSWORD"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func (r *RabbitMQ) Connect(queueName string, durable, autoDelete, exclusive, noWait bool, args map[string]interface{}) error {
	var err error

	_ = godotenv.Load("docker/.env")

	cloudamqp, exists := os.LookupEnv(cloudamqpEnv)
	if exists {
		r.Conn, err = amqp.Dial(cloudamqp)
	} else {
		login := os.Getenv(localUserEnv)
		password := os.Getenv(localPasswordEnv)
		if login == "" || password == "" {
			log.Fatal("RabbitMQ environment variables not set")
		}

		addr, err := config.Get("rabbitmq", "addr")
		port, err := config.Get("rabbitmq", "port")
		if err != nil {
			log.Fatal("Failed to get config for RabbitMQ")
		}

		r.Conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", login, password, addr, port))
	}
	if err != nil {
		return err
	}

	r.Channel, err = r.Conn.Channel()
	if err != nil {
		return err
	}

	r.Queue, err = r.Channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, args)
	return err
}

func New() *RabbitMQ {
	return &RabbitMQ{}
}
