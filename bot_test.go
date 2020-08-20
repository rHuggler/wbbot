package main

import (
	"container/ring"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gocarina/gocsv"
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

func TestRespond(t *testing.T) {
	testCases := []struct {
		want    IRCMessage
		message *IRCMessage
		timeNow time.Time
	}{
		{
			want: IRCMessage{
				Username:  "",
				Command:   priv,
				IsCommand: false,
				Channel:   "channel",
				Message:   "Os próximos World Bosses serão: (00:15) Great Jungle Wurm -- (00:30) Megadestroyer -- (00:45) Shadow Behemoth",
			},
			message: &IRCMessage{
				Username:  "rhuggler1",
				Command:   priv,
				IsCommand: true,
				Channel:   "rhuggler1",
				Message:   "wb",
			},
			timeNow: time.Date(0, 0, 0, 0, 0, 0, 0, time.Local),
		},
		{
			want: IRCMessage{
				Username:  "",
				Command:   priv,
				IsCommand: false,
				Channel:   "channel",
				Message:   "Os próximos World Bosses serão: (23:45) Fire Elemental -- (00:00) Admiral Taidha Covington, Tequatl the Sunless -- (00:15) Great Jungle Wurm",
			},
			message: &IRCMessage{
				Username:  "rhuggler1",
				Command:   priv,
				IsCommand: true,
				Channel:   "rhuggler1",
				Message:   "wb",
			},
			timeNow: time.Date(0, 0, 0, 23, 40, 0, 0, time.Local),
		},
	}

	b := NewBot(BotConfig{
		token:    "token",
		username: "username",
		channel:  "channel",
	})

	csvFile, err := os.OpenFile("world_boss_times.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	defer csvFile.Close()

	events := []*Event{}

	err = gocsv.UnmarshalFile(csvFile, &events)
	if err != nil {
		log.Fatal(err)
	}

	b.Events = events

	r := ring.New(len(b.Events))

	for i := 0; i < r.Len(); i++ {
		r.Value = b.Events[i]
		r = r.Next()
	}

	b.EventRotation = r

	for _, testCase := range testCases {
		got := b.Respond(testCase.message, testCase.timeNow)

		if !reflect.DeepEqual(testCase.want, got) {
			t.Errorf("Expected %#v, got %#v\n", testCase.want, got)
		}
	}

}
