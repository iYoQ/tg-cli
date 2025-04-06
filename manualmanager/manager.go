package manualmanager

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"tg-cli/handlers"
	"tg-cli/info"
	"tg-cli/sender"

	"github.com/zelenin/go-tdlib/client"
)

const NUMBER_OF_CHATS = 5

func Start(my_client *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\n%d recently open chats:\n", NUMBER_OF_CHATS)
	fmt.Println("-----------------------------------------------------------")
	info.GetChats(my_client, NUMBER_OF_CHATS)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. send msg")
		fmt.Println("2. get chat list")
		fmt.Println("9. exit")

		input, err := readInput(reader)
		if err != nil {
			fmt.Println("Failed to read input")
			continue
		}

		choice, err := strconv.ParseInt(input, 10, 32)
		if err != nil {
			fmt.Println("invalid input, enter a number")
			continue
		}

		switch choice {
		case 1:
			createMessage(my_client, reader)
		case 2:
			getChatList(my_client, reader)
		case 9:
			handlers.Shutdown(my_client)
		default:
			fmt.Println("invalid")
		}
	}
}

func createMessage(my_client *client.Client, reader *bufio.Reader) {
	fmt.Println("\nEnter chat id:")

	input, err := readInput(reader)
	if err != nil {
		fmt.Println("Failed to read input")
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
		fmt.Println("Failed to read input")
		return
	}

	attPath := strings.Split(msg, "=")
	if attPath[0] == "ph" {
		sender.SendPhoto(my_client, chatId, strings.Join(attPath[1:], "="))
	} else {
		sender.SendText(my_client, chatId, msg)
	}
}

func getChatList(my_client *client.Client, reader *bufio.Reader) {
	fmt.Println("\nChoose a number of chats:")

	input, err := readInput(reader)
	if err != nil {
		fmt.Println("Failed to read input")
		return
	}

	size64, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		fmt.Println("Invalid size, enter a number")
		return
	}

	size32 := int32(size64)

	info.GetChats(my_client, size32)
}

func readInput(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	return input, nil
}
