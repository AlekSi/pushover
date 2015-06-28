package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/AlekSi/pushover"
)

func main() {
	log.SetFlags(0)

	defApp := os.Getenv("PUSHOVER_APP")
	defDevice := os.Getenv("PUSHOVER_DEVICE")
	defMaxRetries, _ := strconv.Atoi(os.Getenv("PUSHOVER_MAX_RETRIES"))
	defTitle := os.Getenv("PUSHOVER_TITLE")
	defUser := os.Getenv("PUSHOVER_USER")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "%s [flags] [message]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags can also be read from environment variables listed below.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	app := flag.String("app", defApp, "application API token (PUSHOVER_APP)")
	device := flag.String("device", defDevice, "device name to send the message directly to that device, rather than all of the user's devices (PUSHOVER_DEVICE)")
	maxRetries := flag.Int("max-retries", defMaxRetries, "max retries, 0 for unlimited (PUSHOVER_MAX_RETRIES)")
	priority := flag.Int("priority", 0, "priority")
	sound := flag.String("sound", "", "message sound")
	title := flag.String("title", defTitle, "message title (PUSHOVER_TITLE)")
	url := flag.String("url", "", "supplementary URL")
	urlTitle := flag.String("url-title", "", "title for supplementary URL")
	user := flag.String("user", defUser, "user/group key (PUSHOVER_USER)")
	flag.Parse()

	msg := &pushover.Message{
		Device:   *device,
		Message:  strings.Join(flag.Args(), " "),
		Priority: *priority,
		Sound:    *sound,
		Title:    *title,
		URL:      *url,
		URLTitle: *urlTitle,
		User:     *user,
	}
	pushover.DefaultClient.ApplicationToken = *app
	err := pushover.SendWithRetries(msg, *maxRetries)
	if err != nil {
		log.Fatal(err)
	}
}
