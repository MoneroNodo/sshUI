package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
)

var dropToShell *DropToShell = &DropToShell{}

type DropToShell struct {
	current int
}

func NewDropToShell() *DropToShell {
	return dropToShell
}

func (s *DropToShell) Init() tea.Msg {
	return nil
}

func (s *DropToShell) Label() string {
	return "<Drop to Shell>"
}

func (s *DropToShell) View() {}

func (s *DropToShell) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	switch msg := msg.(type) {
	case ScreenActiveChangeMsg:
		if msg.Active {
			switch msg.Screen.(type) {
			case *DropToShell:
				return tea.Quit
			}
		}
	}
	return nil
}

func (s *DropToShell) Items() []ScreenItem {
	return []ScreenItem{}
}

func (s *DropToShell) Current() *int {
	return &s.current
}

func (s *DropToShell) Next() tea.Msg {
	return nil
}

func (s *DropToShell) Prev() tea.Msg {
	return nil
}

func (s *DropToShell) Interact(m tea.Model) tea.Cmd {
	return tea.Quit
}

func (s *DropToShell) PosVertical() gss.Position {
	return 0
}

func (s *DropToShell) PosHorizontal() gss.Position {
	return 0
}

func (s *DropToShell) ItemWidth() int {
	return 0
}

func (s *DropToShell) Vertical() bool {
	return false
}
