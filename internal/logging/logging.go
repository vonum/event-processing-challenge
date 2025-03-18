package logging

import (
	"encoding/json"
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
)

func init() {
  log.SetFormatter(&log.JSONFormatter{})
  log.SetOutput(os.Stdout)
}

func LogInfo(msg string) {
  color.Cyan(msg)
}

func LogSetup(msg string) {
  color.Yellow(msg)
}

func LogEvent(event casino.Event) {
  eventLog, err := json.Marshal(event)
  if err != nil {
    log.WithError(err).Error("Failed to marshal event")
    return
  }

  c := color.New(color.FgGreen).Sprint
  log.WithFields(log.Fields{
    "event": string(eventLog),
  }).Info(c("Finished parsing event."))
}

func LogEventPretty(event casino.Event) {
  color.Green(event.Description)
}

func LogEventMessage(msg string, event *genproto.Event) {
  color.HiBlue(msg, event.String())
  color.HiBlue("\n")
}

func LogError(msg string) {
  color.Red(msg)
}
