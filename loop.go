package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zelenin/go-tdlib/client"
	"tg-cli/sender"
)

func ManualManager(my_client *client.Client) {
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. send msg")
		fmt.Println("2. get chat list")
		fmt.Println("9. exit")

		var choice int
		_, err := fmt.Scanln(&choice)

		if err != nil {
			fmt.Println("invalid")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("\nChoose a chat:")
			var chatId int64

			_, err = fmt.Scanln(&chatId)
			if err != nil {
				fmt.Println("invalid")
				continue
			}

			fmt.Println("\nChoose a msg:")
			var msg string

			_, err = fmt.Scanln(&msg)

			if err != nil {
				fmt.Println("invalid")
				continue
			}

			photoPath := strings.Split(msg, "=")
			if photoPath[0] == "ph" {
				sender.SendPhoto(my_client, chatId, strings.Join(photoPath[1:], "="))
			} else {
				sender.SendText(my_client, chatId, msg)
			}

		case 2:
			fmt.Println("\nChoose a number of chats:")
			var size int32
			_, err = fmt.Scanln(&size)
			if err != nil {
				fmt.Println("invalid")
				continue
			}

			chats, err := my_client.GetChats(context.Background(), &client.GetChatsRequest{Limit: size})
			if err != nil {
				log.Println("smh wrong with chats")
			}

			for _, id := range chats.ChatIds {
				chat, err := my_client.GetChat(context.Background(), &client.GetChatRequest{ChatId: id})
				if err != nil {
					log.Printf("smh wrong with chat %d, error %s", id, err)
					continue
				}
				fmt.Printf("name: %s, id: %d\n", chat.Title, chat.Id)
			}
		case 9:
			my_client.Close(context.Background())
			os.Exit(0)
		default:
			fmt.Println("invalid")
		}
	}
}
