package main

import (
	"log"

	"github.com/SicParv1sMagna/AtomHackMarsService/internal/app"
)

func main() {
	log.Println("Application start!")

	application, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	application.Run()
	log.Println("Application terminated!")
}
