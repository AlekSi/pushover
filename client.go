// Package pushover provides client for Pushover API (http://pushover.net/).
package pushover

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Message priority.
const (
	LowestPriority    = -2 // lowest priority, no notification
	LowPriority       = -1 // low priority, no sound and vibration
	NormalPriority    = 0  // normal priority, default
	HighPriority      = 1  // high priority, always with sound and vibration
	EmergencyPriority = 2  // emergency priority, requires acknowledge
)

// Message sound.
const (
	PushoverSound     = "pushover" // default
	BikeSound         = "bike"
	BugleSound        = "bugle"
	CashregisterSound = "cashregister"
	ClassicalSound    = "classical"
	CosmicSound       = "cosmic"
	FallingSound      = "falling"
	GamelanSound      = "gamelan"
	IncomingSound     = "incoming"
	IntermissionSound = "intermission"
	MagicSound        = "magic"
	MechanicalSound   = "mechanical"
	PianobarSound     = "pianobar"
	SirenSound        = "siren"
	SpacealarmSound   = "spacealarm"
	TugboatSound      = "tugboat"
	AlienSound        = "alien"
	ClimbSound        = "climb"
	PersistentSound   = "persistent"
	EchoSound         = "echo"
	UpdownSound       = "updown"
	VibrateSound      = "vibrate" // vibrate only
	NoneSound         = "none"    // silent
)

// Message to send.
type Message struct {
	// mandatory parameters
	User    string // user/group key
	Message string // message to send

	// optional parameters
	Devices   []string  // device names to send the message directly to that devices, rather than all of the user's devices
	Title     string    // message title, defaults to application name
	URL       string    // supplementary URL
	URLTitle  string    // title for supplementary URL
	Priority  int       // priority, defaults to NormalPriority
	Sound     string    // message sound
	Timestamp time.Time // message time
	HTML      bool      // enable HTML formatting
	Monospace bool      // enable monospace messages

	// for emergency priority only
	Retry    int
	Expire   int
	Callback string
}

// Client represents Pushover API client.
//
// See https://pushover.net/api.
type Client struct {
	appToken string

	m          sync.RWMutex
	httpClient *http.Client
}

// NewClient creates new client.
func NewClient(appToken string) (*Client, error) {
	return &Client{
		appToken: appToken,
	}, nil
}

func (c *Client) SetHTTPClient(client *http.Client) {
	c.m.Lock()
	defer c.m.Unlock()

	c.httpClient = client
}

func (c *Client) http() *http.Client {
	c.m.RLock()
	defer c.m.RUnlock()

	if c.httpClient != nil {
		return c.httpClient
	}
	return http.DefaultClient
}

func (c *Client) sendRequest(ctx context.Context, URL string, data string) error {
	// prepare request
	body := strings.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, "POST", URL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "github.com/AlekSi/pushover")

	// do request and read body
	resp, err := c.http().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// parse response
	var jsonOk bool
	var status float64
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err == nil {
		status, jsonOk = m["status"].(float64)
	}

	if resp.StatusCode == 200 && jsonOk && status == 1.0 {
		return nil
	}

	return fmt.Errorf("%d: %s", resp.StatusCode, b)
}

func (c *Client) makeMessageData(message *Message) string {
	data := make(url.Values)

	// set required parameters
	data.Set("token", c.appToken)
	data.Set("user", message.User)
	data.Set("message", message.Message)

	// set optional parameters
	if len(message.Devices) != 0 {
		data.Set("device", strings.Join(message.Devices, ","))
	}
	if message.Title != "" {
		data.Set("title", message.Title)
	}
	if message.URL != "" {
		data.Set("url", message.URL)
	}
	if message.URLTitle != "" {
		data.Set("url_title", message.URLTitle)
	}
	if message.Priority != 0 {
		data.Set("priority", strconv.Itoa(message.Priority))
	}
	if message.Sound != "" {
		data.Set("sound", message.Sound)
	}
	if !message.Timestamp.IsZero() {
		data.Set("timestamp", strconv.FormatInt(message.Timestamp.Unix(), 10))
	}
	if message.HTML {
		data.Set("html", "1")
	}
	if message.Monospace {
		data.Set("monospace", "1")
	}

	// set parameters for emergency priority
	if message.Priority == EmergencyPriority {
		data.Set("retry", strconv.Itoa(message.Retry))
		data.Set("expire", strconv.Itoa(message.Expire))
		if message.Callback != "" {
			data.Set("callback", message.Callback)
		}
	}

	return data.Encode()
}

// SendMessage sends given message.
func (c *Client) SendMessage(ctx context.Context, message *Message) error {
	return c.sendRequest(ctx, "https://api.pushover.net/1/messages.json", c.makeMessageData(message))
}

// Send is a shortcut for sending a basic message to given user.
func (c *Client) Send(ctx context.Context, user, message string) error {
	m := &Message{
		User:    user,
		Message: message,
	}
	return c.SendMessage(ctx, m)
}

type Glance struct {
	User string

	Device string

	// optional parameters
	Title   *string
	Text    *string
	Subtext *string
	Count   *int
	Percent *int
}

var (
	RemoveCount   = new(int)
	RemovePercent = new(int)
)

func (c *Client) makeGlanceData(glance *Glance) string {
	data := make(url.Values)

	data.Set("token", c.appToken)
	data.Set("user", glance.User)

	if glance.Device != "" {
		data.Set("device", glance.Device)
	}

	if glance.Title != nil {
		data.Set("title", *glance.Title)
	}
	if glance.Text != nil {
		data.Set("text", *glance.Text)
	}
	if glance.Subtext != nil {
		data.Set("subtext", *glance.Subtext)
	}
	if glance.Count != nil {
		var count string
		if glance.Count != RemoveCount {
			count = fmt.Sprint(*glance.Count)
		}
		data.Set("count", count)
	}
	if glance.Percent != nil {
		var percent string
		if glance.Percent != RemovePercent {
			percent = fmt.Sprint(*glance.Percent)
		}
		data.Set("percent", percent)
	}

	return data.Encode()
}

func (c *Client) SendGlance(ctx context.Context, glance *Glance) error {
	return c.sendRequest(ctx, "https://api.pushover.net/1/glances.json", c.makeGlanceData(glance))
}
