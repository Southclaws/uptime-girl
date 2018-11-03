package uptimerobot

import (
	"gopkg.in/resty.v1"
)

// Client represents an UptimeRobot HTTP client.
type Client struct {
	r   *resty.Client
	key string
}

// New creates a new UptimeRobot client with the given API key.
func New(key string) (client *Client, err error) {
	client = &Client{
		r: resty.New().
			SetHostURL("https://api.uptimerobot.com/v2/").
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("cache-control", "no-cache"),
		key: key,
	}

	return
}

// Request returns a pre-made request with an API key and the given body string.
func (c *Client) Request(endpoint string, body string) (response *resty.Response, err error) {
	// marshall body as url.Values
	// inject api_key field
	return c.r.R().SetBody(body).Post(endpoint)
}
