package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"

	"github.com/rlgino/rabbitmq/sender/handler"
)

func main() {
	amqpServerUrl := os.Getenv("AMQP_SERVER_URL")
	amqpServerUrl = "amqp://guest:guest@localhost:5672/"

	connectRabbitMQ, err := amqp.Dial(amqpServerUrl)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()


	app := fiber.New()
	app.Use(logger.New())

	productHandler := handler.NewRegisterProductHandler(connectRabbitMQ)
	app.Get("/send", productHandler.Handle)

	log.Fatal(app.Listen(":3000"))
}
