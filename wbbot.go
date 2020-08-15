package main

import (
	"log"
	"os"

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

	// message := "Hi"

	// err = SendMessage(accessToken, message)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ServeAuthentication()

	// c := cron.New()

	// c.Start()
}
