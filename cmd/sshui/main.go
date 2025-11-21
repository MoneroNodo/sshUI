package main

import (
	"log"
	"math"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/moneronodo/sshui/internal/backend/i_dbus"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/screens"
)

const TabAreaWid int = 12

var prog *tea.Program

type model struct {
	screens     []screens.Screen
	current     int  // cursor
	active      bool // whether focused on cur screen
	width       int
	height      int
	firstBoot   bool
	styles      *base.Styles
	tabsPort    viewport.Model
	contentPort viewport.Model
}

func (m model) Init() tea.Cmd {
	return m.initScreens
}

func (m model) initScreens() tea.Msg {
	var msgs []tea.Msg
	for _, v := range m.screens {
		msg := v.Init()
		msgs = append(msgs, msg)
	}
	screens.UpdateFocus(m.screens[m.current], 0)
	return msgs
}

func (m *model) updateStyles() {
	count := 1
	if len(m.screens) > 0 &&
		len(m.screens[m.current].Items()) > 0 {
		count = len(m.screens[m.current].Items())
	}
	m.styles = base.InitStyles(
		float64(m.screens[m.current].ItemWidth()),
		float64(m.width),
		float64(m.height),
		float64(count),
	)
}

func closePopup() {
	if len(screens.Popups) <= 1 {
		screens.Popups = nil
	} else {
		screens.Popups = screens.Popups[1:]
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	switch mt := msg.(type) {
	case base.ConfigSavedMsg:
		exec.Command("/usr/bin/systemctl", "restart", "monerod")
	case tea.WindowSizeMsg:
		m.width = mt.Width
		m.height = mt.Height
		tabsizeX, tabsizeY := m.styles.TabArea.GetFrameSize()
		contentsizeX, contentsizeY := m.styles.ContentArea.GetFrameSize()
		m.tabsPort.Width = TabAreaWid - tabsizeX
		m.tabsPort.Height = m.height - tabsizeY
		m.contentPort.Width = (m.width - TabAreaWid) - contentsizeX - 4
		m.contentPort.Height = m.height - contentsizeY

		m.tabsPort, cmd = m.tabsPort.Update(msg)
		cmds = append(cmds, cmd)

		m.contentPort, cmd = m.contentPort.Update(msg)
		cmds = append(cmds, cmd)
		m.updateStyles()
	case tea.KeyMsg:
		if len(screens.Popups) > 0 {
			curpopup := screens.Popups[0]
			switch mt.String() {
			case "up":
				return m, curpopup.Prev
			case "down":
				return m, curpopup.Next
			case "esc", "ctrl+c":
				closePopup()
				return m, nil
			case "enter":
				c := curpopup.Interact(m)
				if len(curpopup.Items()) > 0 {
					switch curpopup.Items()[curpopup.Current()].(type) {
					case *screens.ScreenButton:
						closePopup()
					}
				}
				return m, c
			default:
				return m, nil
			}
		}

		curscreen := m.screens[m.current]
		switch mt.String() {
		case "down", "tab":
			if m.active {
				return m, curscreen.Next
			} else {
				m.current++
				if m.current >= len(m.screens) {
					m.current = 0
				}
				m.updateStyles()
				return m, nil
			}
		case "up", "shift+tab":
			if m.active {
				return m, curscreen.Prev
			} else {
				m.current--
				if m.current < 0 {
					m.current = len(m.screens) - 1
				}
				m.updateStyles()
				return m, nil
			}
		case "enter":
			if m.active {
				if len(curscreen.Items()) == 0 {
					return m, nil
				}
				c := curscreen.Interact(m)
				if c == nil {
					c = curscreen.Next
				}
				cmds = append(cmds, c)
			} else {
				m.active = true
				screens.UpdateFocus(curscreen, 0)
				return m, m.sendActive
			}
		case "esc", "ctrl+c":
			if m.active && len(m.screens) > 1 {
				m.active = false
				return m, m.sendActive
			} else {
				return m, tea.Quit
			}
		}
	}
	for _, v := range m.screens {
		for _, i := range v.Items() {
			vi := i.Update(msg, m)
			if vi != nil {
				cmds = append(cmds, vi)
			}
		}
		vu := v.Update(msg, m)
		if vu != nil {
			cmds = append(cmds, v.Update(msg, m))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *model) sendActive() tea.Msg {
	return screens.ScreenActiveChangeMsg{
		Active: m.active,
		Screen: m.screens[m.current],
	}
}

func (m model) renderTabs(col gss.Color) string {
	var sb []string
	for i, s := range m.screens {
		if i == m.current {
			sb = append(sb, m.styles.TabsHg.Background(col).Render(s.Label()))
		} else {
			sb = append(sb, m.styles.Tabs.Foreground(col).Render(s.Label()))
		}
	}
	return gss.JoinVertical(gss.Left, sb...)
}

const colActive = gss.Color(base.CWhite)
const colInactive = gss.Color(base.CGray)

func (m model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	if len(m.screens) <= m.current || len(m.screens) == 0 ||
		m.screens[m.current] == nil {
		return ""
	}
	if m.width < TabAreaWid*2 || m.height < len(m.screens) {
		return "..."
	}
	var (
		it        []string
		tab, cont gss.Color
	)
	if m.active {
		tab = colInactive
		cont = colActive
	} else {
		tab = colActive
		cont = colInactive
	}
	m.screens[m.current].View()
	for _, i := range m.screens[m.current].Items() {
		if i.IsFocus() {
			it = append(it, m.styles.ContentItem.Foreground(cont).BorderForeground(i.GetColor()).Render(i.Render()))
		} else {
			it = append(it, m.styles.ContentItem.Foreground(cont).BorderForeground(colInactive).Render(i.Render()))
		}
	}
	var sv string
	if m.screens[m.current].Vertical() {
		sv = gss.Place(
			m.contentPort.Width,
			m.contentPort.Height,
			m.screens[m.current].PosHorizontal(),
			m.screens[m.current].PosVertical(),
			gss.JoinVertical(gss.Left, it...),
		)
	} else {
		sv = gss.Place(
			m.contentPort.Width,
			m.contentPort.Height,
			m.screens[m.current].PosHorizontal(),
			m.screens[m.current].PosVertical(),
			gss.JoinHorizontal(gss.Top, it...),
		)
	}
	var popups = ""
	if len(screens.Popups) > 0 {
		popups = gss.Place(
			m.width,
			m.height,
			gss.Center,
			gss.Center,
			m.styles.ContentArea.BorderStyle(gss.ThickBorder()).
				Width(4+int(math.Max(float64(screens.Popups[0].Width()), math.Max(20, float64(m.width)*0.6)))).
				Render(screens.Popups[0].Render()),
		)
	}
	return gss.JoinHorizontal(
		gss.Top,
		m.styles.TabArea.Render(m.renderTabs(tab)),
		m.styles.ContentArea.Foreground(cont).BorderForeground(cont).Render(sv),
	) + popups
}

func initModel() model {
	var dump *os.File
	var err error
	dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		os.Exit(1)
	}
	base.Dump = dump
	m := model{
		screens:   []screens.Screen{},
		current:   0,
		active:    false,
		firstBoot: base.IsFirstBoot(),
		styles:    base.InitStyles(1, 0, 0, 1),
	}
	if m.firstBoot {
		m.screens = append(m.screens,
			screens.NewFirstBoot(),
		)
		m.active = true
	} else {
		m.screens = append(m.screens,
			screens.NewDashboard(),
			screens.NewNode(),
			screens.NewSettings(),
			screens.NewSystem(),
			screens.NewLightWallet(),
			screens.NewMoneropay(),
		)
	}
	// Set cursors properly, taking unselectable items into account
	for _, s := range m.screens {
		for _, i := range s.Items() {
			switch p := i.(type) {
			case *screens.ScreenPane:
				screens.WrapPane(p, 0)
			}
		}
		screens.WrapScreen(s, 0)
	}
	return m
}

func main() {
	f, err := tea.LogToFile("debug.log", "dbg:")
	if err != nil {
		log.Fatal("rip")
	}
	defer f.Close()
	prog = tea.NewProgram(initModel(), tea.WithAltScreen())
	go i_dbus.Signals(prog)
	go screens.UpdateRPC(prog)
	go screens.UpdateMpay(prog)
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
