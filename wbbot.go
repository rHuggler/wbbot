package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
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

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err)
	}

	c := cron.New(cron.WithLocation(location))

	fmt.Printf("Adding %d entries\n", len(events))

	botConfig := BotConfig{
		token:    os.Getenv("OAUTH_TOKEN"),
		username: os.Getenv("TWITCH_USERNAME"),
		channel:  os.Getenv("TWITCH_CHANNEL"),
	}

	b := NewBot(botConfig)
	b.Authenticate()

	for _, event := range events {
		time, err := event.GetTime()
		if err != nil {
			log.Fatal(err)
		}

		schedule := fmt.Sprintf("%d %d * * *", time.Minute(), time.Hour())

		message := ""
		if event.HardcoreBoss != "" {
			message = fmt.Sprintf("Active world bosses: %s and %s", event.Boss, event.HardcoreBoss)
		} else {
			message = fmt.Sprintf("Active world boss: %s", event.Boss)
		}

		_, err = c.AddFunc(schedule, func() {
			b.SendMessage(message)
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	c.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	fmt.Println("Terminate by user")
}
