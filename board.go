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
	tabs     []Tab
	quitting bool
	interval int
}

func NewBoard() Board {
	help := help.New()
	help.ShowAll = true
	return Board{
		help:     help,
		focused:  pomodoro,
		interval: 0,
	}
}

func (b *Board) InitTimers() {
	b.cols = []Column{
		NewColumn(pomodoro, time.Second*25, "Pomodoro"),
		NewColumn(short_break, time.Second*5, "Short Break"),
		NewColumn(long_break, time.Second*15, "Long Break"),
	}

	b.tabs = []Tab{
		NewTab("Pomodoro", true),
		NewTab("Short Break", false),
		NewTab("Long Break", false),
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

	tabs := lipgloss.JoinHorizontal(
		lipgloss.Left,
		b.tabs[pomodoro].View(),
		b.tabs[short_break].View(),
		b.tabs[long_break].View(),
	)

	board := lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("202")).Render("GOMODORO!!! // ポモドーロ!!!"),
		lipgloss.JoinVertical(
			lipgloss.Center,
			tabs,
			b.cols[b.focused].View(),
		),
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
	case TimeOut:
		b.interval++
		b.tabs[b.focused].Blur()
		if b.interval == 6 {
			b.focused = b.focused.GetNext()
			b.interval = 0
		} else if b.interval == 5 {
			b.focused = long_break
		} else if b.focused == pomodoro {
			b.focused = b.focused.GetNext()
		} else {
			b.focused = b.focused.GetPrev()
		}
		b.tabs[b.focused].Focus()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.quit):
			b.quitting = true
			return b, tea.Quit
		case key.Matches(msg, Keys.right):
			b.tabs[b.focused].Blur()
			b.focused = b.focused.GetNext()
			b.tabs[b.focused].Focus()
		case key.Matches(msg, Keys.left):
			b.tabs[b.focused].Blur()
			b.focused = b.focused.GetPrev()
			b.tabs[b.focused].Focus()
		case key.Matches(msg, Keys.help):
			Keys.FullHelp()
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
