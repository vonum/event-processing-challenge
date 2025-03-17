package main

import (
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pubsub"
)

func main() {
  const queue = "events"
  config := pubsub.LoadConfig()

  subscriber := pubsub.NewSubscriber(config, queue)
  subscriber.Read()
}
