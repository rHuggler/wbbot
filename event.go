package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Event is a world boss event
type Event struct {
	Time         string `csv:"time"`
	Boss         string `csv:"world_boss"`
	HardcoreBoss string `csv:"hardcore_world_boss"`
}

// GetTime parses HH:MM time to an UTC-3 Time
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

// GetMinutes parses an HH:MM timestamp and return the total of minutes in that timestamp
func (e Event) GetMinutes() (int, error) {
	slices := strings.Split(e.Time, ":")

	if len(slices) < 2 {
		return 0, errors.New("Could not split time")
	}

	hours, err := strconv.Atoi(slices[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(slices[1])
	if err != nil {
		return 0, err
	}

	totalMinutes := hours*60 + minutes

	return totalMinutes, nil
}

// GetMessage formats all event information in one string (current format: "(00:15) Boss, Hardcore Boss")
func (e Event) GetMessage() string {
	bosses := []string{
		e.Boss,
	}

	if e.HardcoreBoss != "" {
		bosses = append(bosses, e.HardcoreBoss)
	}

	joinedBosses := strings.Join(bosses, ", ")

	return fmt.Sprintf("(%s) %s", e.Time, joinedBosses)
}
