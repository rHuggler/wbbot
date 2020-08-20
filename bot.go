package main

import (
	"bufio"
	"container/ring"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"
)

// BotConfig provides Bot configuration parameters
type BotConfig struct {
	token    string
	username string
	channel  string
}

// Bot is a Twitch.tv chatbot
type Bot struct {
	Connection    net.Conn
	Config        BotConfig
	Events        []*Event
	EventRotation *ring.Ring
}

var (
	messageRegex = regexp.MustCompile(`(?::(?P<username>\w+)!\w+@(?:\w+\.?)+\s)?(?P<command>\w+)\s(?:#(?P<channel>\w+)\s?)?(?::(?P<isCommand>\!)?(?P<message>.+\s*)+)?`)
)

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
	b.SendMessage(&IRCMessage{
		Command: pass,
		Message: b.Config.token,
	})

	b.SendMessage(&IRCMessage{
		Username: b.Config.username,
		Command:  nick,
	})

	b.SendMessage(&IRCMessage{
		Command: join,
		Channel: b.Config.channel,
	})
}

// SendMessage sends an IRCMessage to Twitch.tv's IRC
func (b Bot) SendMessage(message *IRCMessage) {
	fmt.Printf("%s\n", message)
	fmt.Fprintf(b.Connection, "%s\r\n", message)
}

// Respond responds to user commands
func (b Bot) Respond(message *IRCMessage, timeNow time.Time) IRCMessage {
	if message.Message == "wb" {
		nextEvents := []string{}

		for i := 0; i < b.EventRotation.Len(); i++ {
			event := b.EventRotation.Value.(*Event)

			eventMinutes, err := event.GetMinutes()
			if err != nil {
				log.Fatal(err)
			}

			minutesNow := timeNow.Hour()*60 + timeNow.Minute()

			if eventMinutes > minutesNow {
				eventMessage := fmt.Sprintf("(%s) %s", event.Time, event.Boss)
				if event.HardcoreBoss != "" {
					eventMessage += fmt.Sprintf(", %s", event.HardcoreBoss)
				}
				nextEvents = append(nextEvents, eventMessage)
				b.EventRotation = b.EventRotation.Next()
				break
			}

			b.EventRotation = b.EventRotation.Next()
		}

		for i := 0; i < 2; i++ {
			event := b.EventRotation.Value.(*Event)

			eventMessage := fmt.Sprintf("(%s) %s", event.Time, event.Boss)
			if event.HardcoreBoss != "" {
				eventMessage += fmt.Sprintf(", %s", event.HardcoreBoss)
			}
			nextEvents = append(nextEvents, eventMessage)

			b.EventRotation = b.EventRotation.Next()
		}

		nextEventsMessage := "Os próximos World Bosses serão: " + strings.Join(nextEvents, " -- ")

		return IRCMessage{
			Command: priv,
			Channel: b.Config.channel,
			Message: nextEventsMessage,
		}
	}

	return IRCMessage{}
}

func (b Bot) parseMessage(line string) *IRCMessage {
	matches := messageRegex.FindStringSubmatch(line)
	if matches == nil {
		return &IRCMessage{}
	}

	message := &IRCMessage{
		Username: matches[1],
		Command:  strToIRCCommand(matches[2]),
		Channel:  matches[3],
		Message:  matches[5],
	}

	if matches[4] == "!" {
		message.IsCommand = true
	}

	return message
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

		message := b.parseMessage(line)

		if message.Command == ping {
			b.SendMessage(&IRCMessage{
				Command: pong,
			})
		}

		if message.IsCommand {
			response := b.Respond(message, time.Now())
			b.SendMessage(&response)
		}

		if os.Getenv("DEBUG") != "" {
			fmt.Println(line)
		}
	}
}
