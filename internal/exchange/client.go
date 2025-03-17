package exchange

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

const BaseUrl = "https://api.exchangerate.host/live"

type Client struct {
  BaseUrl string
  apiKey string
  HttpClient *http.Client
}

func NewClient(apiKey string, timeout int) *Client {
    return &Client{
        BaseUrl: BaseUrl,
        apiKey:  apiKey,
        HttpClient: &http.Client{
            Timeout: time.Duration(timeout) * time.Second,
        },
    }
}

func (c *Client) GetLatestExchangeRate(currencies []string) Quotes {
  curr := strings.Join(casino.Currencies, ",")

  url := fmt.Sprintf("%s?access_key=%s&source=EUR&currencies=%s", c.BaseUrl, c.apiKey, curr)
  req, _ := http.NewRequest("GET", url, nil)

  res, _ := c.HttpClient.Do(req)

  defer res.Body.Close()

  var quotes Quotes

  json.NewDecoder(res.Body).Decode(&quotes)

  return quotes
}
