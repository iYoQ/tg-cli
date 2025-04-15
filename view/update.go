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

func (m model) openChatCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := m.conn.Client.OpenChat(context.Background(), &tdlib.OpenChatRequest{ChatId: m.chatId})
		if err != nil {
			return errMsg(err)
		}

		history, err := getChatHistory(m.conn, m.chatId)
		if err != nil {
			return errMsg(err)
		}

		return chatHistoryMsg(history)
	}
}

func (m model) listenUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		for msg := range m.conn.UpdatesChannel {
			if msg.ChatId == m.chatId {
				from := getUserName(m.conn, msg)
				formatMsg := processMessages(msg, from)
				updateMsg := tdMessageMsg(formatMsg)

				if err := readMessage(m.conn.Client, msg.ChatId, msg.Id); err != nil {
					return errMsg(err)
				}

				return updateMsg
			}
		}
		return nil
	}
}

func closeChat(client *tdlib.Client, chatId int64) error {
	_, err := client.CloseChat(context.Background(), &tdlib.CloseChatRequest{ChatId: chatId})
	if err != nil {
		return err
	}

	return nil
}

func readMessage(client *tdlib.Client, chatId int64, messageId int64) error {
	var messageIds []int64
	messageIds = append(messageIds, messageId)
	_, err := client.ViewMessages(context.Background(), &tdlib.ViewMessagesRequest{
		ChatId:     chatId,
		MessageIds: messageIds,
	})
	if err != nil {
		return err
	}
	return nil
}
