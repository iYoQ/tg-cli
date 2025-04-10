package console

import (
	"fmt"
	"strconv"
	"strings"

	"tg-cli/exchange"

	"github.com/chzyer/readline"
	tdlib "github.com/zelenin/go-tdlib/client"
)

const NUMBER_OF_CHATS = 5

func Start(client *tdlib.Client, reader *readline.Instance, updatesChannel chan *tdlib.Message) {
	fmt.Printf("\n%d recently open chats:\n", NUMBER_OF_CHATS)
	fmt.Println("-----------------------------------------------------------")
	exchange.GetChats(client, NUMBER_OF_CHATS)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. send msg")
		fmt.Println("2. get chat list")
		fmt.Println("3. open chat")
		fmt.Println("9. exit")

		input, err := readInput(reader)
		if err != nil {
			continue
		}

		choice, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			fmt.Println("invalid input, enter a number")
			continue
		}

		switch choice {
		case 1:
			createMessage(client, reader)
		case 2:
			getChatList(client, reader)
		case 3:
			openChat(client, updatesChannel, reader)
		case 9:
			return
		default:
			fmt.Println("invalid")
		}
	}
}

func createMessage(client *tdlib.Client, reader *readline.Instance) {
	fmt.Println("\nEnter chat id:")

	input, err := readInput(reader)
	if err != nil {
		return
	}

	chatId, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println("Invalid id, enter a number")
		return
	}

	fmt.Println("\nChoose a msg:")

	msg, err := readInput(reader)
	if err != nil {
		return
	}

	msgSplit := strings.Split(msg, " ")
	if msgSplit[0] == "/ph" {
		photoPath := msgSplit[1]
		text := strings.Join(msgSplit[2:], " ")
		exchange.SendPhoto(client, chatId, photoPath, text)
	} else {
		exchange.SendText(client, chatId, msg)
	}
	fmt.Println("Message sent")
}

func getChatList(client *tdlib.Client, reader *readline.Instance) {
	fmt.Println("\nChoose a number of chats:")

	input, err := readInput(reader)
	if err != nil {
		return
	}

	size64, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		fmt.Println("Invalid size, enter a number")
		return
	}

	size32 := int32(size64)

	exchange.GetChats(client, size32)
}

func openChat(client *tdlib.Client, updatesChannel chan *tdlib.Message, reader *readline.Instance) {
	fmt.Println("\nEnter chat id:")

	input, err := readInput(reader)
	if err != nil {
		return
	}

	chatId, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		fmt.Println("Invalid id, enter a number")
		return
	}

	exchange.OpenChat(client, chatId, updatesChannel, reader)
}

func readInput(reader *readline.Instance) (string, error) {
	input, err := reader.Readline()
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}

	input = strings.TrimSpace(input)
	return input, nil
}
