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

// MessageClient represents Message API client.
type MessageClient struct {
	ApplicationToken string       // application API token
	HTTPClient       *http.Client // if nil, http.DefaultClient is used
}

// NewMessageClient creates new Message API client.
func NewMessageClient(appToken string) (*MessageClient, error) {
	return &MessageClient{
		ApplicationToken: appToken,
	}, nil
}

func (c *MessageClient) http() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *MessageClient) makeData(message *Message) string {
	data := make(url.Values)

	// set required parameters
	data.Set("token", c.ApplicationToken)
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
func (c *MessageClient) SendMessage(ctx context.Context, message *Message) error {
	// prepare request
	body := strings.NewReader(c.makeData(message))
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.pushover.net/1/messages.json", body)
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

// Send is a shortcut for sending a basic message to given user.
func (c *MessageClient) Send(ctx context.Context, user, message string) error {
	m := &Message{
		User:    user,
		Message: message,
	}
	return c.SendMessage(ctx, m)
}
