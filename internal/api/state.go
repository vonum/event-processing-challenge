package api

type TopStats struct {
  ID int
  Count int
}

type MaterializedStats struct {
  EventsTotal                   int `json:"events_total"`
  EventsPerMinute               float32 `json:"events_per_minute"`
  EventsPerSecondMovingAverage  float32 `json:"events_per_second_moving_average"`

  TopPlayerBets     TopStats `json:"top_player_bets"`
  TopPlayerWins     TopStats `json:"top_player_wins"`
  TopPlayerDeposits TopStats `json:"top_player_deposits"`
}
