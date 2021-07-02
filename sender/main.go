package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	amqpServerUrl := os.Getenv("AMQP_SERVER_URL")
	amqpServerUrl = "amqp://guest:guest@localhost:5672/"

	connectRabbitMQ, err := amqp.Dial(amqpServerUrl)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channel, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	_, err = channel.QueueDeclare("QueueService1", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(logger.New())

	app.Get("/send", func(ctx *fiber.Ctx) error {
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(ctx.Query("msg")),
		}
		if err := channel.Publish("", "QueueService1", false, false, message); err != nil {
			return err
		}

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
