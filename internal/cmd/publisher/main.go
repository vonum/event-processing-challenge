package main

import (
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/generator"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pubsub"
	"golang.org/x/net/context"
)

func main() {
  const url = "amqp://guest:guest@localhost:5672/"
  const queue = "events"

  publisher := pubsub.NewPublisher(url, queue)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  eventCh := generator.Generate(ctx)

  for event := range eventCh {
    // log.Printf("%#v\n", event)
    publisher.Send(&event)
    time.Sleep(3 * time.Second)
  }

  logging.LogInfo("Finished sending messages.")
}
