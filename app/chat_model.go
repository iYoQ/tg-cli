package app

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type chatModel struct {
	viewport viewport.Model
}

func newChatModel(width int, height int) chatModel {
	return chatModel{
		viewport: viewport.New(width-2, height-7),
	}
}

func (m chatModel) Init() tea.Cmd {
	return nil
}
