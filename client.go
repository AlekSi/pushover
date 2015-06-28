package pushover

import (
	"encoding/json"
	"io/ioutil"
	"net"
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
	PushoverSound     = "pushover"
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
	NoneSound         = "none"
)

// Message to send.
type Message struct {
	// mandatory parameters
	User    string // user/group key
	Message string // message to send

	// optional parameters
	Device    string    // device name to send the message directly to that device, rather than all of the user's devices
	Title     string    // message title
	URL       string    // supplementary URL
	URLTitle  string    // title for supplementary URL
	Priority  int       // priority, defaults to NormalPriority
	Timestamp time.Time // message time
	Sound     string    // message sound

	// for emergency priority only
	Retry    int
	Expire   int
	Callback string
}

// Client sends messages.
type Client struct {
	ApplicationToken string       // application API token
	HTTPClient       *http.Client // if nil, http.DefaultClient is used
}

func (c *Client) http() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return http.DefaultClient
}

func (c *Client) makeData(message *Message) string {
	data := make(url.Values)

	// set required parameters
	data.Set("token", c.ApplicationToken)
	data.Set("user", message.User)
	data.Set("message", message.Message)

	// set optional parameters
	if message.Device != "" {
		data.Set("device", message.Device)
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
	if !message.Timestamp.IsZero() {
		data.Set("timestamp", strconv.FormatInt(message.Timestamp.Unix(), 10))
	}
	if message.Sound != "" {
		data.Set("sound", message.Sound)
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

// Send sends message.
// It does not retries failed sends.
// Returns either nil, TemporaryError or FatalError.
func (c *Client) Send(message *Message) error {
	// prepare request
	body := strings.NewReader(c.makeData(message))
	req, err := http.NewRequest("POST", "https://api.pushover.net/1/messages.json", body)
	if err != nil {
		return &FatalError{-1, err.Error()}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "github.com/AlekSi/pushover")

	// do request and read body
	resp, err := c.http().Do(req)
	if err != nil {
		return &TemporaryError{-1, err.Error()}
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &TemporaryError{resp.StatusCode, err.Error()}
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

	if resp.StatusCode/100 == 4 || (jsonOk && status == 0.0) {
		return &FatalError{StatusCode: resp.StatusCode, Message: string(b)}
	}

	return &TemporaryError{StatusCode: resp.StatusCode, Message: string(b)}
}

// SendWithRetries sends message.
// It does retries failed sends for temporary errors up to maxRetries times with 5 second delay.
// Specify maxRetries <= 0 for unlimited retries.
// Returns either nil, TemporaryError (if gave up) or FatalError.
func (c *Client) SendWithRetries(message *Message, maxRetries int) error {
	var i int
	for {
		err := c.Send(message)
		if e, ok := err.(net.Error); !ok || !e.Temporary() {
			return err
		}

		i++
		if maxRetries > 0 && maxRetries == i {
			return err
		}

		time.Sleep(5 * time.Second)
	}
}

// SendMessage sends message to specified user.
// It does not retries failed sends.
// Returns either nil, TemporaryError or FatalError.
func (c *Client) SendMessage(user, message string) error {
	return c.Send(&Message{User: user, Message: message, Timestamp: time.Now()})
}

// SendMessageWithRetries sends message to specified user.
// It does retries failed sends for temporary errors up to maxRetries times with 5 second delay.
// Specify maxRetries <= 0 for unlimited retries.
// Returns either nil, TemporaryError (if gave up) or FatalError.
func (c *Client) SendMessageWithRetries(user, message string, maxRetries int) error {
	return c.SendWithRetries(&Message{User: user, Message: message, Timestamp: time.Now()}, maxRetries)
}
