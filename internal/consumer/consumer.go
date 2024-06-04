package main

import (
  "log"

  "github.com/streadway/amqp"
)

func main() {
  conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
  if err != nil {
    log.Fatalf("Failed to connect to RabbitMQ: %s", err)
  }
  defer conn.Close()

  channel, err := conn.Channel()
  if err != nil {
    log.Fatalf("Failed to open a channel: %s", err)
  }
  defer channel.Close()

  msgs, err := channel.Consume(
    "api_requests",
    "",
    true,
    false,
    false,
    false,
    nil,
  )
  if err != nil {
    log.Fatalf("Failed to register a consumer: %s", err)
  }

  forever := make(chan bool)

  go func() {
    for d := range msgs {
      log.Printf("Received a message: %s", d.Body)
      // Process the message here
    }
  }()

  log.Printf("Waiting for messages. To exit press CTRL+C")
  <-forever
}