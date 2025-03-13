package api

import (
	"encoding/json"
	"net/http"
)

type ApiHandler struct {}

func (h *ApiHandler) Materialize(w http.ResponseWriter, r *http.Request) {
  stats := MaterializedStats{
    EventsTotal: 10,
    EventsPerMinute: 3.2,
    EventsPerSecondMovingAverage: 1.1,
    TopPlayerBets: TopStats{ID: 1, Count: 10},
    TopPlayerWins: TopStats{ID: 1, Count: 10},
    TopPlayerDeposits: TopStats{ID: 1, Count: 10},
  }

  json.NewEncoder(w).Encode(stats)
}


func (h *ApiHandler) Health(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}
