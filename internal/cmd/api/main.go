package main

import (
	"net/http"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/api"
)

func main() {
  const port = ":3000"
  handler := api.NewApiHandler()

  http.HandleFunc("/health", handler.Health)
  http.HandleFunc("/materialize", handler.Materialize)
  http.HandleFunc("/events", handler.PostEvent)

  server := &http.Server{Addr: port}

  server.ListenAndServe()
}
