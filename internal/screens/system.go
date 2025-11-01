package screens

import (
	// "fmt"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/backend"
	"github.com/moneronodo/sshui/internal/base"
	// "github.com/moneronodo/sshui/internal/base"
)

var system *System = &System{}

var (
	rebootButton   *ScreenButton
	shutdownButton *ScreenButton
	recoveryButton *ScreenButton

	sysPane              *ScreenPane
	recoveryFSToggle     *ScreenToggle
	recoveryResyncToggle *ScreenToggle
)

type System struct {
	init    bool
	items   []ScreenItem
	current int
}

func NewSystem() *System {
	return system
}

func (s *System) Init() tea.Msg {
	rebootButton = NewScreenButton("Reboot", gss.Color(base.CYellow), func(sb *ScreenButton) tea.Msg {
		AddPopup(NewDefaultPopupYesNo("Restart", "Are you sure?", gss.Color(base.CBrightRed),
			func(sb *ScreenButton) tea.Msg {
				backend.DbusCall("restart")
				return nil
			},
			nil,
		))
		return nil
	})
	shutdownButton = NewScreenButton("Shutdown", gss.Color(base.CRed), func(sb *ScreenButton) tea.Msg {
		AddPopup(NewDefaultPopupYesNo("Shutdown", "Are you sure?", gss.Color(base.CBrightRed),
			func(sb *ScreenButton) tea.Msg {
				backend.DbusCall("shutdown")
				return nil
			},
			nil,
		))
		return nil
	})
	recoveryFSToggle = NewScreenToggle(
		"Recover Filesystem",
		gss.Color(base.CYellow),
		nil,
	)

	recoveryResyncToggle = NewScreenToggle(
		"Purge & Resync Blockchain",
		gss.Color(base.CYellow),
		nil,
	)
	recoveryButton = NewScreenButton("Start Recovery", gss.Color(base.CBrightPurple), func(sb *ScreenButton) tea.Msg {
		AddPopup(NewDefaultPopupOKCancel("Recovery", "Select your recovery options, then press OK.", gss.Color(base.CYellow),
			func(sb *ScreenButton) tea.Msg {
				backend.DbusCall("startRecovery", recoveryFSToggle.toggled, recoveryResyncToggle.toggled)
				return nil
			},
			nil,
			recoveryFSToggle,
			recoveryResyncToggle,
		))
		return nil
	})

	sysPane = NewScreenPane(
		"Power",
		gss.Color(base.CAqua),
		rebootButton,
		shutdownButton,
		recoveryButton,
	)

	s.items = append(
		s.items,
		sysPane,
	)
	s.init = true
	UpdateFocus(s, 0)
	return nil
}

func (s *System) Label() string {
	return "System"
}

func (s *System) View() {
	if !s.init {
		return
	}
}

func (s *System) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	return nil
}

func (s *System) Items() []ScreenItem {
	return s.items
}

func (s *System) Current() *int {
	return &s.current
}

func (s *System) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *System) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *System) Interact(m tea.Model) tea.Msg {
	return s.items[s.current].Interact(m)
}
