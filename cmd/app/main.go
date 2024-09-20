package main

import (
	"bot_hmb/internal"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Printf("Бот включен")
	internal.NewBotManager().JoinBot()
}
