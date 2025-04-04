package main

import (
	"fmt"
	// "log"
	// "github.com/zelenin/go-tdlib/client"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	my_client := Auth()
	go HandleShutDown((my_client))

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. send msg")
		fmt.Println("2. exit")

		var choice string
		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("invalid")
			continue
		}

		switch choice {
		case "1":
			fmt.Println("\nChoose an chat:")
			var chatId int64

			_, err := fmt.Scanln(&chatId)
			if err != nil {
				fmt.Println("invalid")
				continue
			}

			fmt.Println("\nChoose an msg:")
			var msg string

			_, err = fmt.Scanln(&msg)

			if err != nil {
				fmt.Println("invalid")
				continue
			}

			SendMessage(my_client, chatId, msg)
		case "2":
			my_client.Close(context.Background())
			os.Exit(0)
		default:
			fmt.Println("invalid")
		}
	}

}
