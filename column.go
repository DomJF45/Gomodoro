package main

import (
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const debounceDuration = time.Millisecond * 300

type exitMsg struct {
	timer.Model
}

type SafeTimer struct {
	mu sync.Mutex
	t  timer.Model
}

type Column struct {
	startTime    time.Duration
	timer        timer.Model
	rotation     Rotation
	focus        bool
	width        int
	height       int
	title        string
	progress     progress.Model
	rotationFlag int
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
		startTime:    st,
		timer:        timer.NewWithInterval(st, time.Second),
		rotation:     rotation,
		focus:        focus,
		title:        title,
		progress:     progress.New(progress.WithDefaultGradient()),
		rotationFlag: 0,
	}
}

type TimeOut struct {
	bool
}

func (c Column) Init() tea.Cmd {
	return nil
}

func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var pcmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.SetSize(msg.Width, 20)
	case timer.TickMsg:
		var cmd tea.Cmd
		c.timer, cmd = c.timer.Update(msg)
		c.rotationFlag++
		if c.rotationFlag == int(c.startTime.Seconds()*0.01) {
			pcmd = c.progress.IncrPercent(.01)
			c.rotationFlag = 0
		}
		return c, tea.Batch(cmd, pcmd)
	case timer.StartStopMsg:
		Keys.start.SetEnabled(c.timer.Running())
		Keys.stop.SetEnabled(!c.timer.Running())
		c.timer, cmd = c.timer.Update(msg)
		return c, cmd
	case timer.TimeoutMsg:
		return c, c.MoveToNext(cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.reset):
			c.timer.Timeout = c.startTime
		case key.Matches(msg, Keys.start, Keys.stop):
			// debounce to prevent double time
			return c, tea.Tick(debounceDuration, func(_ time.Time) tea.Msg {
				return exitMsg{c.timer}
			})
		}
	case exitMsg:
		var cmd tea.Cmd
		cmd = c.timer.Toggle()
		return c, cmd
	case progress.FrameMsg:
		progressModel, cmd := c.progress.Update(msg)
		c.progress = progressModel.(progress.Model)
		return c, cmd
	}

	c.timer, cmd = c.timer.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.GetStyle().Render(lipgloss.JoinVertical(lipgloss.Center, c.progress.View(), "Time Remaining: ", lipgloss.JoinHorizontal(lipgloss.Left, lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04b575")).Render(c.timer.View()))))
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

func (c *Column) MoveToNext(cmd tea.Cmd) tea.Cmd {
	c.timer.Timeout = c.startTime
	return tea.Sequence(cmd, func() tea.Msg { return TimeOut{true} })
}
