package main

import (
	"context"
	"github.com/zelenin/go-tdlib/client"
	"log"
)

func SendMessage(my_client *client.Client, chatId int64, msg string) {
	messageContent := &client.InputMessageText{
		Text: &client.FormattedText{
			Text: msg,
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
