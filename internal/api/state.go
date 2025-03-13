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

// ```json
// {
//   "events_total": 12345,
//   "events_per_minute": 123.45,
//   "events_per_second_moving_average": 3.12,
//   "top_player_bets": {
//     "id": 10,
//     "count": 150
//   },
//   "top_player_wins": {
//     "id": 11,
//     "count": 50
//   },
//   "top_player_deposits": {
//     "id": 12,
//     "count": 15000
//   }
// }
// ```

