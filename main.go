package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Not found .env file")
	}

	my_client := Auth()

	MainLoop(my_client)

}
