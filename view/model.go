package view

import (
	"context"
	"fmt"
	"tg-cli/connection"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tdlib "github.com/zelenin/go-tdlib/client"
)

type viewState int

const (
	chatListView viewState = iota
	chatView
)

const (
	historyLength int32 = 50
	chatLength    int32 = 50
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	senders  = make(map[int64]string)
)

type errMsg error
type chatListMsg []*tdlib.Chat
type chatHistoryMsg []string
type tdMessageMsg string

type chatItem struct {
	title string
	id    int64
}

func (c chatItem) Title() string       { return c.title }
func (c chatItem) Description() string { return fmt.Sprintf("ID: %d", c.id) }
func (c chatItem) FilterValue() string { return c.title }

type model struct {
	conn     *connection.Connection
	state    viewState
	chatId   int64
	messages []string
	input    string
	err      error
	chatList list.Model
}

func NewModel(conn *connection.Connection) model {
	chatList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	chatList.SetStatusBarItemName("chat", "chats")
	chatList.SetShowTitle(false)

	return model{
		conn:     conn,
		state:    chatListView,
		chatList: chatList,
	}
}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		chats, err := m.conn.Client.GetChats(context.Background(), &tdlib.GetChatsRequest{Limit: chatLength})
		if err != nil {
			return errMsg(err)
		}

		var results []*tdlib.Chat

		for _, id := range chats.ChatIds {
			chat, err := m.conn.Client.GetChat(context.Background(), &tdlib.GetChatRequest{ChatId: id})
			if err == nil {
				results = append(results, chat)
			}
			senders[id] = chat.Title
		}
		return chatListMsg(results)
	}
}
