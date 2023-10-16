package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Rotation int

const timeout = time.Minute * 25

const (
	Unkown Rotation = iota
	pomodoro
	short_break
	long_break
)

const margin = 5

var (
	timerStyle = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.HiddenBorder())

	timerFocused = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
)

func (r Rotation) GetNext() {
	if r == long_break {
		r = pomodoro
	} else {
		r++
	}
}

func (r Rotation) GetPrev() {
	if r == pomodoro {
		r = long_break
	} else {
		r--
	}
}

type Model struct {
	timer    timer.Model
	keymap   KeyMap
	help     help.Model
	quitting bool
	rotation Rotation
	index    int
}

func (m Model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.reset):
			m.timer.Timeout = timeout
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m Model) HelpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func (m Model) View() string {
	var s string

	model := m.CurrentFocusedModel()
	switch m.rotation {
	case pomodoro:
		s += lipgloss.JoinHorizontal()
	case short_break:
		s += lipgloss.JoinHorizontal()
	case long_break:
		s += lipgloss.JoinHorizontal()
	}
	if m.rotation == pomodoro {
	}

	return s
}

func (m Model) CurrentFocusedModel() string {
	switch m.rotation {
	case pomodoro:
		return "pomodoro"
	case short_break:
		return "short break"
	default:
		return "long break"
	}
}

func InitialModel() Model {
	return Model{
		timer:  timer.NewWithInterval(timeout, time.Second),
		keymap: Keys,
		help:   help.New(),
	}
}

func main() {
	m := InitialModel()
	m.keymap.start.SetEnabled(false)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
