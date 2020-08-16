package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// StartEventsCron creates cron jobs for every event and starts the cron scheduler
func StartEventsCron(events []*Event, b *Bot) error {
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return err
	}

	c := cron.New(cron.WithLocation(location))

	for _, event := range events {
		time, err := event.GetTime()
		if err != nil {
			return err
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
			return err
		}
	}

	fmt.Printf("Added %d cron entries\n", len(events))

	c.Start()
	return nil
}
