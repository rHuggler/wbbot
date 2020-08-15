package main

import (
	"fmt"
	"net"
	"os"
)

// SendMessage sends a message to a Twitch channel
func SendMessage(message string) error {
	c, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		return err
	}

	passMessage := "PASS " + os.Getenv("OAUTH_TOKEN")
	fmt.Fprintf(c, "%s\r\n", passMessage)

	nickMessage := "NICK " + os.Getenv("TWITCH_USERNAME")
	fmt.Fprintf(c, "%s\r\n", nickMessage)

	joinMessage := "JOIN #" + os.Getenv("TWITCH_CHANNEL")
	fmt.Fprintf(c, "%s\r\n", joinMessage)

	privMessage := "PRIVMSG #" + os.Getenv("TWITCH_CHANNEL") + " :" + message
	fmt.Fprintf(c, "%s\r\n", privMessage)

	return nil
}
