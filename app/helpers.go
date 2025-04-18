package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	tdlib "github.com/zelenin/go-tdlib/client"
)

func processMessages(msg *tdlib.Message, from string) string {
	var result string

	dt := parseDate(msg.Date)
	switch content := msg.Content.(type) {
	case *tdlib.MessageText:
		result = fmt.Sprintf("[%s] %s: %s", dt, from, content.Text.Text)
	case *tdlib.MessagePhoto, *tdlib.MessageVideo, *tdlib.MessageAudio:
		result = fmt.Sprintf("%s: [media content]", from)
	}

	return result
}

func getMessagesIds(messages []*tdlib.Message) []int64 {
	if messages == nil {
		return nil
	}

	ids := make([]int64, len(messages))

	for idx, msg := range messages {
		ids[idx] = msg.Id
	}

	return ids
}

func changeView(model tea.Model, newView viewState) (tea.Model, tea.Cmd) {
	return model, tea.Cmd(func() tea.Msg {
		return changeStateMsg{newState: newView}
	})
}

func parseDate(date int32) string {
	tm := time.Unix(int64(date), 0)
	return fmt.Sprint(tm.Format("2006-01-02 15:04:05"))
}
