package view

import (
	"context"
	"fmt"
	"tg-cli/exchange"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
						return m, tea.Batch(m.openChatCmd(), m.listenUpdatesCmd())
					}
				}
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case chatView:
			switch msg.Type {
			case tea.KeyEnter:
				go exchange.SendText(m.conn.Client, m.chatId, m.input)
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
				return m, nil
			default:
				m.input += msg.String()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.chatList.SetSize(msg.Width, msg.Height)

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

func (m model) openChatCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := m.conn.Client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: m.chatId})
		if err != nil {
			return errMsg(err)
		}

		history, err := m.conn.Client.GetChatHistory(context.Background(), &tdlib.GetChatHistoryRequest{
			ChatId:        m.chatId,
			FromMessageId: 0,
			Offset:        0,
			Limit:         m.historyLength,
		})
		if err != nil {
			return errMsg(err)
		}

		var messages []string
		for i := len(history.Messages) - 1; i >= 0; i-- {
			msg := history.Messages[i]
			if text, ok := msg.Content.(*tdlib.MessageText); ok {
				from := fmt.Sprintf("[%s]", senders[msg.SenderId.(*tdlib.MessageSenderUser).UserId])
				messages = append(messages, fmt.Sprintf("%s %s", from, text.Text.Text))
			}
		}
		return chatHistoryMsg(messages)
	}
}

func (m model) listenUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		for msg := range m.conn.UpdatesChannel {
			if msg.ChatId == m.chatId {
				if content, ok := msg.Content.(*tdlib.MessageText); ok {
					from := fmt.Sprintf("[%s]", senders[msg.SenderId.(*tdlib.MessageSenderUser).UserId])
					return tdMessageMsg(fmt.Sprintf("%s %s", from, content.Text.Text))
				}
			}
		}
		return nil
	}
}
