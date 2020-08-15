package main

import (
	"fmt"
	"log"
	"os"
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

	for _, event := range events {
		time, err := event.GetTime()
		if err != nil {
			log.Fatal(err)
		}

		schedule := fmt.Sprintf("0 %d %d * * *", time.Minute(), time.Hour())

		message := ""
		if event.HardcoreBoss != "" {
			message = fmt.Sprintf("Active world bosses: %s and %s", event.Boss, event.HardcoreBoss)
		} else {
			message = fmt.Sprintf("Active world boss: %s", event.Boss)
		}

		fmt.Printf("Crontab: %s\nMessage: %s\n", schedule, message)

		c.AddFunc(schedule, func() {
			SendMessage(message)
		})
	}

	c.Start()
}
