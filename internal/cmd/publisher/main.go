package main

import (
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/generator"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pubsub"
	"golang.org/x/net/context"
)

func main() {
  config := pubsub.LoadConfig()


  publisher := pubsub.NewPublisher(config.RabbitMqAddr, pubsub.EventsQueue)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  eventCh := generator.Generate(ctx)

  for event := range eventCh {
    publisher.Send(&event)
    time.Sleep(300 * time.Millisecond)
  }

  logging.LogInfo("Finished sending messages.")
}
