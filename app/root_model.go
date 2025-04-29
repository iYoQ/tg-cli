package app

import (
	"context"
	"fmt"
	"tg-cli/connection"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func NewRootModel(conn *connection.Connection) rootModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = listSelectedStyle
	delegate.Styles.SelectedDesc = listSelectedStyle

	chatList := list.New([]list.Item{}, delegate, 0, 0)
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
			senders[id] = senderStyle.Render(chat.Title)
		}
		senders[m.conn.GetMe().Id] = meStyle.Render(myIdentifier)
		return chatListMsg(results)
	}
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.chatList.FilterState() != list.Filtering && m.state == chatListView {
				item := m.chatList.SelectedItem().(chatItem)
				if item.haveTopics {
					return openTopics(m, item)
				}
				return openChat(m, item)
			}
		case "ctrl+c", "q":
			if m.state == chatListView {
				return m, tea.Quit
			}
		}

	case changeStateMsg:
		m.state = msg.newState

	case tea.WindowSizeMsg:
		h, v := margStyle.GetFrameSize()
		m.chatList.SetSize(msg.Width-h, msg.Height-v)

	case errMsg:
		m.err = msg

	case chatListMsg:
		var item chatItem
		items := make([]list.Item, 0, len(msg))
		for _, chat := range msg {
			if chat.ViewAsTopics {
				item = chatItem{title: chat.Title, id: chat.Id, haveTopics: true}
			} else {
				item = chatItem{title: chat.Title, id: chat.Id}
			}
			items = append(items, item)
		}
		m.chatList.SetItems(items)
	}

	var cmd tea.Cmd
	switch m.state {
	case chatListView:
		updatedModel, newCmd := m.chatList.Update(msg)
		m.chatList = updatedModel
		cmd = newCmd
	case topicsView:
		updatedModel, newCmd := m.topics.Update(msg)
		m.topics = updatedModel.(topicsModel)
		cmd = newCmd
	case chatView:
		updatedModel, newCmd := m.chat.Update(msg)
		m.chat = updatedModel.(chatModel)
		cmd = newCmd
	}

	return m, cmd
}

func (m rootModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	switch m.state {
	case chatListView:
		return margStyle.Render(m.chatList.View())
	case topicsView:
		return margStyle.Render(m.topics.View())
	case chatView:
		return m.chat.View()
	}
	return "Loading..."
}

func openChat(m rootModel, chat chatItem) (tea.Model, tea.Cmd) {
	m.state = chatView
	m.chat = newChatModel(m.chatList.Width(), m.chatList.Height(), chat.id, m.conn)
	chatCmd := m.chat.Init()
	return m, chatCmd
}

func openTopics(m rootModel, chat chatItem) (tea.Model, tea.Cmd) {
	m.state = topicsView
	m.topics = newTopicsModel(chat.id, m.conn)

	size := tea.WindowSizeMsg{
		Width:  m.chatList.Width(),
		Height: m.chatList.Height(),
	}

	chatCmd := m.topics.Init()

	return m, tea.Batch(chatCmd, func() tea.Msg { return size })
}
