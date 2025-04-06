package main

import (
	"context"
	"log"
	"os"

	"github.com/zelenin/go-tdlib/client"
)

func SendPhoto(my_client *client.Client, chatId int64, photoPath string) {
	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		log.Printf("File does not exist: %s", photoPath)
		return
	}

	messageContent := &client.InputMessagePhoto{
		Photo: &client.InputFileLocal{
			Path: photoPath,
		},
	}

	_, err := my_client.SendMessage(context.Background(), &client.SendMessageRequest{
		ChatId:              chatId,
		InputMessageContent: messageContent,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
	log.Println("Message sent")
}
