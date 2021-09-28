package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func main() {
	amqpServerUrl := os.Getenv("AMQP_SERVER_URL")
	amqpServerUrl = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(amqpServerUrl)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare("logs_topic", "topic", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
	}

	q, err := channel.QueueDeclare("product_notifier.user.publish_on_product_created", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = channel.QueueBind(
		q.Name,                   // queue name
		"rlgino.product_creator.1.event.product.created", // routing key
		"logs_topic",             // exchange
		false,
		nil)
	if err != nil {
		log.Println(err.Error())
	}

	messages, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	forever := make(chan bool)

	go func() {
		for message := range messages {

			productCreated := struct {
				ProdID int `json:"prod_id"`
				UserID int `json:"user_id"`
			}{}
			json.Unmarshal(message.Body, &productCreated)
			// For example, show received message in a console.
			log.Printf(" > Notifier that user %d create a product %d\n", productCreated.UserID, productCreated.ProdID)
		}
	}()

	<-forever
}
