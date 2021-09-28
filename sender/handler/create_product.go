package handler

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

type RegisterProduct struct {
	channel *amqp.Channel
}

func NewRegisterProductHandler(connectRabbitMQ *amqp.Connection) RegisterProduct {
	channel, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	err = channel.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		panic(err)
	}
	return RegisterProduct{
		channel: channel,
	}
}

func (handler *RegisterProduct) Handle(ctx *fiber.Ctx) error {
	newProduct := struct {
		ID     int     `json:"id"`
		Name   string  `json:"name"`
		Price  float32 `json:"price"`
		UserID int     `json:"user_id"`
	}{}
	err2 := json.Unmarshal(ctx.Body(), &newProduct)
	if err2 != nil {
		return err2
	}
	// Save new product....
	log.Printf("Save product: [%d %s]\n", newProduct.ID, newProduct.Name)

	productCreated := struct {
		ProdID int `json:"prod_id"`
		UserID int `json:"user_id"`
	}{
		ProdID: newProduct.ID,
		UserID: newProduct.UserID,
	}
	event, _ := json.Marshal(productCreated)
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        event,
	}
	if err := handler.channel.Publish(
		"logs_topic",
		"rlgino.product_creator.1.event.product.created",
		false,
		false,
		message); err != nil {
		return err
	}

	return nil
}
