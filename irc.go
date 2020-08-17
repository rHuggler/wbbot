package main

import (
	"fmt"
	"log"
)

// IRCCommand is an enum for Twitch.tv's IRC commands
type IRCCommand int

const (
	unknown IRCCommand = iota
	ping
	pong
	pass
	nick
	join
	priv
)

// IRCMessage is an incoming message on Twitch.tv's IRC
type IRCMessage struct {
	Username  string
	Command   IRCCommand
	IsCommand bool
	Channel   string
	Message   string
}

func strToIRCCommand(str string) IRCCommand {
	mapper := map[string]IRCCommand{
		"PING":    ping,
		"PONG":    pong,
		"PASS":    pass,
		"NICK":    nick,
		"JOIN":    join,
		"PRIVMSG": priv,
	}

	command, ok := mapper[str]
	if !ok {
		return unknown
	}

	return command
}

func (i IRCCommand) String() string {
	mapper := map[IRCCommand]string{
		ping: "PING",
		pong: "PONG",
		pass: "PASS",
		nick: "NICK",
		join: "JOIN",
		priv: "PRIVMSG",
	}

	command, ok := mapper[i]
	if !ok {
		log.Fatal("Po ta unkjwno parcça")
	}

	return command
}

func (i IRCMessage) String() string {
	switch i.Command {
	case pong:
		return fmt.Sprintf("%s :tmi.twitch.tv", i.Command)
	case pass:
		return fmt.Sprintf("%s %s", i.Command, i.Message)
	case nick:
		return fmt.Sprintf("%s %s", i.Command, i.Username)
	case join:
		return fmt.Sprintf("%s #%s", i.Command, i.Channel)
	case priv:
		return fmt.Sprintf("%s #%s :%s", i.Command, i.Channel, i.Message)
	default:
		return ""
	}
}
