package db

import "time"

type User struct {
  Email string
  LastSignedInAt time.Time
}
