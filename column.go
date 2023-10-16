package main

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Column struct {
	startTime time.Duration
	timer     timer.Model
	rotation  Rotation
	focus     bool
	width     int
	height    int
}

func (c *Column) Focus() {
	c.focus = true
}

func (c *Column) Blur() {
	c.focus = false
}

func (c *Column) IsFocused() bool {
	return c.focus
}

func NewColumn(rotation Rotation, st time.Duration) Column {
	var focus bool
	if rotation == pomodoro {
		focus = true
	}
	return Column{
		startTime: st,
		timer:     timer.NewWithInterval(st, time.Second),
		rotation:  rotation,
		focus:     focus,
	}
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		c.timer, cmd = c.timer.Update(msg)
		return c, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		c.timer, cmd = c.timer.Update(msg)
		Keys.stop.SetEnabled(c.timer.Running())
		Keys.start.SetEnabled(!c.timer.Running())
		return c, cmd
	case timer.TimeoutMsg:
		// switch to next
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.reset):
			c.timer.Timeout = c.startTime
		case key.Matches(msg, Keys.start, Keys.stop):
			return c, c.timer.Toggle()
		}
	}
	c.timer, cmd = c.timer.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.GetStyle().Render(c.timer.View())
}

func (c *Column) GetStyle() lipgloss.Style {
	if c.focus == true {
		return lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
	} else {
		return lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.HiddenBorder())
	}
}

func (c *Column) SetSize(width int) {
	c.width = width / margin
}
