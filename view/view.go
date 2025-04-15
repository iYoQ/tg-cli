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
		return docStyle.Render(m.chatList.View())

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
