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
			from := getUserName(conn, msg)
			messages = append(messages, fmt.Sprintf("%s %s", from, text.Text.Text))
		}
	}

	return messages, nil
}

func getUserName(conn *connection.Connection, msg *tdlib.Message) string {
	var userLastName string
	var from string

	userId := msg.SenderId.(*tdlib.MessageSenderUser).UserId
	userName := senders[userId]
	if userName == "" {
		unkUser, err := conn.Client.GetUser(context.Background(), &tdlib.GetUserRequest{UserId: userId})
		if err != nil {
			userName = "unk"
		} else {
			userName, userLastName = unkUser.FirstName, unkUser.LastName
		}
	}
	if userLastName == "" {
		from = fmt.Sprintf("[%s]", userName)
	} else {
		from = fmt.Sprintf("[%s %s]", userName, userLastName)
	}

	return from
}
