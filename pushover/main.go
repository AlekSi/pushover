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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage\n")
		fmt.Fprintf(os.Stderr, "%s [flags] [message]:\n", os.Args[0])
		flag.PrintDefaults()
	}
	app := flag.String("app", "", "application API token")
	user := flag.String("user", "", "user/group key")
	flag.Parse()

	message := strings.Join(flag.Args(), " ")
	pushover.DefaultClient.ApplicationToken = *app
	err := pushover.SendMessage(*user, message)
	if err != nil {
		log.Fatal(err)
	}
}
