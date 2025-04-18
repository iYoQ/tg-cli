package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	tdlib "github.com/zelenin/go-tdlib/client"
	"golang.org/x/term"
)

func processMessages(msg *tdlib.Message, from string, width int) string {
	var result string

	switch content := msg.Content.(type) {
	case *tdlib.MessageText:
		result = formatMessage(content.Text.Text, from, msg.Date, width)
	case *tdlib.MessagePhoto, *tdlib.MessageVideo, *tdlib.MessageAudio:
		result = formatMessage("[media content]", from, msg.Date, width)
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

func formatMessage(msg string, from string, unixDate int32, width int) string {
	dt := parseDate(unixDate)

	return fmt.Sprintf("[%s] %s: %s", dt, from, msg)
}

func wrapMessage(msg string) string {
	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = maxScreenChat
	}
	formattedMessage := wordwrap.String(msg, min(termWidth, maxScreenChat))
	return formattedMessage
}
