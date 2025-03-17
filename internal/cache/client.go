package cache

import (
	ctx "context"
	"encoding/json"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/exchange"
	"github.com/redis/go-redis/v9"
)

const LockKey = "quotes:lock"
const QoutesKey = "qoutes:data"
const LockTimeoutS = 5
const QoutesTimeoutS = 60

// Caching strategy
// If data is not present, a lock is acquired to prevent multiple workers
// From spamming the external api
// If a worker manages to acquire a lock, it will update the cache and release it
// If a worker doesn't manage to acquire a lock, it will poll the cache
// All quotes are updated at the same time for simplicity

// Data is stored as json since it was decided to store all quotes together for simplicity

type Client struct {
  client *redis.Client
}

func NewClient(address string) *Client {
  client := redis.NewClient(&redis.Options{
    Addr: address,
    Password: "",
    DB: 0,
  })

  return &Client{client}
}

func (c *Client) CacheQuotes(qoutes *exchange.Quotes) error {
  qoutesJson := marshalQuotes(*qoutes)
  err := c.client.Set(
    ctx.Background(),
    QoutesKey,
    qoutesJson,
    QoutesTimeoutS * time.Second,
  ).Err()

  if err != nil {
    panic(err)
  }

  return err
}

func (c *Client) ReadQuotes() (*exchange.Quotes, error) {
  qoutesJson, err := c.client.Get(ctx.Background(), QoutesKey).Result()
  if err != nil {
    return nil, err
  }
  return unmarshalQuotes([]byte(qoutesJson)), nil
}

func (c *Client) AcquireLock() (bool, error) {
  acquired, err := c.client.SetNX(
    ctx.Background(),
    LockKey,
    "1",
    LockTimeoutS * time.Second,
  ).Result()

  return acquired, err
}

func (c *Client) ReleaseLock() {
  c.client.Del(ctx.Background(), LockKey)
}

func marshalQuotes(qoutes exchange.Quotes) string {
  qoutesJson, _ := json.Marshal(qoutes)
  return string(qoutesJson)
}

func unmarshalQuotes(qoutesJson []byte) *exchange.Quotes {
  var qoutes exchange.Quotes
  json.Unmarshal(qoutesJson, &qoutes)
  return &qoutes
}
