package app

import (
	"context"
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
	topicsView
)

const (
	historyLength int32 = 50
	chatLength    int32 = 50
	maxScreenChat int   = 150
	loadMessages  int32 = 20
)

const (
	myIdentifier      string = "You"
	unknownIdentifier string = "Unknown"
)

var (
	margStyle         = lipgloss.NewStyle().Margin(1, 2)
	inputStyle        = lipgloss.NewStyle().Background(lipgloss.Color("253")).Foreground(lipgloss.Color("232")).PaddingLeft(2).PaddingRight(2)
	listSelectedStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("#b39ddb")).Foreground(lipgloss.Color("#b39ddb")).Padding(0, 0, 0, 1)
	meStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#b39ddb"))
	senderStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffaf00"))
	unkSenderStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Margin(1, 0)
)

var senders = make(map[int64]string)

type errMsg error
type chatListMsg []list.Item
type chatHistoryMsg []string
type tdMessageMsg string
type topicListMsg []*tdlib.ForumTopic
type openChatMsg struct {
	chatId   int64
	threadId int64
}
type changeStateMsg struct {
	newState viewState
}

type chatItem struct {
	title      string
	id         int64
	haveTopics bool
}

type topicItem struct {
	chatId   int64
	threadId int64
	title    string
}

func (c chatItem) Title() string       { return c.title }
func (c chatItem) Description() string { return fmt.Sprintf("ID: %d", c.id) }
func (c chatItem) FilterValue() string { return c.title }

func (c topicItem) Title() string       { return c.title }
func (c topicItem) Description() string { return fmt.Sprintf("ID: %d", c.chatId) }
func (c topicItem) FilterValue() string { return c.title }

type rootModel struct {
	conn     *connection.Connection
	state    viewState
	err      error
	chatList list.Model
	chat     chatModel
	topics   topicsModel
}

type chatModel struct {
	conn            *connection.Connection
	viewport        viewport.Model
	messages        []string
	chatId          int64
	threadId        int64
	input           string
	err             errMsg
	atTop           bool
	chatLoadSize    int32
	newChatLoadSize int32
	init            bool
	msgChan         chan tdMessageMsg
	ctxForMsg       context.Context
	cancelMsgChan   context.CancelFunc
}

type topicsModel struct {
	conn             *connection.Connection
	topicList        list.Model
	superGroupChatId int64
}
