package exchange

import (
	"context"
	"fmt"
	"log"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func GetChats(client *tdlib.Client, size int32) {
	chats, err := client.GetChats(context.Background(), &tdlib.GetChatsRequest{Limit: size})
	if err != nil {
		log.Printf("Failed GetChats: %s", err)
		return
	}

	for _, id := range chats.ChatIds {
		chat, err := client.GetChat(context.Background(), &tdlib.GetChatRequest{ChatId: id})
		if err != nil {
			log.Printf("Failder get chat %d, error %s", id, err)
			continue
		}

		fmt.Printf("name: %s, id: %d\n", chat.Title, chat.Id)
		fmt.Println("-----------------------------------------------------------")
	}
}

// Переработать этот пиздец, добавить идентификатор того кто отправлял сообщение
func GetMessages(client *tdlib.Client, chatId int64, updatesChannel chan *tdlib.Message) {
	_, err := client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: chatId})
	if err != nil {
		log.Printf("Failed open chat %d, error: %s", chatId, err)
		return
	}

	messages, err := client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: 0,
		Offset:        0,
		Limit:         1,
	})
	if err != nil {
		log.Printf("Cannot receive last message, error: %s", err)
		return
	}

	moreMsg, err := client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: messages.Messages[0].Id,
		Offset:        0,
		Limit:         10,
	})
	if err != nil {
		log.Printf("Cannot receive messages, error: %s", err)
		return
	}

	fmt.Println("-----------------------------------------------------------")
	fmt.Printf("Chat history, last %d messages\n", moreMsg.TotalCount)
	fmt.Println("-----------------------------------------------------------")

	for _, message := range moreMsg.Messages {
		switch content := message.Content.(type) {
		case *tdlib.MessageText:
			fmt.Println(content.Text.Text)
			fmt.Println("-----------------------------------------------------------")
		}
	}

	for {
		message, ok := <-updatesChannel
		if !ok {
			break
		}

		switch content := message.Content.(type) {
		case *tdlib.MessageText:
			fmt.Printf("%s\n", content.Text.Text)
		}

	}
}
