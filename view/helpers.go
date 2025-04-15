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
		from := getUserName(conn, msg)
		formatMsg := processMessages(msg, from)
		messages = append(messages, formatMsg)
	}

	return messages, nil
}

func getUserName(conn *connection.Connection, msg *tdlib.Message) string {
	var userLastName string
	var from string
	var senderId int64

	switch sender := msg.SenderId.(type) {
	case *tdlib.MessageSenderUser:
		senderId = sender.UserId
	case *tdlib.MessageSenderChat:
		senderId = sender.ChatId
	}

	userName := senders[senderId]
	if userName == "" {
		unkUser, err := conn.Client.GetUser(context.Background(), &tdlib.GetUserRequest{UserId: senderId})
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

func processMessages(msg *tdlib.Message, from string) string {
	var result string
	switch content := msg.Content.(type) {
	case *tdlib.MessageText:
		result = fmt.Sprintf("%s %s", from, content.Text.Text)
	case *tdlib.MessagePhoto, *tdlib.MessageVideo, *tdlib.MessageAudio:
		result = fmt.Sprintf("%s [media content]", from)
	}

	return result
}
