package main

import "github.com/charmbracelet/lipgloss"

type Tab struct {
	title   string
	focused bool
}

func (t *Tab) Focus() {
	t.focused = true
}

func (t *Tab) Blur() {
	t.focused = false
}

func (t *Tab) IsFocused() bool {
	return t.focused
}

func NewTab(t string, f bool) Tab {
	return Tab{
		title:   t,
		focused: f,
	}
}

func (t Tab) GetStyle() lipgloss.Style {
	if t.focused {
		return lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("202")).
			Padding(0, 3).
			Underline(true).
			MarginTop(1).
			MarginRight(1)
	} else {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(1)
	}
}

func (t Tab) View() string {
	return t.GetStyle().Render(t.title)
}
