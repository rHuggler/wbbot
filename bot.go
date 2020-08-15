package main

import (
	"fmt"
	"log"
	"net"
)

// BotConfig provides Bot configuration parameters
type BotConfig struct {
	token    string
	username string
	channel  string
}

// Bot is a Twitch.tv chatbot
type Bot struct {
	Connection net.Conn
	Config     BotConfig
}

// NewBot returns a Bot pointer
func NewBot(config BotConfig) *Bot {
	c, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		Connection: c,
		Config:     config,
	}
}

// Authenticate authenticates against Twitch.tv's IRC server and joins the channel
func (b Bot) Authenticate() {
	passMessage := fmt.Sprintf("PASS %s", b.Config.token)
	fmt.Fprintf(b.Connection, "%s\r\n", passMessage)

	nickMessage := fmt.Sprintf("NICK %s", b.Config.username)
	fmt.Fprintf(b.Connection, "%s\r\n", nickMessage)

	joinMessage := fmt.Sprintf("JOIN #%s", b.Config.channel)
	fmt.Fprintf(b.Connection, "%s\r\n", joinMessage)
}

// SendMessage sends a message to the channel
func (b Bot) SendMessage(message string) {
	privMessage := fmt.Sprintf("PRIVMSG #%s :%s", b.Config.channel, message)
	fmt.Fprintf(b.Connection, "%s\r\n", privMessage)
}
