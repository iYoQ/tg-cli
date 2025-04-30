package requests

import (
	"context"
	"log"
	"os"

	tdlib "github.com/zelenin/go-tdlib/client"
)

type Params struct {
	ChatId   int64
	ThreadId int64
	Msg      string
	FilePath string
}

func SendText(client *tdlib.Client, params Params) {
	messageContent := &tdlib.InputMessageText{
		Text: &tdlib.FormattedText{
			Text: params.Msg,
		},
	}

	messageRequest := buildRequest(messageContent, params)
	sendMessage(client, messageRequest)
}

func SendPhoto(client *tdlib.Client, params Params) error {
	if _, err := os.Stat(params.FilePath); os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}

	messageContent := &tdlib.InputMessagePhoto{
		Photo: &tdlib.InputFileLocal{
			Path: params.FilePath,
		},
		Caption: &tdlib.FormattedText{
			Text: params.Msg,
		},
	}

	messageRequest := buildRequest(messageContent, params)
	sendMessage(client, messageRequest)
	return nil
}

func SendFile(client *tdlib.Client, params Params) error {
	if _, err := os.Stat(params.FilePath); os.IsNotExist(err) {
		return err
	} else if err != nil {
		return err
	}

	messageContent := &tdlib.InputMessageDocument{
		Document: &tdlib.InputFileLocal{
			Path: params.FilePath,
		},
		Caption: &tdlib.FormattedText{
			Text: params.Msg,
		},
	}

	messageRequest := buildRequest(messageContent, params)
	sendMessage(client, messageRequest)
	return nil
}

func buildRequest(content tdlib.InputMessageContent, params Params) *tdlib.SendMessageRequest {
	return &tdlib.SendMessageRequest{
		ChatId:              params.ChatId,
		InputMessageContent: content,
		MessageThreadId:     params.ThreadId,
	}
}

func sendMessage(client *tdlib.Client, messageRequest *tdlib.SendMessageRequest) {
	_, err := client.SendMessage(context.Background(), messageRequest)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
}
