package exchange

import "fmt"

type Quotes struct {
  Timestamp int             `json:"timestamp"`
  Quotes map[string]float32 `json:"quotes"`
}

func (q *Quotes) GetQuote(currency string) float32 {
  curr := fmt.Sprintf("EUR%s", currency)
  return q.Quotes[curr]
}
