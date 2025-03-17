package pubsub

import (
	"context"
	"time"
	"google.golang.org/protobuf/proto"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
)

type Publisher struct {
  Channel *amqp.Channel
  Queue *amqp.Queue
}

func NewPublisher(url, queue string) *Publisher {
  conn, _ := amqp.Dial(url)
  logging.LogSetup("Successfully connected to broker.")

  ch, _ := conn.Channel()
  logging.LogSetup("Successfully opened a channel.")

  q := DeclareQueue(ch, queue)

  return &Publisher{Channel: ch, Queue: &q}
}

func (p *Publisher) Send(event *casino.Event) {
  ctx, cancel := context.WithTimeout(
    context.Background(),
    PublishTimeoutS * time.Second,
  )
  defer cancel()

  event_msg := genproto.Event{
    Id: int64(event.ID),
    PlayerId: int64(event.PlayerID),
    GameId: int64(event.GameID),
    Type: event.Type,
    Amount: int64(event.Amount),
    Currency: event.Currency,
    HasWon: event.HasWon,
    CreatedAt: int64(event.CreatedAt.Unix()),
  }
  logging.LogEventMessage("Sending event message:", &event_msg)

  body, _ := proto.Marshal(&event_msg)

  p.Channel.PublishWithContext(
    ctx,
    "",           // exchange
    p.Queue.Name, // routing key
    false,        // mandatory
    false,        // immediate
    amqp.Publishing {
      DeliveryMode: amqp.Persistent, // survives rabbitmq restart
      ContentType: "text/plain",
      Body:        body,
    },
  )

  logging.LogEventMessage("Successfully sent event message:", &event_msg)
  logging.LogInfo("\n")
}
