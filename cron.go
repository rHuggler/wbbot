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
		eventTime, err := event.GetTime()
		if err != nil {
			return err
		}

		// Subtracts 10 minutes
		eventTime = eventTime.Add(-10 * time.Minute)

		schedule := fmt.Sprintf("%d %d * * *", eventTime.Minute(), eventTime.Hour())

		message := ""
		if event.HardcoreBoss != "" {
			message = fmt.Sprintf("Os seguintes World Bosses estarão ativos em 10 minutos: %s and %s", event.Boss, event.HardcoreBoss)
		} else {
			message = fmt.Sprintf("O seguinte World Boss estará ativo em 10 minutos: %s", event.Boss)
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
