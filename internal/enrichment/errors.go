package enrichment

import "fmt"

type FetchingQuotesTimeoutError struct {
  Timeout int
}

func (e *FetchingQuotesTimeoutError) Error() string {
  return fmt.Sprintf("Failed to fetch quotes after %ds.", e.Timeout)
}
