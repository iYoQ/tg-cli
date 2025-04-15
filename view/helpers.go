package view

import (
	"context"
	"fmt"
	"slices"
	"tg-cli/connection"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func getChatHistory(conn *connection.Connection, chatId int64) ([]string, error) {
	_, err := conn.Client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: chatId})
	if err != nil {
		return nil, err
	}

	var history []*tdlib.Message
	fromMessageId := int64(0)

	for int32(len(history)) < chatLength {

		batch, err := conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
			ChatId:        chatId,
			FromMessageId: fromMessageId,
			Offset:        0,
			Limit:         min(pageSize, chatLength-int32(len(history))),
		})
		if err != nil {
			return nil, err
		}

		if len(batch.Messages) == 0 {
			break
		}

		history = append(history, batch.Messages...)
		fromMessageId = batch.Messages[len(batch.Messages)-1].Id
	}

	messagesIds := getMessagesIds(history)
	readMessages(conn.Client, chatId, messagesIds)

	var messages []string
	for _, msg := range slices.Backward(history) {
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

func closeChat(client *tdlib.Client, chatId int64) error {
	_, err := client.CloseChat(context.Background(), &tdlib.CloseChatRequest{ChatId: chatId})
	if err != nil {
		return err
	}

	return nil
}

func getMessagesIds(messages []*tdlib.Message) []int64 {
	ids := make([]int64, len(messages))

	for idx, msg := range messages {
		ids[idx] = msg.Id
	}

	return ids
}

func readMessages(client *tdlib.Client, chatId int64, messageIds []int64) error {
	_, err := client.ViewMessages(context.Background(), &tdlib.ViewMessagesRequest{
		ChatId:     chatId,
		MessageIds: messageIds,
	})
	if err != nil {
		return err
	}
	return nil
}
