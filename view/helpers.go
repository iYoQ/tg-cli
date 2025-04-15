package view

import (
	"context"
	"fmt"
	"slices"
	"tg-cli/connection"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func getChatHistory(conn *connection.Connection, chatId int64) ([]string, error) {
	lastMssage, err := conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: 0,
		Offset:        0,
		Limit:         1,
	})
	if err != nil {
		return nil, err
	}

	history, err := conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: lastMssage.Messages[0].Id,
		Offset:        0,
		Limit:         chatLength,
	})
	if err != nil {
		return nil, err
	}

	var messages []string
	for _, msg := range slices.Backward(history.Messages) {
		if text, ok := msg.Content.(*tdlib.MessageText); ok {
			from := fmt.Sprintf("[%s]", senders[msg.SenderId.(*tdlib.MessageSenderUser).UserId])
			messages = append(messages, fmt.Sprintf("%s %s", from, text.Text.Text))
		}
	}

	return messages, nil
}
