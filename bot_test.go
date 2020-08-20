package main

import (
	"reflect"
	"testing"
)

func TestParseMessage(t *testing.T) {
	testMessage := ":rhuggler1!rhuggler1@rhuggler1.tmi.twitch.tv PRIVMSG #rhuggler1 :foda meu amigo"
	b := NewBot(BotConfig{
		token:    "token",
		username: "user",
		channel:  "chan",
	})

	want := &IRCMessage{
		Username:  "rhuggler1",
		Command:   priv,
		IsCommand: false,
		Channel:   "rhuggler1",
		Message:   "foda meu amigo",
	}

	got := b.parseMessage(testMessage)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Expected %#v, got %#v", want, got)
	}
}

func TestParseMessageTable(t *testing.T) {
	testCases := []struct {
		want    *IRCMessage
		message string
	}{
		{
			want: &IRCMessage{
				Username:  "rhuggler1",
				Command:   priv,
				IsCommand: true,
				Channel:   "rhuggler1",
				Message:   "wb",
			},
			message: ":rhuggler1!rhuggler1@rhuggler1.tmi.twitch.tv PRIVMSG #rhuggler1 :!wb",
		},
		{
			want: &IRCMessage{
				Username:  "worldbossbot",
				Command:   join,
				IsCommand: false,
				Channel:   "rhuggler1",
				Message:   "",
			},
			message: ":worldbossbot!worldbossbot@worldbossbot.tmi.twitch.tv JOIN #rhuggler1",
		},
		{
			want: &IRCMessage{
				Username:  "",
				Command:   ping,
				IsCommand: false,
				Channel:   "",
				Message:   "tmi.twitch.tv",
			},
			message: "PING :tmi.twitch.tv",
		},
	}

	b := NewBot(BotConfig{
		token:    "token",
		username: "username",
		channel:  "channel",
	})

	for _, testCase := range testCases {
		got := b.parseMessage(testCase.message)

		if !reflect.DeepEqual(testCase.want, got) {
			t.Errorf("Expected %#v, got %#v", testCase.want, got)
		}
	}
}
