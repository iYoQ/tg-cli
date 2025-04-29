package app

import (
	"fmt"
	"strings"
	"tg-cli/connection"
	"tg-cli/requests"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func newChatModel(width int, height int, chatId int64, threadId int64, conn *connection.Connection) chatModel {
	vp := viewport.New(width, height)
	vp.SetContent("")

	return chatModel{
		viewport:        vp,
		chatId:          chatId,
		conn:            conn,
		atTop:           false,
		chatLoadSize:    20,
		newChatLoadSize: 20,
		init:            true,
		threadId:        threadId,
	}
}

func (m chatModel) Init() tea.Cmd {
	return tea.Batch(m.openChatCmd(m.chatLoadSize), m.listenUpdatesCmd())
}

func (m chatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			switch checkCommand(m.input) {
			case "photo":
				path, caption, err := formatCommand(m.input, "photo")
				if err != nil {
					return m, func() tea.Msg { return errMsg(err) }
				}
				go requests.SendPhoto(m.conn.Client, m.chatId, path, caption)
				if m.chatId != m.conn.GetMe().Id {
					message := formatMessage("[media content]", senders[m.conn.GetMe().Id], int32(time.Now().Unix()))
					m.messages = append(m.messages, message)
				}
			case "file":
				path, caption, err := formatCommand(m.input, "file")
				if err != nil {
					return m, func() tea.Msg { return errMsg(err) }
				}
				go requests.SendFile(m.conn.Client, m.chatId, path, caption)
				if m.chatId != m.conn.GetMe().Id {
					message := formatMessage("[media content]", senders[m.conn.GetMe().Id], int32(time.Now().Unix()))
					m.messages = append(m.messages, message)
				}
			default:
				go requests.SendText(m.conn.Client, m.chatId, m.input)
				if m.chatId != m.conn.GetMe().Id {
					message := formatMessage(m.input, senders[m.conn.GetMe().Id], int32(time.Now().Unix()))
					m.messages = append(m.messages, message)
				}
			}

			m.input = ""
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			closeChat(m.conn.Client, m.chatId)
			return changeView(m, chatListView)
		case tea.KeyRunes, tea.KeySpace:
			m.input += msg.String()
		case tea.KeyCtrlDown:
			m.input += "\n"
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height

	case tdMessageMsg:
		m.messages = append(m.messages, string(msg))
		m.viewport.SetContent(m.renderMessages())
		m.viewport.GotoBottom()
		cmds = append(cmds, m.listenUpdatesCmd())

	case chatHistoryMsg:
		prevLineCount := m.viewport.TotalLineCount()
		prevYOffset := m.viewport.YOffset

		m.messages = msg
		m.viewport.SetContent(m.renderMessages())
		if m.init {
			m.viewport.GotoBottom()
			m.init = false
		} else {
			newLines := m.viewport.TotalLineCount() - prevLineCount
			m.viewport.YOffset = prevYOffset + newLines
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	if m.viewport.YOffset == 0 && !m.atTop {
		m.atTop = true
		m.newChatLoadSize += loadMessages
		cmds = append(cmds, m.openChatCmd(m.newChatLoadSize))
	} else if m.viewport.YOffset > 0 && m.atTop {
		m.atTop = false
	}

	return m, tea.Batch(cmds...)
}

func (m chatModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	wrappedInput := wrapMessage(m.input)

	help := "[/f] - send file, [/f] send photo, [Ctrl+C]/[Esc] return"

	newStr := fmt.Sprintf("%s\n%s\n%s", m.viewport.View(), inputStyle.Render("> "+wrappedInput), helpStyle.Render(help))

	return chatStyle.Render(newStr)
}

func (m chatModel) listenUpdatesCmd() tea.Cmd {
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

func (m chatModel) openChatCmd(chatLoadSize int32) tea.Cmd {
	chatLoadSize = max(m.chatLoadSize, chatLoadSize)
	return func() tea.Msg {
		history, err := getChatHistory(m.conn.Client, m.chatId, m.threadId, chatLoadSize)
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
