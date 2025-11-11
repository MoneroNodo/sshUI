package screens

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/moneronodo/sshui/internal/backend/daemonrpc"
	"github.com/moneronodo/sshui/internal/base"
	rpc_model "github.com/moneronodo/sshui/internal/model/daemonrpc"
	dbus_model "github.com/moneronodo/sshui/internal/model/dbus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var dash *Dashboard = &Dashboard{}

var (
	nodeText     *ScreenLabel
	hardwareText *ScreenLabel
	serviceText  *ScreenLabel

	nodePane     *ScreenPane
	hardwarePane *ScreenPane
	servicePane  *ScreenPane
)

type hardwareStatus struct {
	cpuUsg    float32
	cpuGhz    float32
	temp      float32
	ram       float32
	ramTotal  float32
	ssd       float32
	ssdTotal  float32
	emmc      float32
	emmcTotal float32
	uptime    string
}

type serviceStatus struct {
	monerod   string
	tor       string
	i2pd      string
	moneroLws string
	sshd      string
	moneropay string
}

type Dashboard struct {
	init       bool
	getInfo    rpc_model.DaemonResponseBodyGetInfo
	getVersion rpc_model.DaemonResponseBodyGetVersion
	hardware   hardwareStatus
	service    serviceStatus
	items      []ScreenItem
	current    int
}

func NewDashboard() *Dashboard {
	return dash
}

func (s *Dashboard) Init() tea.Msg {
	var (
		height = 9
		color  = "3"
	)
	hardwareText = NewScreenLabel("", gss.Color(color))
	serviceText = NewScreenLabel("", gss.Color(color))
	nodeText = NewScreenLabel("", gss.Color(color))
	hardwareText.Style = hardwareText.Style.Height(height)
	serviceText.Style = serviceText.Style.Height(height)
	nodeText.Style = nodeText.Style.Height(height)

	hardwarePane = NewScreenPane("Hardware", gss.Color(base.CPurple))
	servicePane = NewScreenPane("Services", gss.Color(base.CPurple))
	nodePane = NewScreenPane("Node Status", gss.Color(base.CPurple))
	hardwarePane.Items = append(hardwarePane.Items, hardwareText)
	servicePane.Items = append(servicePane.Items, serviceText)
	nodePane.Items = append(nodePane.Items, nodeText)

	s.items = append(s.items, nodePane, servicePane, hardwarePane)
	s.init = true
	return nil
}

func (s *Dashboard) Label() string {
	return "Dashboard"
}

var labelInner = gss.NewStyle().Padding(1)

func (s *Dashboard) getSyncStatus() string {
	var (
		hei     int  = s.getInfo.Height
		tarHei  int  = s.getInfo.TargetHeight
		syn     bool = s.getInfo.Synchronized
		syncing bool = s.getInfo.BusySyncing
	)
	synPer := fmt.Sprintf("%.0f", (float32(hei)/float32(tarHei))*100)
	if syn || hei == tarHei {
		return "Synchronized (100%)"
	} else if syncing {
		return fmt.Sprintf("Synchronizing (%s)", synPer)
	} else {
		return "Not synchronizing"
	}
}

func (s *Dashboard) View() {
	if !s.init {
		return
	}
	var (
		update string = "Unknown"
		online string = "Unknown"
	)
	if s.getInfo.Offline {
		online = "Not connected"
	} else {
		online = "Connected"
	}
	if s.getInfo.UpdateAvailable {
		update = "Update available"
	} else {
		update = "Up to date"
	}
	nodeText.label = gss.JoinHorizontal(
		gss.Top,
		labelInner.Render(s.getSyncStatus()+`
Block height
Version
Out peers
In  peers
Update
Network`),
		labelInner.Render(fmt.Sprintf(
			`
: %d
: %s
: %d
: %d
: %s
: %s`,
			s.getInfo.Height,
			s.getInfo.Version,
			s.getInfo.OutgoingConnectionsCount,
			s.getInfo.IncomingConnectionsCount,
			update,
			online,
		)),
	)
	hardwareText.label = gss.JoinHorizontal(
		gss.Top,
		labelInner.Render(`CPU
Temperature
RAM
Blockchain
Storage
Uptime`),
		labelInner.Render(fmt.Sprintf(
			`: %.1f Ghz (%.0f%%)
: %.0f°C
: %.1f/%.1f GB (%.0f%%)
: %.1f/%.1f TB (%.0f%%)
: %.1f/%.1f GB (%.0f%%)
: %s`,
			s.hardware.cpuGhz,
			s.hardware.cpuUsg,
			s.hardware.temp,
			s.hardware.ram,
			s.hardware.ramTotal,
			s.hardware.ram/s.hardware.ramTotal*100,
			s.hardware.ssd,
			s.hardware.ssdTotal,
			s.hardware.ssd/s.hardware.ssdTotal*100,
			s.hardware.emmc,
			s.hardware.emmcTotal,
			s.hardware.emmc/s.hardware.emmcTotal*100,
			s.hardware.uptime,
		)))
	serviceText.label = gss.JoinHorizontal(
		gss.Top,
		labelInner.Render(`Node
Tor
I2P
LWS
MoneroPay`),
		labelInner.Render(fmt.Sprintf(
			`: %s
: %s
: %s
: %s
: %s`,
			s.service.monerod,
			s.service.tor,
			s.service.i2pd,
			s.service.moneroLws,
			s.service.moneropay,
		)))
}

