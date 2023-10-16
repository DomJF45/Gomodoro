package main

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Board struct {
	help     help.Model
	loaded   bool
	focused  Rotation
	cols     []Column
	quitting bool
}

func NewBoard() Board {
	help := help.New()
	help.ShowAll = true
	return Board{
		help:    help,
		focused: pomodoro,
	}
}

func (b Board) Init() tea.Cmd {
	return nil
}

func (b Board) View() string {
	if b.quitting {
		return ""
	}

	if !b.loaded {
		return "loading..."
	}

	board := lipgloss.JoinHorizontal(
		lipgloss.Center,
		b.cols[pomodoro].View(),
		b.cols[short_break].View(),
		b.cols[long_break].View(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, board, b.help.View(Keys))
}

func (b Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
}
