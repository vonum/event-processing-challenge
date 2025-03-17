package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type ApiHandler struct {
  StartTime time.Time
  NEvents int
  UserBets map[int]int
  UserWins map[int]int
  UserDeposits map[int]int
  TPBets TopStats
  TPWins TopStats
  TPDeposits TopStats
}

func NewApiHandler() *ApiHandler {
  return &ApiHandler{
    StartTime: time.Now(),
    NEvents: 0,
    UserBets: make(map[int]int),
    UserWins: make(map[int]int),
    UserDeposits: make(map[int]int),
    TPBets: TopStats{},
    TPWins: TopStats{},
    TPDeposits: TopStats{},
  }
}

func (h *ApiHandler) Materialize(w http.ResponseWriter, r *http.Request) {
  stats := MaterializedStats{
    EventsTotal: h.NEvents,
    EventsPerMinute: h.eventsPerMinute(),
    EventsPerSecondMovingAverage: 1.1,
    TopPlayerBets: h.TPBets,
    TopPlayerWins: h.TPWins,
    TopPlayerDeposits: h.TPDeposits,
  }

  json.NewEncoder(w).Encode(stats)
}

func (h *ApiHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Received event")
  decoder := json.NewDecoder(r.Body)
  var event casino.Event
  decoder.Decode(&event)
  fmt.Println(event)

  h.NEvents++

  switch event.Type {
  case "bet":
    h.updateStats(event.PlayerID, 1, h.UserBets, &h.TPBets)
    // bets := h.UserBets[event.PlayerID] + 1
    // h.UserBets[event.PlayerID] = bets
    // if bets > h.TPBets.Count {
    //   h.TPBets.ID = event.PlayerID
    //   h.TPBets.Count = bets
    // }
  case "deposit":
    h.updateStats(event.PlayerID, event.AmountEUR, h.UserDeposits, &h.TPDeposits)
    // deposits := h.UserDeposits[event.PlayerID] + event.AmountEUR
    // h.UserDeposits[event.PlayerID] = deposits
    // if deposits > h.TPDeposits.Count {
    //   h.TPDeposits.ID = event.PlayerID
    //   h.TPDeposits.Count = deposits
    // }
  }

  if event.HasWon {
    h.updateStats(event.PlayerID, 1, h.UserWins, &h.TPWins)
    // wins := h.UserWins[event.PlayerID] + 1
    // h.UserWins[event.PlayerID] = wins
    // if wins > h.TPWins.Count {
    //   h.TPWins.ID = event.PlayerID
    //   h.TPWins.Count = wins
    // }
  }

}

func (h *ApiHandler) eventsPerMinute() float32 {
  t := time.Now()
  timedelta := t.Sub(h.StartTime).Minutes()

  return float32(h.NEvents) / float32(timedelta)
}

func (h *ApiHandler) updateStats(
  playerId int,
  value int,
  data map[int]int,
  stats *TopStats,
) {
  v := data[playerId] + value
  data[playerId] = v
  if v > stats.Count {
    stats.ID = playerId
    stats.Count = v
  }
}

func (h *ApiHandler) Health(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}
