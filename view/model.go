package view

import (
	"context"
	"tg-cli/connection"

	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

type viewState int

const (
	chatListView viewState = iota
	chatView
)

type errMsg error
type chatList []*tdlib.Chat
type chatHistoryMsg []string
type tdMessageMsg string

const LEN = 10

type model struct {
	conn     *connection.Connection
	state    viewState
	chats    []*tdlib.Chat
	selected int
	chatId   int64
	messages []string
	input    string
	err      error
}

func NewModel(conn *connection.Connection) model {
	return model{
		conn:  conn,
		state: chatListView,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		chats, err := m.conn.Client.GetChats(context.Background(), &tdlib.GetChatsRequest{Limit: LEN})
		if err != nil {
			return errMsg(err)
		}

		var results []*tdlib.Chat

		for _, id := range chats.ChatIds {
			chat, err := m.conn.Client.GetChat(context.Background(), &tdlib.GetChatRequest{ChatId: id})
			if err == nil {
				results = append(results, chat)
			}
		}
		return chatList(results)
	}
}
