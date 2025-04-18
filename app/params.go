package app

import (
	"fmt"
	"tg-cli/connection"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
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

const (
	myIdentifier      string = "You"
	unknownIdentifier string = "Unknown"
)

var (
	margStyle = lipgloss.NewStyle().Margin(1, 2)
	// chatStyle         = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#b39ddb"))
	inputStyle        = lipgloss.NewStyle().Background(lipgloss.Color("253")).Foreground(lipgloss.Color("232")).PaddingLeft(2).PaddingRight(2)
	listSelectedStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("#b39ddb")).Foreground(lipgloss.Color("#b39ddb")).Padding(0, 0, 0, 1)
	meStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#b39ddb"))
	senderStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaf00"))
	unkSenderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
)

var senders = make(map[int64]string)

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
	err      error
	chatList list.Model
	chat     chatModel
}

type chatModel struct {
	viewport viewport.Model
	messages []string
	chatId   int64
	conn     *connection.Connection
	input    string
	err      errMsg
}
