package requests

import (
	"testing"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func TestBuildRequest(t *testing.T) {
	content := &tdlib.InputMessageText{Text: &tdlib.FormattedText{Text: "hi"}}

	req := buildRequest(content, Params{ChatId: 7, ThreadId: 11})

	if req.ChatId != 7 {
		t.Errorf("ChatId = %d, want 7", req.ChatId)
	}
	if req.MessageThreadId != 11 {
		t.Errorf("MessageThreadId = %d, want 11", req.MessageThreadId)
	}
	if req.InputMessageContent != tdlib.InputMessageContent(content) {
		t.Errorf("InputMessageContent = %v, want %v", req.InputMessageContent, content)
	}
}
