package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel    *amqp.Channel
	routingKey string
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	defer conn.Close()

	// declare channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Cannot create publisher: %s", err.Error())
	}

	defer ch.Close()

	if err != nil {
		log.Fatal("Unable to declare channel")
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Fatal("Unable to declare queue")
		return nil, err
	}

	return &Publisher{channel: ch, routingKey: q.Name}, nil
}

func (p *Publisher) Publish(body string) {
	err := p.channel.Publish(
		"",           // exchange
		p.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		log.Println(err)
		log.Fatal("Failed to publish a message")
	}
}
