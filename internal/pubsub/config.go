package pubsub

import "os"

const EventsQueue = "events"
const DeadLetterQueue = "dead_letter"

type Config struct {
  ApiKey string
  RabbitMqAddr string
  RedisAddr string
  PgConn string
  ApiUrl string
}

func LoadConfig() *Config {
  apiKey := getEnv("API_KEY", "1b894e89bd173b9bc1e5e3d55bb85c04")
  rabbitMqAddr := getEnv(
    "RABBITMQ_ADDRESS",
    "amqp://guest:guest@localhost:5672/",
  )
  redisAddr := getEnv("REDIS_ADDRESS", "localhost:6379")
  pgConn := getEnv(
    "POSTGRES_CONNECTION_STRING",
    "host=127.0.0.1 port=5432 user=casino password=casino sslmode=disable",
  )

  apiUrl := getEnv("API_URL", "http://localhost:3000")

  return &Config{apiKey, rabbitMqAddr, redisAddr, pgConn, apiUrl}
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
