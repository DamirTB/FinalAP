package sender

import (
    "github.com/streadway/amqp"
    "errors"
)

var channel *amqp.Channel

// InitRabbitMQ initializes the RabbitMQ channel
func InitRabbitMQ() error {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    channel = ch

    return nil
}

// PublishMessage publishes a message to RabbitMQ
func PublishMessage(message string) error {
    if channel == nil {
        return errors.New("RabbitMQ channel not initialized")
    }

    err := channel.Publish(
        "",
        "api_requests",
        false,
        false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(message),
        },
    )
    if err != nil {
        return err
    }
    return nil
}

// CloseRabbitMQ closes the RabbitMQ channel
func CloseRabbitMQ() {
    if channel != nil {
        channel.Close()
    }
}
