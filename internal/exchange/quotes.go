package exchange

type Quotes struct {
  Timestamp int             `json:"timestamp"`
  Qoutes map[string]float32 `json:"quotes"`
}
