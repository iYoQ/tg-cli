package app

import (
	"context"
	"fmt"
	"strings"
	"tg-cli/connection"
	"tg-cli/requests"

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
	pageSize      int32 = 50
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	senders  = make(map[int64]string)
)

type errMsg error
type chatListMsg []*tdlib.Chat
type chatHistoryMsg []string
type tdMessageMsg string
type changeStateMsg struct {
	newState viewState
}

type chatItem struct {
	title string
	id    int64
}

func (c chatItem) Title() string       { return c.title }
func (c chatItem) Description() string { return fmt.Sprintf("ID: %d", c.id) }
func (c chatItem) FilterValue() string { return c.title }

type rootModel struct {
	conn     *connection.Connection
	state    viewState
	chatId   int64
	messages []string
	input    string
	err      error
	chatList list.Model
	chat     chatModel
}

func NewRootModel(conn *connection.Connection) rootModel {
	chatList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	chatList.SetStatusBarItemName("chat", "chats")
	chatList.SetShowTitle(false)

	return rootModel{
		conn:     conn,
		state:    chatListView,
		chatList: chatList,
	}
}

func (m rootModel) Init() tea.Cmd {
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

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case chatListView:
			switch msg.String() {
			case "enter":
				if m.chatList.FilterState() != list.Filtering {
					item, ok := m.chatList.SelectedItem().(chatItem)
					if ok {
						m.chatId = item.id
						m.state = chatView
						m.chat = newChatModel(m.chatList.Width(), m.chatList.Height())
						chatCmd := m.chat.Init()
						return m, chatCmd
						// return m, tea.Batch(m.openChatCmd(), m.listenUpdatesCmd())
					}
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case chatView:
			switch msg.Type {
			case tea.KeyEnter:
				go requests.SendText(m.conn.Client, m.chatId, m.input)
				if m.chatId != m.conn.GetMe().Id {
					m.messages = append(m.messages, fmt.Sprintf("You: %s", m.input))
				}
				m.input = ""
			case tea.KeyBackspace:
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case tea.KeyCtrlC:
				m.state = chatListView
				closeChat(m.conn.Client, m.chatId)
				return m, nil
			case tea.KeyRunes:
				m.input += msg.String()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.chatList.SetSize(msg.Width-h, msg.Height-v)

	case errMsg:
		m.err = msg

	case chatListMsg:
		items := make([]list.Item, 0, len(msg))
		for _, chat := range msg {
			items = append(items, chatItem{title: chat.Title, id: chat.Id})
		}
		m.chatList.SetItems(items)

	case chatHistoryMsg:
		m.messages = msg

	case tdMessageMsg:
		m.messages = append(m.messages, string(msg))
		return m, m.listenUpdatesCmd()
	}

	var cmd tea.Cmd
	m.chatList, cmd = m.chatList.Update(msg)
	if m.chatList.FilterState() == list.Filtering {
		return m, cmd
	}

	return m, nil
}

func (m rootModel) openChatCmd() tea.Cmd {
	return func() tea.Msg {
		history, err := getChatHistory(m.conn.Client, m.chatId)
		if err != nil {
			return errMsg(err)
		}

		return chatHistoryMsg(history)
	}
}

func (m rootModel) listenUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		for msg := range m.conn.UpdatesChannel {
			if msg.ChatId == m.chatId {
				from := getUserName(m.conn.Client, msg)
				formatMsg := processMessages(msg, from)
				updateMsg := tdMessageMsg(formatMsg)

				messageIds := make([]int64, 1)
				messageIds[0] = msg.Id

				if err := readMessages(m.conn.Client, msg.ChatId, messageIds); err != nil {
					return errMsg(err)
				}

				return updateMsg
			}
		}
		return nil
	}
}

func (m rootModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	switch m.state {
	case chatListView:
		return docStyle.Render(m.chatList.View())

	case chatView:
		var b strings.Builder

		for _, msg := range m.messages {
			b.WriteString(msg + "\n")
		}

		b.WriteString("\n> " + m.input)
		return b.String()
	}

	return "Loading..."

}
