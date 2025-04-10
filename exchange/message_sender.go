package exchange

import (
	"context"
	"log"
	"os"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func SendText(client *tdlib.Client, chatId int64, msg string) {
	messageContent := &tdlib.InputMessageText{
		Text: &tdlib.FormattedText{
			Text: msg,
		},
	}

	messageRequest := buildRequest(chatId, messageContent)
	sendMessage(client, messageRequest)
}

func SendPhoto(client *tdlib.Client, chatId int64, photoPath string, text string) {
	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		log.Printf("File does not exist: %s", photoPath)
		return
	} else if err != nil {
		log.Printf("Error checking file: %v", err)
		return
	}

	messageContent := &tdlib.InputMessagePhoto{
		Photo: &tdlib.InputFileLocal{
			Path: photoPath,
		},
		Caption: &tdlib.FormattedText{
			Text: text,
		},
	}

	messageRequest := buildRequest(chatId, messageContent)
	sendMessage(client, messageRequest)
}

func buildRequest(chatId int64, content tdlib.InputMessageContent) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:              chatId,
		InputMessageContent: content,
	}
}

func sendMessage(client *tdlib.Client, messageRequest *tdlib.SendMessageRequest) {
	_, err := client.SendMessage(context.Background(), messageRequest)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
}
