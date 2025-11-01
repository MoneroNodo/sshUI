package screens

import (
	// "fmt"

	"fmt"
	"net"

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
	clearnetPort int = 18089
	onionPort    int = 18089
	i2pPort      int = 18089
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

func (s *Node) Init() tea.Msg {
	clearnetAddr = getClearnetIp()
	onionAddr = base.GetVal("tor_address").(string)
	i2pAddr = base.GetVal("i2p_address").(string)

	clearnetLabel = NewScreenLabel(fmt.Sprintf("%s:%d", clearnetAddr, clearnetPort), gss.Color(base.CPurple))
	onionLabel = NewScreenLabel(fmt.Sprintf("%s:%d", onionAddr, onionPort), gss.Color(base.CPurple))
	i2pLabel = NewScreenLabel(fmt.Sprintf("%s:%d", i2pAddr, i2pPort), gss.Color(base.CPurple))

	hiddenRpcToggle = NewScreenToggle("Enable Hidden RPC", gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Msg {
			base.SetConfig("anon_rpc", toggled)
			return nil
		})
	torToggle = NewScreenToggle("Enable Tor", gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Msg {
			base.SetConfig("tor_enabled", toggled)
			return nil
		})
	torAllToggle = NewScreenToggle("Route All Through Tor", gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Msg {
			base.SetConfig("torproxy_enabled", toggled)
			return nil
		})
	i2pToggle = NewScreenToggle("Enable I2P", gss.Color(base.CGreen),
		func(sb *ScreenToggle, toggled bool) tea.Msg {
			base.SetConfig("i2p_enabled", toggled)
			return nil
		})

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
	UpdateFocus(s, 0)
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

func (s *Node) Interact(m tea.Model) tea.Msg {
	return s.items[s.current].Interact(m)
}
