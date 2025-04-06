package main

import (
	"context"
	"fmt"
	"os"

	"github.com/zelenin/go-tdlib/client"
)

func MainLoop(my_client *client.Client) {
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

			_, err = fmt.Scanln(&chatId)
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
