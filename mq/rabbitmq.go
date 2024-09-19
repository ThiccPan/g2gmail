package mq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitQueueDeclare(ch *amqp.Channel) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		"hello", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	return q, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
