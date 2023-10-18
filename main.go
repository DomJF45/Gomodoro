package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Rotation int

const timeout = time.Minute * 25

const (
	pomodoro Rotation = iota
	short_break
	long_break
)

const margin = 3

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
	if _, err := tea.NewProgram(b, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
