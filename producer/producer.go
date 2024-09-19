package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/thiccpan/g2gmail/model"
	"github.com/thiccpan/g2gmail/mq"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := mq.InitQueueDeclare(ch)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := model.Task{
		Value: bodyFrom(os.Args),
	}

	data, _ := json.Marshal(body)
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        data,
		})

	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func bodyFrom(args []string) string {
	var s string
	log.Println(args[1:])
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
