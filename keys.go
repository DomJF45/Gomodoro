package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	help  key.Binding
	start key.Binding
	stop  key.Binding
	reset key.Binding
	left  key.Binding
	right key.Binding
	quit  key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.help, k.quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.quit, k.help},
		{k.start, k.stop, k.reset},
		{k.right, k.left},
	}
}

var Keys = KeyMap{
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "start"),
	),
	stop: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "stop"),
	),
	reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset"),
	),
	left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/l", "move left"),
	),
	right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
