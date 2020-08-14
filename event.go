package main

import "time"

// Event is a world boss event
type Event struct {
	Time         string `csv:"time"`
	Boss         string `csv:"world_boss"`
	HardcoreBoss string `csv:"hardcore_world_boss"`
}

// GetTime parses HH:MM time to an UTC-3 datetime
func (e Event) GetTime() (time.Time, error) {
	eventTime, err := time.Parse("15:04", e.Time)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Time{}, err
	}

	currentTime := time.Now()

	eventTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), eventTime.Hour(), eventTime.Minute(), 0, 0, location)

	return eventTime, nil
}
