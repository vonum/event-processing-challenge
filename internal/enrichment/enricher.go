package enrichment

import (
	"fmt"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/cache"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/db"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/exchange"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
)

const EURCurrency = "EUR"
const ApiTimeout = 5
const TickDurationMS = 200

type Enricher struct {
  exchangeClient *exchange.Client
  cacheClient *cache.Client
  dbClient *db.Client
}

func NewEnricher(
  apiKey string,
  redisAddr string,
  pgHost string,
  pgUser string,
  pgPassword string,
  pgPort int,
) *Enricher {
  exchangeClient := exchange.NewClient(apiKey, ApiTimeout)
  cacheClient := cache.NewClient(redisAddr)
  dbClient := db.NewClient(pgHost, pgUser, pgPassword, pgPort)

  return &Enricher{exchangeClient, cacheClient, dbClient}
}

func (e *Enricher) Enrich(event *casino.Event) error {
  if event.Currency == EURCurrency {
    logging.LogInfo(fmt.Sprintf("Event with ID %d is already in EUR.", event.ID))
    event.AmountEUR = event.Amount
  } else {
    logging.LogInfo(fmt.Sprintf("Getting EUR amount value for event %d.", event.ID))
    exchangeRate := e.getExchangeRate(event.Currency)
    logging.LogInfo(fmt.Sprintf("Exchange rate is %f.", exchangeRate))
    event.AmountEUR = event.Amount * int(exchangeRate)
  }

  logging.LogInfo(fmt.Sprintf("\nGetting details for user %d.", event.PlayerID))
  user, err := e.dbClient.GetUser(event.PlayerID)

  if err != nil {
    logging.LogInfo(fmt.Sprintf("\nDetails for player %d not found.", event.PlayerID))
  } else {
    logging.LogInfo(
      fmt.Sprintf(
        "\nPlayer %d found with email %s and last sign in at: %s.",
        event.PlayerID,
        user.Email,
        user.LastSignedInAt,
      ),
    )
    event.Player = casino.Player{Email: user.Email, LastSignedInAt: user.LastSignedInAt}
  }

  return nil
}

func (e *Enricher) getExchangeRate(currency string) float32 {
  ticker := time.NewTicker(TickDurationMS * time.Millisecond)

  for {
    select {
    case <- ticker.C:
      // try to read from cache and return value
      logging.LogInfo(fmt.Sprintf("Reading value from cache for %s.", currency))
      quotes, err := e.cacheClient.ReadQuotes()

      if err == nil {
        fmt.Println(quotes)
        exchangeRate := quotes.GetQuote(currency)
        logging.LogInfo(fmt.Sprintf("Found value in cache for %s: %f.", currency, exchangeRate))

        return exchangeRate
      }

      logging.LogInfo(fmt.Sprintf("Value not found in cache for %s.", currency))

      // if it's not in cache
      // try to acquire lock for consuming the api
      // and update the cache
      logging.LogInfo("Acquiring lock for consuming the api.")
      ackquiredLock, err := e.cacheClient.AcquireLock()

      if ackquiredLock {
        logging.LogInfo("Acquired lock for consuming the api.")
        quotes := e.exchangeClient.GetLatestExchangeRate(casino.Currencies)
        logging.LogInfo("Got currencies exchange rates from the api.")

        logging.LogInfo("Caching currencies")
        err := e.cacheClient.CacheQuotes(&quotes)
        if err != nil {
          fmt.Println(err)
        }

        e.cacheClient.ReleaseLock()
        logging.LogInfo("Released lock.")

        fmt.Println(quotes)
        exchangeRate := quotes.GetQuote(currency)
        return exchangeRate
      }

      logging.LogInfo(fmt.Sprintf("Failed to ackquire lock, retrying in %dms.", TickDurationMS))


      // keep trying
    }
  }
}
