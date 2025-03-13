package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
  Channel *amqp.Channel
  Queue *amqp.Queue
}

func NewPublisher(url, queue string) *Publisher {
  conn, _ := amqp.Dial(url)
  logging.LogInfo("Successfully connected to broker")

  ch, _ := conn.Channel()
  logging.LogInfo("Successfully opened a channel")

  q, _ := ch.QueueDeclare(
    queue,   // queue name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )

  return &Publisher{Channel: ch, Queue: &q}
}

func (p *Publisher) Send(event *casino.Event) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  event_msg := genproto.Event{
    Id: int32(event.ID),
    PlayerId: int32(event.PlayerID),
    Type: event.Type,
    Amount: int32(event.Amount),
    Currency: event.Currency,
    HasWon: event.HasWon,
    CreatedAt: int32(event.CreatedAt.Unix()),
  }
  logging.LogEventMessage("Sending event message", &event_msg)

  body, _ := proto.Marshal(&event_msg)
  fmt.Println("BODY")
  fmt.Println(body)

  p.Channel.PublishWithContext(ctx,
    "",           // exchange
    p.Queue.Name, // routing key
    false,        // mandatory
    false,        // immediate
    amqp.Publishing {
      ContentType: "text/plain",
      Body:        body,
    },
  )

  logging.LogEventMessage("Successfully sent event message", &event_msg)
}
