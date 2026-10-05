package app

import (
	"strings"
	"testing"
	"time"

	tdlib "github.com/zelenin/go-tdlib/client"
)

func TestFormatMessage(t *testing.T) {
	date := time.Date(2020, 1, 2, 3, 4, 5, 0, time.Local)
	unixDate := int32(date.Unix())

	tests := []struct {
		name string
		msg  string
		from string
		want string
	}{
		{
			name: "simple message",
			msg:  "hello",
			from: "Alice",
			want: "[2020-01-02 03:04:05] Alice: hello",
		},
		{
			name: "multiline message",
			msg:  "line1\nline2",
			from: "Bob",
			want: "[2020-01-02 03:04:05] Bob: line1\n" + strings.Repeat(" ", len("2020-01-02 03:04:05")+3) + "line2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatMessage(tt.msg, tt.from, unixDate); got != tt.want {
				t.Errorf("formatMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAddIndenting(t *testing.T) {
	indent := strings.Repeat(" ", len("12345678")+3)

	tests := []struct {
		name string
		msg  string
		want string
	}{
		{name: "no newline", msg: "abc", want: "abc"},
		{name: "empty", msg: "", want: ""},
		{name: "single newline", msg: "a\nb", want: "a\n" + indent + "b"},
		{name: "double newline kept", msg: "a\n\nb", want: "a\n\n" + indent + "b"},
		{name: "trailing newline", msg: "a\n", want: "a\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addIndenting(tt.msg, "12345678"); got != tt.want {
				t.Errorf("addIndenting() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	now := time.Now()

	today := now.Truncate(time.Second)
	if got := parseDate(int32(today.Unix())); got != today.Format("15:04:05") {
		t.Errorf("parseDate(today) = %q, want %q", got, today.Format("15:04:05"))
	}

	prevYear := now.AddDate(-1, 0, 0)
	want := time.Unix(prevYear.Unix(), 0).Format("2006-01-02 15:04:05")
	if got := parseDate(int32(prevYear.Unix())); got != want {
		t.Errorf("parseDate(prevYear) = %q, want %q", got, want)
	}

	sameYear := time.Date(now.Year(), 1, 15, 10, 11, 12, 0, time.Local)
	if sameYear.YearDay() == now.YearDay() {
		sameYear = time.Date(now.Year(), 7, 20, 10, 11, 12, 0, time.Local)
	}
	want = sameYear.Format("01/02 15:04:05")
	if got := parseDate(int32(sameYear.Unix())); got != want {
		t.Errorf("parseDate(sameYear) = %q, want %q", got, want)
	}
}

func TestCheckCommand(t *testing.T) {
	tests := []struct {
		msg  string
		want string
	}{
		{msg: "/p /tmp/a.png", want: "photo"},
		{msg: "/p /tmp/a.png caption here", want: "photo"},
		{msg: "/f /tmp/doc.pdf", want: "file"},
		{msg: "/p", want: ""},
		{msg: "/f", want: ""},
		{msg: "hello", want: ""},
		{msg: "", want: ""},
		{msg: "  /p x", want: ""},
	}

	for _, tt := range tests {
		if got := checkCommand(tt.msg); got != tt.want {
			t.Errorf("checkCommand(%q) = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestFormatCommand(t *testing.T) {
	tests := []struct {
		name       string
		msg        string
		cmdType    string
		wantPath   string
		wantText   string
		wantErrGen bool
	}{
		{name: "photo with caption", msg: "/p /tmp/a.png nice pic", cmdType: "photo", wantPath: "/tmp/a.png", wantText: "nice pic"},
		{name: "file without caption", msg: "/f /tmp/doc.pdf", cmdType: "file", wantPath: "/tmp/doc.pdf", wantText: ""},
		{name: "unknown command type", msg: "/x path text", cmdType: "text", wantErrGen: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, text, err := formatCommand(tt.msg, tt.cmdType)
			if tt.wantErrGen {
				if err == nil {
					t.Errorf("formatCommand(%q, %q) expected error, got nil", tt.msg, tt.cmdType)
				}
				return
			}
			if err != nil {
				t.Fatalf("formatCommand(%q, %q) unexpected error: %v", tt.msg, tt.cmdType, err)
			}
			if path != tt.wantPath {
				t.Errorf("path = %q, want %q", path, tt.wantPath)
			}
			if text != tt.wantText {
				t.Errorf("text = %q, want %q", text, tt.wantText)
			}
		})
	}
}

func TestProcessMessages(t *testing.T) {
	unixDate := int32(time.Date(2020, 1, 2, 3, 4, 5, 0, time.Local).Unix())
	prefix := "[2020-01-02 03:04:05] "

	tests := []struct {
		name    string
		content tdlib.MessageContent
		want    string
	}{
		{
			name:    "text",
			content: &tdlib.MessageText{Text: &tdlib.FormattedText{Text: "hi"}},
			want:    prefix + "Alice: hi",
		},
		{
			name:    "photo with caption",
			content: &tdlib.MessagePhoto{Caption: &tdlib.FormattedText{Text: "cap"}},
			want:    prefix + "Alice: [media content] cap",
		},
		{
			name:    "photo without caption",
			content: &tdlib.MessagePhoto{},
			want:    prefix + "Alice: [media content]",
		},
		{
			name:    "animated emoji",
			content: &tdlib.MessageAnimatedEmoji{Emoji: "\U0001F525"},
			want:    prefix + "Alice: \U0001F525",
		},
		{
			name:    "video",
			content: &tdlib.MessageVideo{},
			want:    prefix + "Alice: [media content]",
		},
		{
			name:    "unsupported type",
			content: &tdlib.MessageSticker{},
			want:    prefix + "Alice: [smh]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &tdlib.Message{Id: 1, Date: unixDate, Content: tt.content}
			if got := processMessages(msg, "Alice"); got != tt.want {
				t.Errorf("processMessages() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetMessagesIds(t *testing.T) {
	if got := getMessagesIds(nil); got != nil {
		t.Errorf("getMessagesIds(nil) = %v, want nil", got)
	}

	got := getMessagesIds([]*tdlib.Message{{Id: 3}, {Id: 1}, {Id: 2}})
	want := []int64{3, 1, 2}
	if len(got) != len(want) {
		t.Fatalf("getMessagesIds() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("getMessagesIds()[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestWrapMessage(t *testing.T) {
	input := strings.TrimRight(strings.Repeat("word ", 100), " ")

	got := wrapMessage(input)
	if got == input {
		t.Fatalf("wrapMessage() did not wrap long input")
	}
	for i, line := range strings.Split(got, "\n") {
		if len(line) > maxScreenChat {
			t.Errorf("wrapMessage() line %d length = %d, want <= %d", i, len(line), maxScreenChat)
		}
	}
}
