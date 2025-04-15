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
type chatListMsg []*tdlib.Chat
type chatHistoryMsg []string
type tdMessageMsg string

type model struct {
	conn          *connection.Connection
	state         viewState
	chats         []*tdlib.Chat
	selected      int
	chatId        int64
	messages      []string
	input         string
	err           error
	historyLength int32
}

func NewModel(conn *connection.Connection) model {
	return model{
		conn:          conn,
		state:         chatListView,
		historyLength: 10,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		chats, err := m.conn.Client.GetChats(context.Background(), &tdlib.GetChatsRequest{Limit: m.historyLength})
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
		return chatListMsg(results)
	}
}
