package db

import (
	"database/sql"
	"fmt"
  _ "github.com/lib/pq"
)

type Client struct {
  db *sql.DB
}

func NewClient(host, user, password string, port int) *Client {
  dataSource := fmt.Sprintf(
    "host=%s port=%d user=%s password=%s sslmode=disable",
    host,
    port,
    user,
    password,
  )

  db, _ := sql.Open("postgres", dataSource)
  return &Client{db}
}

func (c *Client) GetUser(userId int) (*User, error) {
  var user User

  query := `
    SELECT email, last_signed_in_at FROM players
    WHERE id = $1
  `

  row := c.db.QueryRow(query, userId)
  err := row.Scan(&user.Email, &user.LastSignedInAt)

  if err != nil {
    return nil, err
  }

  return &user, nil
}

func (c *Client) Close() {
  c.db.Close()
}
