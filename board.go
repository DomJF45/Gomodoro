package main

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

func (b *Board) InitTimers() {
	b.cols = []Column{
		NewColumn(pomodoro, time.Minute*25),
		NewColumn(short_break, time.Minute*5),
		NewColumn(long_break, time.Minute*15),
	}
	b.cols[pomodoro].title = "Pomodoro"
	b.cols[short_break].title = "Short Break"
	b.cols[long_break].title = "Long Break"
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		var cmds []tea.Cmd
		b.help.Width = msg.Width - margin
		for i := 0; i < len(b.cols); i++ {
			var res tea.Model
			res, cmd = b.cols[i].Update(msg)
			b.cols[i] = res.(Column)
			cmds = append(cmds, cmd)
		}
		b.loaded = true
		return b, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.quit):
			b.quitting = true
			return b, tea.Quit
		case key.Matches(msg, Keys.right):
			b.cols[b.focused].Blur()
			b.focused = b.focused.GetNext()
			b.cols[b.focused].Focus()
		case key.Matches(msg, Keys.left):
			b.cols[b.focused].Blur()
			b.focused = b.focused.GetPrev()
			b.cols[b.focused].Focus()
		}
	}
	res, cmd := b.cols[b.focused].Update(msg)
	if _, ok := res.(Column); ok {
		b.cols[b.focused] = res.(Column)
	} else {
		return res, cmd
	}
	return b, cmd
}
