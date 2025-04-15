package view

import (
	"fmt"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func processMessages(msg *tdlib.Message, from string) string {
	var result string
	switch content := msg.Content.(type) {
	case *tdlib.MessageText:
		result = fmt.Sprintf("%s %s", from, content.Text.Text)
	case *tdlib.MessagePhoto, *tdlib.MessageVideo, *tdlib.MessageAudio:
		result = fmt.Sprintf("%s [media content]", from)
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
