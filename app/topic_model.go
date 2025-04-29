package app

import (
	"context"
	"tg-cli/connection"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func newTopicsModel(width int, height int, chatId int64, conn *connection.Connection) topicsModel {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = listSelectedStyle
	delegate.Styles.SelectedDesc = listSelectedStyle

	topicsList := list.New([]list.Item{}, delegate, width, height)
	topicsList.SetStatusBarItemName("topic", "topics")
	topicsList.SetShowTitle(false)

	return topicsModel{
		conn:             conn,
		topicList:        topicsList,
		superGroupChatId: chatId,
	}
}

func (m topicsModel) Init() tea.Cmd {
	return func() tea.Msg {
		result, _ := m.conn.Client.GetForumTopics(context.Background(), &tdlib.GetForumTopicsRequest{ChatId: m.superGroupChatId, Limit: 100})

		return topicListMsg(result.Topics)
	}
}

func (m topicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			item := m.topicList.SelectedItem().(topicItem)
			return m, func() tea.Msg { return openChatMsg{chatId: item.chatId, threadId: item.threadId} }

		case "ctrl+c", "q":
			return changeView(m, chatListView)
		}

	case topicListMsg:
		var item topicItem
		items := make([]list.Item, 0, len(msg))
		for _, topic := range msg {
			item = topicItem{chatId: m.superGroupChatId, threadId: topic.Info.MessageThreadId, title: topic.Info.Name}

			items = append(items, item)
		}
		m.topicList.SetItems(items)

	case tea.WindowSizeMsg:
		h, v := margStyle.GetFrameSize()
		m.topicList.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.topicList, cmd = m.topicList.Update(msg)
	return m, cmd
}

func (m topicsModel) View() string {
	return m.topicList.View()
}

func (m topicsModel) openTopic(topic topicItem) (tea.Model, tea.Cmd) {
	return m, nil
}
