package main

import (
	"flag"
	"fmt"
	"github.com/AlekSi/pushover"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)

	defApp := os.Getenv("PUSHOVER_APP")
	defUser := os.Getenv("PUSHOVER_USER")
	defDevice := os.Getenv("PUSHOVER_DEVICE")
	defTitle := os.Getenv("PUSHOVER_TITLE")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "%s [flags] [message]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags can also be read from environment variables listed below.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	app := flag.String("app", defApp, "application API token (PUSHOVER_APP)")
	user := flag.String("user", defUser, "user/group key (PUSHOVER_USER)")
	device := flag.String("device", defDevice, "device name to send the message directly to that device, rather than all of the user's devices (PUSHOVER_DEVICE)")
	title := flag.String("title", defTitle, "message title (PUSHOVER_TITLE)")
	url := flag.String("url", "", "supplementary URL")
	urlTitle := flag.String("url-title", "", "title for supplementary URL")
	priority := flag.Int("priority", 0, "priority")
	sound := flag.String("sound", "", "message sound")
	flag.Parse()

	msg := &pushover.Message{
		User:     *user,
		Message:  strings.Join(flag.Args(), " "),
		Device:   *device,
		Title:    *title,
		URL:      *url,
		URLTitle: *urlTitle,
		Priority: *priority,
		Sound:    *sound,
	}
	pushover.DefaultClient.ApplicationToken = *app
	err := pushover.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
