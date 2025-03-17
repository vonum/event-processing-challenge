package pubsub

import (
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/enrichment"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
	ampq "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Subscriber struct {
  Channel *ampq.Channel
  Queue *ampq.Queue
  Enricher *enrichment.Enricher
}

func NewSubscriber(
  url string,
  queue string,
) *Subscriber {
  conn, _ := ampq.Dial(url)
  logging.LogSetup("Successfully connected to broker")

  ch, _ := conn.Channel()
  logging.LogSetup("Successfully opened a channel")

  q, _ := ch.QueueDeclare(
    "events", // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )

  e := enrichment.NewEnricher(
    "1b894e89bd173b9bc1e5e3d55bb85c04",
    "localhost:6379",
    "127.0.0.1",
    "casino",
    "casino",
    5432,
  )

  return &Subscriber{Channel: ch, Queue: &q, Enricher: e}
}

func (s *Subscriber) Read() {
  msgs, _ := s.Channel.Consume(
    s.Queue.Name, // queue
    "",           // consumer
    true,         // auto-ack
    false,        // exclusive
    false,        // no-local
    false,        // no-wait
    nil,          // args
  )

  var forever chan struct{}

  go func() {
    var eventMsg genproto.Event

    for d := range msgs {
      err := proto.Unmarshal(d.Body, &eventMsg)
      if err != nil {
        logging.LogError("Failed to parse event message.")
        logging.LogError(err.Error())
      } else {
        logging.LogEventMessage("\nReceived event message:", &eventMsg)
        event := casino.Event{
          ID: int(eventMsg.Id),
          PlayerID: int(eventMsg.PlayerId),
          GameID: int(eventMsg.GameId),
          Type: eventMsg.Type,
          Amount: int(eventMsg.Amount),
          Currency: eventMsg.Currency,
          HasWon: eventMsg.HasWon,
          CreatedAt: time.Unix(eventMsg.CreatedAt, 0),
        }

        err := s.Enricher.Enrich(&event)
        if err != nil {
          logging.LogError("Failed to enrich event.")
          logging.LogError(err.Error())
          continue
        }

        logging.LogEventPretty(event)
      }
    }
  }()

  logging.LogSetup(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}
