package db

import "time"

type User struct {
  email string
  last_signed_in_at time.Time
}
