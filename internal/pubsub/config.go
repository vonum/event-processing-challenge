package pubsub

import "os"

const EventsQueue = "events"
const DeadLetterQueue = "dead_letter"

type Config struct {
  ApiKey string
  RabbitMqAddr string
  RedisAddr string
  PgConn string
}

func LoadConfig() *Config {
  apiKey := getEnv("API_KEY", "")
  rabbitMqAddr := getEnv(
    "RABBITMQ_ADDRESS",
    "amqp://guest:guest@localhost:5672/",
  )
  redisAddr := getEnv("REDIS_ADDRESS", "localhost:6379")
  pgConn := getEnv(
    "POSTGRES_CONNECTION_STRING",
    "host=127.0.0.1 port=5432 user=casino password=casino sslmode=disable",
  )

  return &Config{apiKey, rabbitMqAddr, redisAddr, pgConn}
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
