package uptimerobot

import (
	"net/url"

	"gopkg.in/resty.v1"
)

// Client represents an UptimeRobot HTTP client.
type Client struct {
	r   *resty.Client
	key string
}

// New creates a new UptimeRobot client with the given API key.
func New(key string) (client *Client) {
	return &Client{
		r: resty.New().
			SetHostURL("https://api.uptimerobot.com/v2/").
			SetHeader("content-type", "application/x-www-form-urlencoded").
			SetHeader("cache-control", "no-cache"),
		key: key,
	}
}

// Response is returned by all API calls
type Response struct {
	Stat       string
	Pagination struct {
		Offset int
		Limit  int
		Total  int
	}
	Monitors []Monitor
}

// Monitor represents a single uptime monitor
type Monitor struct {
	ID             float64
	URL            string
	Port           string
	Status         float64
	Interval       float64
	FriendlyName   string
	CreateDatetime float64
	Type           float64
	KeywordValue   string
	KeywordType    string
	HTTPUsername   string
	HTTPPassword   string
	SubType        string
}

// GetMonitors returns all monitors for an account
func (c *Client) GetMonitors() (monitors []Monitor, err error) {
	var r Response
	_, err = c.r.R().SetBody(c.withAuth(url.Values{})).SetResult(&r).Post("getMonitors")
	if err != nil {
		return
	}

	return r.Monitors, nil
}

func (c Client) withAuth(params url.Values) (urlEncoded string) {
	params.Set("api_key", c.key)
	return params.Encode()
}
