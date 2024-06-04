package rabbitmq

import (
  "github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func InitRabbitMQ() error {
  var err error
  conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
  if err != nil {
    return err
  }

  channel, err = conn.Channel()
  if err != nil {
    return err
  }

  _, err = channel.QueueDeclare(
    "api_requests",
    true,
    false,
    false,
    false,
    nil,
  )
  if err != nil {
    return err
  }
  return nil
}

func PublishMessage(message string) error {
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

func CloseRabbitMQ() {
  if conn != nil {
    conn.Close()
  }
  if channel != nil {
    channel.Close()
  }
}