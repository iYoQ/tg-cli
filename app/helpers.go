package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tg-cli/connection"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	tdlib "github.com/zelenin/go-tdlib/client"
	"golang.org/x/term"
)

func processMessages(msg *tdlib.Message, from string) string {
	var result string

	switch content := msg.Content.(type) {
	case *tdlib.MessageText:
		result = formatMessage(content.Text.Text, from, msg.Date)
	case *tdlib.MessagePhoto:
		var tmpText string
		if content.Caption != nil {
			tmpText = fmt.Sprintf("[media content] %s", content.Caption.Text)
		}
		result = formatMessage(tmpText, from, msg.Date)
	case *tdlib.MessageDocument:
		var tmpText string
		if content.Caption != nil {
			tmpText = fmt.Sprintf("[media content] %s", content.Caption.Text)
		}
		result = formatMessage(tmpText, from, msg.Date)
	case *tdlib.MessageAnimatedEmoji:
		result = formatMessage(content.Emoji, from, msg.Date)
	case *tdlib.MessageVideo, *tdlib.MessageAudio:
		result = formatMessage("[media content]", from, msg.Date)
	default:
		result = formatMessage("[smh]", from, msg.Date)
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

func formatMessage(msg string, from string, unixDate int32) string {
	dt := parseDate(unixDate)
	msg = addIndenting(msg, dt)

	return fmt.Sprintf("[%s] %s: %s", dt, from, msg)
}

func formatCommand(msg string, cmdType string) (string, string, error) {
	msgTrim := strings.TrimSpace(msg)
	msgSplit := strings.Split(msgTrim, " ")

	if cmdType == "photo" || cmdType == "file" {
		path := msgSplit[1]
		text := strings.Join(msgSplit[2:], " ")
		return path, text, nil
	} else {
		return "", "", errors.New("error in formatCommand")
	}
}

func parseDate(date int32) string {
	tm := time.Unix(int64(date), 0)
	now := time.Now()

	if tm.Year() == now.Year() && tm.YearDay() == now.YearDay() {
		return tm.Format("15:04:05")
	}

	if tm.Year() == now.Year() {
		return tm.Format("01/02 15:04:05")
	}

	return tm.Format("2006-01-02 15:04:05")
}

func wrapMessage(msg string) string {
	termWidth, _, err := term.GetSize(0)
	if err != nil {
		termWidth = maxScreenChat
	}
	formattedMessage := wordwrap.String(msg, min(termWidth, maxScreenChat))
	return formattedMessage
}

func addIndenting(msg string, str string) string {
	var result strings.Builder
	indentation := len(str) + 3
	for idx := 0; idx < len(msg); idx++ {
		if msg[idx] == '\n' && idx+1 < len(msg) && msg[idx+1] != '\n' {
			result.WriteString("\n")
			result.WriteString(strings.Repeat(" ", indentation))
		} else {
			result.WriteByte(msg[idx])
		}
	}
	return result.String()
}

func checkCommand(msg string) string {
	msgSplit := strings.Split(msg, " ")
	if msgSplit[0] == "/p" && len(msgSplit) > 1 {
		return "photo"
	} else if msgSplit[0] == "/f" && len(msgSplit) > 1 {
		return "file"
	} else {
		return ""
	}
}

func listenFromUpdateChannel(conn *connection.Connection, currentChatId int64, msgChannel chan<- tdMessageMsg, ctx context.Context) {
	for {
		select {
		case msg := <-conn.UpdatesChannel:
			if msg.ChatId != currentChatId {
				continue
			}
			if msg.IsOutgoing && currentChatId != conn.GetMe().Id {
				continue
			}

			from := getUserName(conn.Client, msg)
			formatMsg := processMessages(msg, from)

			readMessages(conn.Client, msg.ChatId, []int64{msg.Id})

			msgChannel <- tdMessageMsg(formatMsg)

		case <-ctx.Done():
			return
		}
	}
}
