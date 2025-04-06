package info

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
	}
}
