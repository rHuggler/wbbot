package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
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

// Authenticate authenticates against Twitch.tv's IRC server and joins a channel
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

// Listen listens for messages from Twitch.tv's IRC
func (b Bot) Listen() {
	reader := bufio.NewReader(b.Connection)
	tp := textproto.NewReader(reader)

	for {
		line, err := tp.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if strings.HasPrefix(line, "PING") {
			pongMessage := "PONG :tmi.twitch.tv"
			fmt.Fprintf(b.Connection, "%s\r\n", pongMessage)
		}

		if os.Getenv("DEBUG") != "" {
			fmt.Println(line)
		}
	}
}
