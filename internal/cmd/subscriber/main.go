package main

import (
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pubsub"
)

func main() {
  config := pubsub.LoadConfig()

  subscriber := pubsub.NewSubscriber(config, pubsub.EventsQueue, pubsub.DeadLetterQueue)
  subscriber.Read()
}
