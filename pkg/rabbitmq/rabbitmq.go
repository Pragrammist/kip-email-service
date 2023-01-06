package rabbitmq

import (
	"fmt"
	"go-email/config"
	"log"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg *config.Config) *amqp.Connection {
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Rabbit.User,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		strconv.Itoa(cfg.Rabbit.Port),
	)

	conn, err := amqp.Dial(connAddr)

	if err != nil {
		log.Println(err)
		log.Fatal("Cannot connect to rabbitmq")
	}

	return conn
}
