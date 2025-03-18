package enrichment

import (
	"fmt"
	"time"
	"github.com/dustin/go-humanize"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/cache"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/db"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/exchange"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/logging"
)

const EURCurrency = "EUR"
const ApiTimeoutS = 3
const TickDurationMS = 200
const QuoteTimeoutS = 5

type Enricher struct {
  exchangeClient *exchange.Client
  cacheClient *cache.Client
  dbClient *db.Client
}

func NewEnricher(
  apiKey string,
  redisAddr string,
  pgConn string,
) *Enricher {
  exchangeClient := exchange.NewClient(apiKey, ApiTimeoutS)
  cacheClient := cache.NewClient(redisAddr)
  dbClient := db.NewClient(pgConn)

  return &Enricher{exchangeClient, cacheClient, dbClient}
}

func (e *Enricher) Enrich(event *casino.Event) error {
  if event.Currency == EURCurrency {
    logging.LogInfo(fmt.Sprintf("Event with ID %d is already in EUR.", event.ID))
    event.AmountEUR = event.Amount
  } else {
    logging.LogInfo(fmt.Sprintf("Getting EUR amount value for event %d.", event.ID))

    exchangeRate, err := e.getExchangeRate(event.Currency)
    if err != nil {
      return err
    }

    event.AmountEUR = int(float32(event.Amount) / exchangeRate)
  }

  logging.LogInfo(fmt.Sprintf("Getting details for user %d.\n", event.PlayerID))
  user, err := e.dbClient.GetUser(event.PlayerID)

  if err != nil {
    logging.LogInfo(fmt.Sprintf("Details for player %d not found.\n", event.PlayerID))
  } else {
    logging.LogInfo(
      fmt.Sprintf(
        "Player %d found with email %s and last sign in at: %s.\n",
        event.PlayerID,
        user.Email,
        user.LastSignedInAt,
      ),
    )
    event.Player = casino.Player{Email: user.Email, LastSignedInAt: user.LastSignedInAt}
  }

  event.Description = e.formatDescription(event)

  return nil
}

func (e *Enricher) getExchangeRate(currency string) (float32, error) {
  ticker := time.NewTicker(TickDurationMS * time.Millisecond)
  timeout := time.After(QuoteTimeoutS * time.Second)

  for {
    select {
    case <- timeout:
      return 0, &FetchingQuotesTimeoutError{QuoteTimeoutS}
    case <- ticker.C:
      // try to read from cache and return value
      logging.LogInfo(fmt.Sprintf("Reading value from cache for %s.", currency))
      quotes, err := e.cacheClient.ReadQuotes()

      if err == nil {
        exchangeRate := quotes.GetQuote(currency)
        logging.LogInfo(fmt.Sprintf("Found value in cache for %s: %f.", currency, exchangeRate))

        return exchangeRate, nil
      }

      logging.LogInfo(fmt.Sprintf("Value not found in cache for %s.", currency))

      // if it's not in cache
      // try to acquire lock for consuming the api
      // and update the cache
      logging.LogInfo("Acquiring lock for consuming the api.")
      ackquiredLock, err := e.cacheClient.AcquireLock()

      if ackquiredLock {
        logging.LogInfo("Getting exchange rates from the api.")
        quotes, err := e.exchangeClient.GetLatestExchangeRate(casino.Currencies)
        if err != nil {
          return 0, err
        }

        logging.LogInfo("Caching currencies.")
        err = e.cacheClient.CacheQuotes(quotes)

        if err != nil {
          logging.LogError("Failed to update cache.")
        }

        e.cacheClient.ReleaseLock()
        logging.LogInfo("Released lock.")

        exchangeRate := quotes.GetQuote(currency)
        return exchangeRate, nil
      }

      logging.LogInfo(fmt.Sprintf("Failed to ackquire lock, retrying in %dms.", TickDurationMS))


      // keep trying
    }
  }
}

func (e *Enricher) formatDescription(event *casino.Event) string {
  var dsc string
  switch event.Type {
  case "game_start":
    title := casino.Games[event.GameID].Title
    dsc = fmt.Sprintf(
      "Player #%d started playing a game \"%s\" on %s.",
      event.PlayerID,
      title,
      formatTime(event.CreatedAt),
    )
  case "bet":
    title := casino.Games[event.GameID].Title
    dsc = fmt.Sprintf(
      "Player #%s placed a bet of %d%s (%d EUR) on a game \"%s\" on %s.",
      event.Player.Email,
      event.Amount,
      event.Currency,
      event.AmountEUR,
      title,
      formatTime(event.CreatedAt),
    )
  case "deposit":
    dsc = fmt.Sprintf(
      "Player #%d made a deposit of %d EUR on %s.",
      event.PlayerID,
      event.AmountEUR,
      formatTime(event.CreatedAt),
    )
  case "game_stop":
    title := casino.Games[event.GameID].Title
    dsc = fmt.Sprintf(
      "Player #%d stopped playing a game \"%s\" on %s.",
      event.PlayerID,
      title,
      formatTime(event.CreatedAt),
    )
  default:
    dsc = fmt.Sprintf("Unknown event %s", event.Type)
  }

  return dsc
}

func formatTime(t time.Time) string {
  return fmt.Sprintf(
    "%s %s, %d at %02d:%02d UTC",
    t.Format("January"),
    humanize.Ordinal(t.Day()),
    t.Year(),
    t.Hour(),
    t.Minute(),
  )
}
