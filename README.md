# API Golang con RabbitMQ

## Endpoints

### Product
POST
```bash
curl --location --request GET 'localhost:3000/send' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 1,
    "name": "Teclado",
    "price": 200,
    "user_id": 5
}'
```

## Queues
* com.rlgino.go.product_created

## Sources
* [📈 Working with RabbitMQ in Golang by examples](https://dev.to/koddr/working-with-rabbitmq-in-golang-by-examples-2dcn)