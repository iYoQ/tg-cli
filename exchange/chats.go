package exchange

import (
	"context"
	"fmt"
	"log"
	"strings"
	"tg-cli/connection"

	"github.com/chzyer/readline"
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
func OpenChat(conn *connection.Connection, chatId int64, reader *readline.Instance) {
	_, err := conn.Client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: chatId})
	if err != nil {
		log.Printf("Failed open chat %d, error: %s", chatId, err)
		return
	}

	messages, err := conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: 0,
		Offset:        0,
		Limit:         1,
	})
	if err != nil {
		log.Printf("Cannot receive last message, error: %s", err)
		return
	}

	moreMsg, err := conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
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

	// сделать в reverse порядке?
	for _, message := range moreMsg.Messages {
		switch content := message.Content.(type) {
		case *tdlib.MessageText:
			fmt.Println(content.Text.Text)
			fmt.Println("-----------------------------------------------------------")
		}
	}

	inputChannel := make(chan string)

	go func() {
		for {
			msg, err := reader.Readline()
			if err != nil {
				if err.Error() != "Interrupt" {
					fmt.Printf("Failed to read input, error: %#v\n", err)
				}
				msg = "exit"
			}

			inputChannel <- msg
			if msg == "exit" {
				close(inputChannel)
				return
			}
		}
	}()

	for {
		select {
		case message, ok := <-conn.UpdatesChannel:
			if !ok {
				inputChannel <- "Channel is closed"
				return
			}

			if message.ChatId == chatId {
				switch content := message.Content.(type) {
				case *tdlib.MessageText:
					fmt.Printf("%s\n", content.Text.Text)
				}
			}
		case msg := <-inputChannel:
			if msg == "exit" {
				return
			}

			msgSplit := strings.Split(msg, " ")
			if msgSplit[0] == "/ph" {
				photoPath := msgSplit[1]
				text := strings.Join(msgSplit[2:], " ")
				SendPhoto(conn.Client, chatId, photoPath, text)
			} else {
				SendText(conn.Client, chatId, msg)
			}
		}

	}

}
