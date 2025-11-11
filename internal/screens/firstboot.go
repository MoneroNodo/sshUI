package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/backend/i_dbus"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/model/dbus"
)

var firstBoot *FirstBoot = &FirstBoot{}

var (
	pane           *ScreenPane
	password       *ScreenInputField
	passwordRepeat *ScreenInputField
	setPassword    *ScreenButton
	passwordMsg    *ScreenLabel
)

type FirstBoot struct {
	init    bool
	set     bool
	items   []ScreenItem
	current int
}

func NewFirstBoot() *FirstBoot {
	return firstBoot
}

func pwdField(lab string) *ScreenInputField {
	p := NewScreenInputField("", lab, gss.Color(base.CBrightBlack))
	p.Delegate.Width = 24
	p.Delegate.EchoMode = textinput.EchoPassword
	p.Delegate.EchoCharacter = '•'
	return p
}

func (s *FirstBoot) Init() tea.Msg {
	password = pwdField("Password")

	passwordRepeat = pwdField("Repeat Password")

	setPassword = NewScreenButton("Submit", gss.Color(base.CBrightGreen),
		func(sb *ScreenButton) tea.Cmd {
			if password.Delegate.Value() == passwordRepeat.Delegate.Value() {
				i_dbus.Call("setPassword", password.Delegate.Value())
			}
			return nil
		},
	)

	passwordMsg = NewScreenLabel("", gss.Color(base.CBrightRed))

	pane = NewScreenPane("Set your user password", gss.Color(base.CBrightBlue))
	pane.Items = append(pane.Items, password, passwordRepeat, passwordMsg, setPassword)
	s.items = append(s.items, pane)
	s.init = true
	return nil
}

func (s *FirstBoot) PosVertical() gss.Position {
	return gss.Center
}

func (s *FirstBoot) PosHorizontal() gss.Position {
	return gss.Center
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
	if !s.init {
		return nil
	}
	if password.Delegate.Value() == passwordRepeat.Delegate.Value() {
		passwordMsg.label = ""
		setPassword.enabled = len([]byte(password.Delegate.Value())) > 0
	} else {
		passwordMsg.label = "Passwords do not match."
		setPassword.enabled = false
	}
	switch msg := msg.(type) {
	case dbus.DbusSignalMsg:
		switch msg.Signal.(type) {
		case dbus.PasswordChangeStatus:
			AddPopup(
				NewDefaultPopupOK(
					"Password set!",
					"Password changed. Your device will now reboot.",
					gss.Color(base.CGreen),
					func(sb *ScreenButton) tea.Cmd {
						i_dbus.Call("restart")
						return nil
					},
				),
			)
		}
	}
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

func (s *FirstBoot) Interact(m tea.Model) tea.Cmd {
	return s.items[s.current].Interact(m)
}

func (s *FirstBoot) ItemWidth() int {
	return 30
}

func (s *FirstBoot) Vertical() bool {
	return true
}

