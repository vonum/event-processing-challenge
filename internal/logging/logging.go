package logging

import (
	"fmt"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
)

func LogInfo(msg string) {
  color.Cyan(msg)
}

func LogEvent() {
}

func LogEventPretty(event casino.Event) {
  var msg string
  switch event.Type {
  case "game_start":
    title := casino.Games[event.GameID].Title
    msg = fmt.Sprintf(
      "Player #%d started playing a game \"%s\" on %s.",
      event.PlayerID,
      title,
      formatTime(event.CreatedAt),
    )
  case "bet":
    title := casino.Games[event.GameID].Title
    msg = fmt.Sprintf(
      "Player #%s placed a bet of %d%s (%d EUR) on a game \"%s\" on %s.",
      event.Player.Email,
      event.Amount,
      event.Currency,
      event.AmountEUR,
      title,
      formatTime(event.CreatedAt),
    )
  case "deposit":
    msg = fmt.Sprintf(
      "Player #%d made a deposit of %d EUR on %s.",
      event.PlayerID,
      event.AmountEUR,
      formatTime(event.CreatedAt),
    )
  case "game_stop":
    title := casino.Games[event.GameID].Title
    msg = fmt.Sprintf(
      "Player #%d stopped playing a game \"%s\" on %s.",
      event.PlayerID,
      title,
      formatTime(event.CreatedAt),
    )
  default:
    msg = fmt.Sprintf("Unknown event %s", event.Type)
  }

  color.Green(msg)
}

func LogEventMessage(msg string, event *genproto.Event) {
  color.Cyan(msg, event.String(), "\n")
}

func LogError(msg string) {
  color.Red(msg)
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

// Player #10 started playing a game "Rocket Dice" on January 10th, 2022 at 12:34 UTC.
// Player #11 (john@example.com) placed a bet of 5 USD (4.68 EUR) on a game "It's bananas!" on February 2nd, 2022 at 23:45 UTC.
// Player #12 made a deposit of 100 EUR on February 3rd, 2022 at 12:12 UTC.
