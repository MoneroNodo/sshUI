package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/base"
)

var firstBoot *FirstBoot = &FirstBoot{}

type FirstBoot struct {
	init           bool
	password       *ScreenInputField
	passwordRepeat *ScreenInputField
	set            bool
	items          []ScreenItem
	current        int
}

func NewFirstBoot() *FirstBoot {
	return firstBoot
}

func (s *FirstBoot) Init() tea.Msg {
	s.password = NewScreenInputField("", "Password", gss.Color(base.CBrightBlack))
	s.passwordRepeat = NewScreenInputField("", "Repeat Password", gss.Color(base.CBrightBlack))
	s.items = append(s.items, s.password, s.passwordRepeat)
	s.init = true
	return nil
}

func (s *FirstBoot) Label() string {
	return "Initial Setup"
}

func (s *FirstBoot) View() {
	if !s.init {
		return
	}
	var buf []string
	if s.set {
		buf = []string{"Password set! Your device will now reboot."}
	} else {
		for _, i := range s.items {
			buf = append(buf, i.Render())
		}
	}
	//return gss.JoinHorizontal(gss.Center, buf...)
}

func (s *FirstBoot) Update(msg tea.Msg, m tea.Model) tea.Cmd {

	return nil
}

func (s *FirstBoot) Items() []ScreenItem {
	return s.items
}

func (s *FirstBoot) Current() *int {
	return &s.current
}

func (s *FirstBoot) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *FirstBoot) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *FirstBoot) Interact(m tea.Model) tea.Msg {
	return s.items[s.current].Interact(m)
}
