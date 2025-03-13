package main

import (
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pubsub"
)

func main() {
  const url = "amqp://guest:guest@localhost:5672/"
  const queue = "events"

  subscriber := pubsub.NewSubscriber(url, queue)
  subscriber.Read()
}
