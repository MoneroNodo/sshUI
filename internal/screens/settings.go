package screens

import (
	// "fmt"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/base"
)

var settings *Settings = &Settings{}

var (
	inPeerButton    *ScreenButton
	outPeerButton   *ScreenButton
	upSpeedButton   *ScreenButton
	downSpeedButton *ScreenButton
	rpcUserButton   *ScreenButton
	rpcPassButton   *ScreenButton
	banlistButton   *ScreenButton

	settingsDataPane    *ScreenPane
	settingsPrivacyPane *ScreenPane

	privateRPCToggle *ScreenToggle
)

type Settings struct {
	init    bool
	items   []ScreenItem
	current int
}

func NewSettings() *Settings {
	return settings
}

func newBanlistToggle(label, val string) *ScreenToggle {
	toggle := NewScreenToggle(label, gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Cmd {
			base.SetBanlistConfig(val, toggled)
			return base.SaveConfigFile
		})
	toggle.toggled, _ = base.GetVal("banlists", val).(bool)
	return toggle
}

func (s *Settings) Init() tea.Msg {

	inPeerButton = newInputIntBtn("Incoming Peers", "in_peers")
	outPeerButton = newInputIntBtn("Outgoing Peers", "out_peers")
	upSpeedButton = newInputIntBtn("Upload Speed (kB/s)", "limit_rate_up")
	downSpeedButton = newInputIntBtn("Download Speed (kB/s)", "limit_rate_down")

	privateRPCToggle = newToggle("RPC Authentication", "rpc_enabled")
	rpcUserButton = newInputStrBtn("RPC Username", "rpcu", false)
	rpcPassButton = newInputStrBtn("RPC Password", "rpcp", true)
	banlistButton = NewScreenButton("Banlist Settings", gss.Color(base.CBrightYellow),
		func(sb *ScreenButton) tea.Cmd {
			var (
				boogToggle   = newBanlistToggle("Boog900", "boog900")
				dnsToggle    = newBanlistToggle("DNS", "dns")
				guixmrToggle = newBanlistToggle("gui.xmr.pm", "gui-xmr-pm")
			)
			AddPopup(NewDefaultPopupOK("Banlist Settings", "", gss.Color(base.CBrightYellow),
				func(sb *ScreenButton) tea.Cmd {
					return nil
				},
				boogToggle,
				dnsToggle,
				guixmrToggle,
			))
			return nil
		})

	settingsDataPane = NewScreenPane(
		"Data",
		gss.Color(base.CBrightAqua),
		inPeerButton,
		outPeerButton,
		upSpeedButton,
		downSpeedButton,
	)
	settingsPrivacyPane = NewScreenPane(
		"Privacy",
		gss.Color(base.CBrightPurple),
		privateRPCToggle,
		rpcUserButton,
		rpcPassButton,
		NewScreenHr(20, gss.Color(base.CBrightBlack)),
		banlistButton,
	)
	s.items = append(s.items, settingsDataPane, settingsPrivacyPane)
	s.init = true
	return nil
}

func (s *Settings) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	return nil
}

func (s *Settings) View() {
	if !s.init {
		return
	}

}

func (s *Settings) Label() string {
	return "Settings"
}

func (s *Settings) Items() []ScreenItem {
	return s.items
}

func (s *Settings) Current() *int {
	return &s.current
}

func (s *Settings) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *Settings) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *Settings) Interact(m tea.Model) tea.Cmd {
	return s.items[s.current].Interact(m)
}

func (s *Settings) PosVertical() gss.Position {
	return gss.Position(0.8)
}

func (s *Settings) PosHorizontal() gss.Position {
	return gss.Center
}

func (s *Settings) ItemWidth() int {
	return 4
}

func (s *Settings) Vertical() bool {
	return false
}

