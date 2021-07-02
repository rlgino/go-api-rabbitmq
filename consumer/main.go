package main

import (
	"github.com/rlgino/rabbitmq/config"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main(){
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

	messages, err := channel.Consume(config.QUEUE, "", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
	}


	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever
}
