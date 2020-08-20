package main

import (
	"container/ring"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	botConfig := BotConfig{
		token:    os.Getenv("OAUTH_TOKEN"),
		username: os.Getenv("TWITCH_USERNAME"),
		channel:  os.Getenv("TWITCH_CHANNEL"),
	}

	b := NewBot(botConfig)
	b.Authenticate()

	b.Events = events

	r := ring.New(len(b.Events))

	for i := 0; i < r.Len(); i++ {
		r.Value = b.Events[i]
		r = r.Next()
	}

	b.EventRotation = r

	go b.Listen()

	err = StartEventsCron(events, b)
	if err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	fmt.Println("Terminate by user")
}
