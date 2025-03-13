package logging

import (
	"github.com/Bitstarz-eng/event-processing-challenge/internal/genproto"
	"github.com/fatih/color"
)

func LogInfo(msg string) {
  color.Cyan(msg)
}

func LogEvent() {

}

func LogEventMessage(msg string, event *genproto.Event) {
  color.Cyan(msg, event.String())
}

func LogError(msg string) {
  color.Red(msg)
}
