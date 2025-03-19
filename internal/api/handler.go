package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

const WindowSize = 10

type ApiHandler struct {
  lock *sync.Mutex
  StartTime time.Time
  NEvents int
  UserBets map[int]int
  UserWins map[int]int
  UserDeposits map[int]int
  TPBets TopStats
  TPWins TopStats
  TPDeposits TopStats
  EventsPerSecond map[int64]int
}

func NewApiHandler() *ApiHandler {
  return &ApiHandler{
    lock: &sync.Mutex{},
    StartTime: time.Now(),
    NEvents: 0,
    UserBets: make(map[int]int),
    UserWins: make(map[int]int),
    UserDeposits: make(map[int]int),
    TPBets: TopStats{},
    TPWins: TopStats{},
    TPDeposits: TopStats{},
    EventsPerSecond: make(map[int64]int),
  }
}

func (h *ApiHandler) Materialize(w http.ResponseWriter, r *http.Request) {
  var window_size int
  param := r.URL.Query().Get("window_size")
  if param == "" {
    window_size = WindowSize
  } else {
    window_size, _ = strconv.Atoi(param)
  }

  stats := MaterializedStats{
    EventsTotal: h.NEvents,
    EventsPerMinute: h.eventsPerMinute(),
    EventsPerSecondMovingAverage: h.smaPerSecond(window_size),
    TopPlayerBets: h.TPBets,
    TopPlayerWins: h.TPWins,
    TopPlayerDeposits: h.TPDeposits,
  }

  json.NewEncoder(w).Encode(stats)
}

func (h *ApiHandler) PostEvent(w http.ResponseWriter, r *http.Request) {
  h.lock.Lock()
  defer h.lock.Unlock()

  decoder := json.NewDecoder(r.Body)
  var event casino.Event
  decoder.Decode(&event)

  h.NEvents++

  switch event.Type {
  case "bet":
    h.updateStats(event.PlayerID, 1, h.UserBets, &h.TPBets)
  case "deposit":
    h.updateStats(event.PlayerID, event.AmountEUR, h.UserDeposits, &h.TPDeposits)
  }

  if event.HasWon {
    h.updateStats(event.PlayerID, 1, h.UserWins, &h.TPWins)
  }

  second := time.Now().Unix()
  h.EventsPerSecond[second]++
}

func (h *ApiHandler) eventsPerMinute() float32 {
  t := time.Now()
  timedelta := t.Sub(h.StartTime).Minutes()

  return float32(h.NEvents) / float32(timedelta)
}

func (h *ApiHandler) smaPerSecond(seconds int) float32 {
  t := time.Now().Unix()
  s := 0
  for i := 0; i < seconds; i++ {
    s += h.EventsPerSecond[t - int64(i)]
  }

  if s == 0 {
    return 0
  }

  return float32(s) / float32(seconds)
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
