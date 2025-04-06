package sender

import (
	"context"
	"log"
	"os"

	"github.com/zelenin/go-tdlib/client"
)

func SendText(my_client *client.Client, chatId int64, msg string) {
	messageContent := &client.InputMessageText{
		Text: &client.FormattedText{
			Text: msg,
		},
	}

	messageRequest := buildRequest(chatId, messageContent)
	sendMessage(my_client, messageRequest)
}

func SendPhoto(my_client *client.Client, chatId int64, photoPath string) {
	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		log.Printf("File does not exist: %s", photoPath)
		return
	} else if err != nil {
		log.Printf("Error checking file: %v", err)
		return
	}

	messageContent := &client.InputMessagePhoto{
		Photo: &client.InputFileLocal{
			Path: photoPath,
		},
	}

	messageRequest := buildRequest(chatId, messageContent)
	sendMessage(my_client, messageRequest)
}

func buildRequest(chatId int64, content client.InputMessageContent) *client.SendMessageRequest {
	return &client.SendMessageRequest{
		ChatId:              chatId,
		InputMessageContent: content,
	}
}

func sendMessage(my_client *client.Client, messageRequest *client.SendMessageRequest) {
	_, err := my_client.SendMessage(context.Background(), messageRequest)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}

	log.Println("Message sent")
}
