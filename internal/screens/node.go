package screens

import (
	// "fmt"

	"fmt"
	"net"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/base"
)

func getClearnetIp() string {
	conn, err := net.Dial("udp", "192.168.1.1:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

var node *Node = &Node{}

const (
	nodePort int = 18089
)

var (
	clearnetAddr string
	onionAddr    string
	i2pAddr      string

	clearnetLabel *ScreenLabel
	onionLabel    *ScreenLabel
	i2pLabel      *ScreenLabel

	clearnetPane *ScreenPane
	onionPane    *ScreenPane
	i2pPane      *ScreenPane

	torToggle       *ScreenToggle
	torAllToggle    *ScreenToggle
	hiddenRpcToggle *ScreenToggle
	i2pToggle       *ScreenToggle
)

type Node struct {
	init    bool
	items   []ScreenItem
	current int
}

func NewNode() *Node {
	return node
}

var valueStyle = gss.NewStyle().Foreground(gss.Color(base.CWhite))

func newInputIntBtn(label, val string) *ScreenButton {
	v, _ := base.GetVal(val).(float64)
	str := fmt.Sprintf("%s: %s", label, valueStyle.Render(strconv.Itoa(int(v))))
	var btn *ScreenButton
	btn = NewScreenButton(str, gss.Color(base.CGreen),
		func(sb *ScreenButton) tea.Cmd {
			v, _ := base.GetVal(val).(float64)
			in := NewScreenInputField(strconv.Itoa(int(v)), strconv.Itoa(int(v)), gss.Color(base.CGray))
			AddPopup(NewDefaultPopupOKCancel(label, "Set new value", gss.Color(base.CGreen),
				func(sb *ScreenButton) tea.Cmd {
					i, err := strconv.Atoi(in.Delegate.Value())
					if err == nil {
						base.SetConfig(val, i)
						base.SaveConfigFile()
					}
					btn.label = fmt.Sprintf("%s: %s", label, valueStyle.Render(in.Delegate.Value()))
					return nil
				}, nil,
				in,
			))
			return nil
		})
	return btn
}

func newInputStrBtn(label, val string, secret bool) *ScreenButton {
	var (
		str string
		btn *ScreenButton
	)
	if secret {
		str = label
	} else {
		v, _ := base.GetVal(val).(string)
		str = fmt.Sprintf("%s: %s", label, valueStyle.Render(v))
	}
	btn = NewScreenButton(str, gss.Color(base.CGreen),
		func(sb *ScreenButton) tea.Cmd {
			v, _ := base.GetVal(val).(string)
			in := NewScreenInputField(v, v, gss.Color(base.CGray))
			AddPopup(NewDefaultPopupOKCancel(label, "Set new value", gss.Color(base.CBrightGreen),
				func(sb *ScreenButton) tea.Cmd {
					base.SetConfig(val, in.Delegate.Value())
					if !secret {
						btn.label = fmt.Sprintf("%s: %s", label, valueStyle.Render(in.Delegate.Value()))
					}
					return base.SaveConfigFile
				}, nil,
				in,
			))
			return nil
		})
	return btn
}

func newToggle(label, val string) *ScreenToggle {
	toggle := NewScreenToggle(label, gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Cmd {
			base.SetConfig(val, toggled)
			return base.SaveConfigFile
		})
	toggle.toggled, _ = base.GetVal(val).(bool)
	return toggle
}

func (s *Node) Init() tea.Msg {
	hiddenRpcToggle = newToggle("Hidden RPC", "anon_rpc")
	torToggle = newToggle("Enable Tor", "tor_enabled")
	torAllToggle = newToggle("Route All Through Tor", "tor_global_enabled")
	i2pToggle = newToggle("Enable I2P", "i2p_enabled")

	clearnetAddr = getClearnetIp()
	onionAddr, _ = base.GetVal("tor_address").(string)
	i2pAddr, _ = base.GetVal("i2p_address").(string)

	clearnetLabel = NewScreenLabel(fmt.Sprintf("%s:%d", clearnetAddr, nodePort), gss.Color(base.CPurple))
	onionLabel = NewScreenLabel(fmt.Sprintf("%s:%d", onionAddr, nodePort), gss.Color(base.CPurple))
	i2pLabel = NewScreenLabel(fmt.Sprintf("%s:%d", i2pAddr, nodePort), gss.Color(base.CPurple))

	clearnetPane = NewScreenPane(
		"Clearnet",
		gss.Color(base.CBlue),
		clearnetLabel,
		hiddenRpcToggle,
	)
	onionPane = NewScreenPane(
		"Tor",
		gss.Color(base.CBlue),
		onionLabel,
		torToggle,
		torAllToggle,
	)
	i2pPane = NewScreenPane(
		"I2P",
		gss.Color(base.CBlue),
		i2pLabel,
		i2pToggle,
	)

	clearnetPane.Style = clearnetPane.Style.Width(40)
	onionPane.Style = onionPane.Style.Width(70)
	i2pPane.Style = i2pPane.Style.Width(70)

	s.items = append(s.items, clearnetPane, onionPane, i2pPane)
	s.init = true
	return nil
}

func (s *Node) Label() string {
	return "Node"
}

func (s *Node) View() {
	if !s.init {
		return
	}
}

func (s *Node) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	return nil
}

func (s *Node) Items() []ScreenItem {
	return s.items
}

func (s *Node) Current() *int {
	return &s.current
}

func (s *Node) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *Node) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *Node) Interact(m tea.Model) tea.Cmd {
	return s.items[s.current].Interact(m)
}

func (s *Node) PosVertical() gss.Position {
	return gss.Position(0.8)
}

func (s *Node) PosHorizontal() gss.Position {
	return gss.Center
}

func (s *Node) ItemWidth() int {
	return 4
}

func (s *Node) Vertical() bool {
	return true
}
