package pubsub

import (
	"fmt"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
	ampq "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type Subscriber struct {
  Channel *ampq.Channel
  Queue *ampq.Queue
}

func NewSubscriber(url, queue string) *Subscriber {
  conn, _ := ampq.Dial(url)
  logging.LogInfo("Successfully connected to broker")

  ch, _ := conn.Channel()
  logging.LogInfo("Successfully opened a channel")

  q, _ := ch.QueueDeclare(
    "events", // name
    false,   // durable
    false,   // delete when unused
    false,   // exclusive
    false,   // no-wait
    nil,     // arguments
  )

  return &Subscriber{Channel: ch, Queue: &q}
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
    var event genproto.Event

    for d := range msgs {
      err := proto.Unmarshal(d.Body, &event)
      if err != nil {
        fmt.Println("Failed to parse message", err)
      } else {
        logging.LogEventMessage("Received event message:", &event)
      }
    }
  }()

  fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")
  <-forever
}
