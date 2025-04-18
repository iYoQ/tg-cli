package app

import (
	"fmt"
	"strings"
	"tg-cli/connection"
	"tg-cli/requests"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func newChatModel(width int, height int, chatId int64, conn *connection.Connection) chatModel {
	vp := viewport.New(width, height)
	vp.SetContent("")

	return chatModel{
		viewport: vp,
		chatId:   chatId,
		conn:     conn,
	}
}

func (m chatModel) Init() tea.Cmd {
	return tea.Batch(m.openChatCmd(), m.listenUpdatesCmd())
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			go requests.SendText(m.conn.Client, m.chatId, m.input)
			message := formatMessage(m.input, senders[m.conn.GetMe().Id], int32(time.Now().Unix()))
			m.messages = append(m.messages, message)
			m.input = ""
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			closeChat(m.conn.Client, m.chatId)
			return changeView(m, chatListView)
		case tea.KeyRunes:
			m.input += msg.String()
		case tea.KeyCtrlDown:
			m.input += "\n"
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

	case tdMessageMsg:
		m.messages = append(m.messages, string(msg))
		cmds = append(cmds, m.listenUpdatesCmd())

	case chatHistoryMsg:
		m.messages = msg
	}

	var cmd tea.Cmd
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m chatModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	wrappedInput := wrapMessage(m.input)

	newStr := fmt.Sprintf("%s\n%s", m.viewport.View(), inputStyle.Render("> "+wrappedInput))

	return chatStyle.Render(newStr)
}

func (m chatModel) listenUpdatesCmd() tea.Cmd {
	return func() tea.Msg {
		for msg := range m.conn.UpdatesChannel {
			if msg.ChatId == m.chatId {
				if sender, ok := msg.SenderId.(*tdlib.MessageSenderUser); ok {
					if sender.UserId == m.conn.GetMe().Id {
						continue
					}
				}
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

func (m chatModel) openChatCmd() tea.Cmd {
	return func() tea.Msg {
		history, err := getChatHistory(m.conn.Client, m.chatId)
		if err != nil {
			return errMsg(err)
		}

		return chatHistoryMsg(history)
	}
}

func (m chatModel) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {

		wrappedMessage := wrapMessage(msg)
		b.WriteString(wrappedMessage + "\n")
	}

	return b.String()
}
