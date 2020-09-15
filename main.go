package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

type config struct {
	Username string
	Password string
	Channels []string
	Message  string
}

func parseConfig() (config, error) {
	c := config{}
	var channels string
	flag.StringVar(&c.Username, "username", "", "twitch username")
	flag.StringVar(&c.Password, "password", "", "twitch 'password' (oauth token)")
	flag.StringVar(&channels, "channels", "", "list of channels to join")
	flag.StringVar(&c.Message, "message", "Please chat in our Discord chat at https://discord.gg/ponyfest instead of here. Thanks!", "Message to send in reply")
	flag.Parse()

	if c.Username == "" {
		return c, fmt.Errorf("--username is required")
	}
	if c.Password == "" {
		return c, fmt.Errorf("--password is required")
	}
	c.Channels = strings.Split(channels, ",")
	return c, nil
}

func main() {
	c, err := parseConfig()
	if err != nil {
		log.Fatalf("error: %v.\n", err)
	}
	client := twitch.NewClient(c.Username, c.Password)

	client.OnConnect(func() {
		log.Printf("Connected! Joining %s...", strings.Join(c.Channels, ", "))
		client.Join(c.Channels...)
	})
	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		client.Say(m.Channel, fmt.Sprintf("/delete %s", m.ID))
		client.Say(m.Channel, c.Message)
	})

	if err := client.Connect(); err != nil {
		log.Fatalf("Couldn't connect to Twitch: %v.\n", err)
	}
}
