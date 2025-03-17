package pubsub

import (
  ampq "github.com/rabbitmq/amqp091-go"
)

const PublishTimeoutS = 5

func DeclareQueue(channel *ampq.Channel, name string) ampq.Queue {
  queue, _ := channel.QueueDeclare(
    name,
    true,  // durable
    false, // delete when unused
    false, // exclusive
    false, // no-wait
    nil,   // arguments
  )

  return queue
}
