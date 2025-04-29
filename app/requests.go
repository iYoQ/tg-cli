package app

import (
	"context"
	"slices"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func getChatHistory(client *tdlib.Client, chatId int64, threadId int64, chatLoadSize int32) ([]string, error) {
	_, err := client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: chatId})
	if err != nil {
		return nil, err
	}

	var history []*tdlib.Message
	fromMessageId := int64(0)

	for int32(len(history)) < chatLoadSize {

		batch, err := client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
			ChatId:        chatId,
			FromMessageId: fromMessageId,
			Offset:        0,
			Limit:         chatLoadSize - int32(len(history)),
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
	if messagesIds != nil {
		readMessages(client, chatId, messagesIds)
	}

	var messages []string
	for _, msg := range slices.Backward(history) {
		if threadId > 0 {
			if threadId == msg.MessageThreadId {
				from := getUserName(client, msg)
				formatMsg := processMessages(msg, from)
				messages = append(messages, formatMsg)
			}
		} else {
			from := getUserName(client, msg)
			formatMsg := processMessages(msg, from)
			messages = append(messages, formatMsg)
		}
	}

	return messages, nil
}

func getUserName(client *tdlib.Client, msg *tdlib.Message) string {
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
		unkUser, err := client.GetUser(context.Background(), &tdlib.GetUserRequest{UserId: senderId})
		if err != nil {
			userName = unkSenderStyle.Render(unknownIdentifier)
			from = userName
		} else {
			userName, userLastName = unkUser.FirstName, unkUser.LastName
			if userLastName == "" {
				senders[senderId] = senderStyle.Render(userName)
			} else {
				senders[senderId] = senderStyle.Render(userName, userLastName)
			}
		}
	}

	if from == "" {
		from = senders[senderId]
	}

	return from
}

func closeChat(client *tdlib.Client, chatId int64) error {
	_, err := client.CloseChat(context.Background(), &tdlib.CloseChatRequest{ChatId: chatId})
	if err != nil {
		return err
	}

	return nil
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
