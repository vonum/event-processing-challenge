package pubsub

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/enrichment"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
)

type Subscriber struct {
  Channel *amqp.Channel
  EventsQueue *amqp.Queue
  DeadLetterQueue *amqp.Queue
  Enricher *enrichment.Enricher
  ApiUrl string
}

func NewSubscriber(
  config *Config,
  eventsQueue string,
  deadLetterQueue string,
) *Subscriber {
  conn, _ := amqp.Dial(config.RabbitMqAddr)
  logging.LogSetup("Successfully connected to broker")

  ch, _ := conn.Channel()
  logging.LogSetup("Successfully opened a channel")

  q := DeclareQueue(ch, eventsQueue)
  dlq := DeclareQueue(ch, deadLetterQueue)

  e := enrichment.NewEnricher(
    config.ApiKey,
    config.RedisAddr,
    config.PgConn,
  )

  return &Subscriber{
    Channel: ch,
    EventsQueue: &q,
    DeadLetterQueue: &dlq,
    Enricher: e,
    ApiUrl: config.ApiUrl,
  }
}

func (s *Subscriber) Read() {
  msgs, _ := s.Channel.Consume(
    s.EventsQueue.Name, // queue
    "",                 // consumer
    false,              // auto-ack
    false,              // exclusive
    false,              // no-local
    false,              // no-wait
    nil,                // args
  )

  forever := make(chan struct{})
  // var forever chan struct{}

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

          logging.LogInfo("Publishing event to dlq.")
          s.PublishToDLQ(d.Body)
          d.Ack(false)

          continue
        }

        // acknowledge that messages are handled
        // once messages are acked, they are deleted
        // by default, message will be routed to a different consumer
        // if the connection is lost or if timeout is exceeded
        logging.LogEventPretty(event)
        logging.LogEvent(event)

        logging.LogInfo("Posting event to api.")
        s.PostEvent(&event)

        d.Ack(false)
      }
    }
  }()

  logging.LogSetup(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}

func (s *Subscriber) PublishToDLQ(msg []byte) {
  ctx, cancel := context.WithTimeout(
    context.Background(),
    PublishTimeoutS * time.Second,
  )
  defer cancel()

  s.Channel.PublishWithContext(
    ctx,
    "",                     // exchange
    s.DeadLetterQueue.Name, // routing key
    false,                  // mandatory
    false,                  // immediate
    amqp.Publishing {
      DeliveryMode: amqp.Persistent, // survives rabbitmq restart
      ContentType: "text/plain",
      Body:        msg,
    },
  )
}

func (s *Subscriber) PostEvent(event *casino.Event) {
  // error handling
  apiUrl, _ := url.JoinPath(s.ApiUrl, "events")
  logging.LogInfo(apiUrl)

  jsonValue, _ := json.Marshal(event)
  http.Post(apiUrl, "application/json", bytes.NewBuffer(jsonValue))
}
