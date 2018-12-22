package uptimerobot

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Southclaws/qstring"
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
	Error struct {
		Type          string `json:"type"`
		ParameterName string `json:"parameter_name"`
		PassedValue   string `json:"passed_value"`
		Message       string `json:"message"`
	}
	Monitors []Monitor
	Monitor  struct {
		ID     int
		Status int
	}
}

// Type represents a monitor type
type Type int

// Valid monitor types
const (
	MonitorTypeInvalid Type = 0
	MonitorTypeHTTP    Type = 1
	MonitorTypeKeyword Type = 2
	MonitorTypePing    Type = 3
	MonitorTypePort    Type = 4
)

// SubType is used only for "Port monitoring (monitor>type = 4)" and shows which
// pre-defined port/service is monitored or if a custom port is monitored.
type SubType int

// Valid subtypes
const (
	SubTypeHTTP       SubType = 1 // 80
	SubTypeHTTPS      SubType = 2 // 443
	SubTypeFTP        SubType = 3 // 21
	SubTypeSMTP       SubType = 4 // 25
	SubTypePOP3       SubType = 5 // 110
	SubTypeIMAP       SubType = 6 // 143
	SubTypeCustomPort SubType = 99
)

// Monitor represents a single uptime monitor
type Monitor struct {
	ID  int    `json:"id"                        qstring:"id"`
	URL string `json:"url,omitempty"             qstring:"url,omitempty"`
	// Port           int     `json:"port,omitempty"            qstring:"port,omitempty"`
	Status         int    `json:"status,omitempty"          qstring:"status,omitempty"`
	Interval       int    `json:"interval,omitempty"        qstring:"interval,omitempty"`
	FriendlyName   string `json:"friendly_name,omitempty"   qstring:"friendly_name,omitempty"`
	CreateDatetime int    `json:"create_datetime,omitempty" qstring:"create_datetime,omitempty"`
	Type           Type   `json:"type,omitempty"            qstring:"type,omitempty"`
	// KeywordValue   string  `json:"keyword_value,omitempty"   qstring:"keyword_value,omitempty"`
	// KeywordType    int     `json:"keyword_type,omitempty"    qstring:"keyword_type,omitempty"`
	// HTTPUsername   string  `json:"http_username,omitempty"   qstring:"http_username,omitempty"`
	// HTTPPassword   string  `json:"http_password,omitempty"   qstring:"http_password,omitempty"`
	// SubType SubType `json:"sub_type,omitempty"        qstring:"sub_type,omitempty"`
}

// Validate ensures a monitor payload has all required fields
func (m Monitor) Validate() (err error) {
	if m.URL == "" {
		err = errors.New("missing url")
	}
	if m.FriendlyName == "" {
		err = errors.New("missing friendly_name")
	}
	if m.Type == 0 {
		err = errors.New("missing type")
	}
	return
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

// NewMonitor creates a new uptime monitor
func (c *Client) NewMonitor(monitor Monitor) (id int, err error) {
	if err = monitor.Validate(); err != nil {
		return
	}
	params, err := qstring.Marshal(&monitor)
	if err != nil {
		return
	}
	var r Response
	_, err = c.r.R().SetBody(c.withAuth(params)).SetResult(&r).Post("newMonitor")
	if err != nil {
		return
	}

	if r.Stat != "ok" {
		err = errors.New(r.Error.Message)
	}

	id = r.Monitor.ID

	return
}

// DeleteMonitor removes an uptime monitor
func (c *Client) DeleteMonitor(id int) (err error) {
	var r Response
	_, err = c.r.R().SetBody(c.withAuth(url.Values{"id": []string{fmt.Sprint(id)}})).SetResult(&r).Post("deleteMonitor")
	if err != nil {
		return
	}

	if r.Stat != "ok" {
		err = errors.New(r.Error.Message)
	}

	return
}

func (c Client) withAuth(params url.Values) (urlEncoded string) {
	params.Set("api_key", c.key)
	return params.Encode()
}
