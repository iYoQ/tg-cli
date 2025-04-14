package view

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	switch m.state {
	case chatListView:
		var b strings.Builder
		b.WriteString("Select a chat:\n\n")
		for i, c := range m.chats {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}
			fmt.Fprintf(&b, "%s %s (%d)\n", cursor, c.Title, c.Id)
		}
		return b.String()

	case chatView:
		var b strings.Builder

		for _, msg := range m.messages {
			b.WriteString(msg + "\n")
		}

		b.WriteString("\n> " + m.input)
		return b.String()
	}

	return "Loading..."

}