func (s *Dashboard) Items() []ScreenItem {
	return s.items
}

func (s *Dashboard) Current() *int {
	return &s.current
}

func (s *Dashboard) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	switch m := msg.(type) {
	case dbus_model.DbusSignalMsg:
		switch sig := m.Signal.(type) {
		case dbus_model.ServiceStatusReadyNotification:
			updateServices(s, sig.Message)
		case dbus_model.HardwareStatusReadyNotification:
			updateStatuses(s, sig.Message)
		}
	case rpc_model.DaemonRPCMsg:
		resp := m.Response.(rpc_model.DaemonRPCResponse)
		if resp.Error != nil {
			// TODO handle
		} else {
			switch resp.Result.(type) {
			case *rpc_model.DaemonResponseBodyGetInfo:
				s.getInfo = *resp.Result.(*rpc_model.DaemonResponseBodyGetInfo)
			case *rpc_model.DaemonResponseBodyGetVersion:
				s.getVersion = *resp.Result.(*rpc_model.DaemonResponseBodyGetVersion)
			}
		}
		return nil
	}
	return nil
}

func updateServices(s *Dashboard, str string) {
	spl := strings.Split(str, "\n")
	c := cases.Title(language.Und)
	if len(spl) < 6 {
		spew.Fprint(base.Dump, "updateHardware: split too small: ")
		spew.Fdump(base.Dump, spl)
		return
	}
	s.service.monerod = c.String(strings.Split(spl[0], ":")[1])
	s.service.tor = c.String(strings.Split(spl[1], ":")[1])
	s.service.i2pd = c.String(strings.Split(spl[2], ":")[1])
	s.service.moneroLws = c.String(strings.Split(spl[3], ":")[1])
	s.service.sshd = c.String(strings.Split(spl[4], ":")[1])
	s.service.moneropay = c.String(strings.Split(spl[5], ":")[1])

}

func convFl(val string) float32 {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return -1
	}
	return float32(f)
}

func updateStatuses(s *Dashboard, str string) {
	spl := strings.Split(str, "\n")
	if len(spl) < 10 {
		spew.Fprint(base.Dump, "updateServices: split too small: ")
		spew.Fdump(base.Dump, spl)
		return
	}
	s.hardware.cpuUsg = convFl(spl[0])
	s.hardware.cpuGhz = convFl(spl[1])
	s.hardware.ram = convFl(spl[2])
	s.hardware.ramTotal = convFl(spl[3])
	s.hardware.temp = convFl(spl[4])
	s.hardware.ssd = convFl(spl[5])
	s.hardware.ssdTotal = convFl(spl[6])
	s.hardware.emmc = convFl(spl[7])
	s.hardware.emmcTotal = convFl(spl[8])
	s.hardware.uptime = spl[9]
}

func _updateRpc(prog *tea.Program) {
	j := *daemonrpc.DaemonPost("http://127.0.0.1:18081/json_rpc", daemonrpc.DaemonRequestBodyGetInfo())
	prog.Send(rpc_model.DaemonRPCMsg{
		Response: j,
	})
}

func UpdateRPC(prog *tea.Program) {
	rpcTick := time.NewTicker(5 * time.Second)
	defer rpcTick.Stop()
	_updateRpc(prog)
	for range rpcTick.C {
		_updateRpc(prog)
	}
}

func (s *Dashboard) Next() tea.Msg {
	return &FocusChangeMsg{}
}

func (s *Dashboard) Prev() tea.Msg {
	return &FocusChangeMsg{}
}

func (s *Dashboard) Interact(m tea.Model) tea.Cmd {
	return nil
}

func (s *Dashboard) PosVertical() gss.Position {
	return gss.Position(0.8)
}

func (s *Dashboard) PosHorizontal() gss.Position {
	return gss.Center
}

func (s *Dashboard) ItemWidth() int {
	return 4
}

func (s *Dashboard) Vertical() bool {
	return false
}
