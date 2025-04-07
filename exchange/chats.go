package exchange

import (
	"context"
	"fmt"
	"log"

	"github.com/zelenin/go-tdlib/client"
)

func GetChats(my_client *client.Client, size int32) {
	chats, err := my_client.GetChats(context.Background(), &client.GetChatsRequest{Limit: size})
	if err != nil {
		log.Println("smh wrong with chats")
		return
	}

	for _, id := range chats.ChatIds {
		chat, err := my_client.GetChat(context.Background(), &client.GetChatRequest{ChatId: id})
		if err != nil {
			log.Printf("smh wrong with chat %d, error %s", id, err)
			continue
		}

		fmt.Printf("name: %s, id: %d\n", chat.Title, chat.Id)
		fmt.Println("-----------------------------------------------------------")
	}
}

// Переработать этот пиздец, добавить идентификатор того кто отправлял сообщение
func GetMessages(my_client *client.Client, chatId int64) {
	_, err := my_client.OpenChat(context.Background(), &client.OpenChatRequest{ChatId: chatId})
	if err != nil {
		log.Printf("Cannot open chat, %s", err)
		return
	}

	messages, err := my_client.GetChatHistory(context.Background(), &client.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: 0,
		Offset:        0,
		Limit:         1,
	})
	if err != nil {
		log.Printf("Cannot receive messages, %s", err)
		return
	}

	moreMsg, err := my_client.GetChatHistory(context.Background(), &client.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: messages.Messages[0].Id,
		Offset:        -1,
		Limit:         10,
	})
	if err != nil {
		log.Printf("Cannot receive messages, %s", err)
		return
	}

	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("Chat history, last %d messages\n", moreMsg.TotalCount)
	fmt.Println("-----------------------------------------------------------")

	for _, message := range moreMsg.Messages {
		switch content := message.Content.(type) {
		case *client.MessageText:
			fmt.Println(content.Text.Text)
			fmt.Println("-----------------------------------------------------------")
		}
	}
}
