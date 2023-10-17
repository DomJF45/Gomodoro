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
	title     string
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

func NewColumn(rotation Rotation, st time.Duration, title string) Column {
	var focus bool
	if rotation == pomodoro {
		focus = true
	}
	return Column{
		startTime: st,
		timer:     timer.NewWithInterval(st, time.Second),
		rotation:  rotation,
		focus:     focus,
		title:     title,
	}
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.SetSize(msg.Width, 20)
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
	return c.GetStyle().Render("Time Remaining: ", lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04b575")).Render(c.timer.View()))
}

func (c *Column) GetTextStyle() lipgloss.Style {
	if c.focus {
		return lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).MarginBottom(10)
	} else {
		return lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).MarginBottom(10)
	}
}

func (c *Column) GetStyle() lipgloss.Style {
	return lipgloss.NewStyle().Width(c.width).Height(c.height).Align(lipgloss.Top, lipgloss.Left).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("202"))
}

func (c *Column) SetSize(width, height int) {
	c.width = width / margin
	c.height = height / margin
}
