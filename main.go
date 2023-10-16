package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Rotation int

const timeout = time.Minute * 25

const (
	pomodoro Rotation = iota
	short_break
	long_break
)

const margin = 5

var (
	timerStyle = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.HiddenBorder())

	timerFocused = lipgloss.NewStyle().Width(15).Height(5).Align(lipgloss.Center, lipgloss.Center).BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("69"))
)

func (r Rotation) GetNext() Rotation {
	if r == long_break {
		return pomodoro
	} else {
		return r + 1
	}
}

func (r Rotation) GetPrev() Rotation {
	if r == pomodoro {
		return long_break
	} else {
		return r - 1
	}
}

func main() {
	b := NewBoard()
	b.InitTimers()
	if _, err := tea.NewProgram(b).Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
