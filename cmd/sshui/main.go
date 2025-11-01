package main

import (
	"log"
	"math"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"

	// "github.com/davecgh/go-spew/spew"
	"github.com/moneronodo/sshui/internal/backend"
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
	return msgs
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// if base.Dump != nil {
	//   spew.Fdump(base.Dump, msg)
	// }
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	switch mt := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = mt.Width
		m.height = mt.Height
		tabsizeX, tabsizeY := m.styles.TabArea.GetFrameSize()
		contentsizeX, contentsizeY := m.styles.ContentArea.GetFrameSize()
		m.tabsPort.Width = TabAreaWid - tabsizeX
		m.tabsPort.Height = m.height - tabsizeY
		m.contentPort.Width = (m.width - TabAreaWid) - contentsizeX - 1
		m.contentPort.Height = m.height - contentsizeY

		m.tabsPort, cmd = m.tabsPort.Update(msg)
		cmds = append(cmds, cmd)

		m.contentPort, cmd = m.contentPort.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		if len(screens.Popups) > 0 {
			curpopup := screens.Popups[0]
			switch mt.String() {
			case "up":
				return m, curpopup.Prev
			case "down":
				return m, curpopup.Next
			case "enter":
				curpopup.Interact(m)
				if len(curpopup.Items()) > 0 {
					switch curpopup.Items()[curpopup.Current()].(type) {
					case *screens.ScreenButton:
						if len(screens.Popups) <= 1 {
							screens.Popups = nil
						} else {
							screens.Popups = screens.Popups[1:]
						}
					}
				}
				return m, nil
			default:
				return m, nil
			}
		}

		curscreen := m.screens[m.current]
		switch mt.String() {
		case "down":
			if m.active {
				return m, curscreen.Next
			} else {
				m.current++
				if m.current >= len(m.screens) {
					m.current = 0
				}
				return m, nil
			}
		case "up":
			if m.active {
				return m, curscreen.Prev
			} else {
				m.current--
				if m.current < 0 {
					m.current = len(m.screens) - 1
				}
				return m, nil
			}
		case "enter":
			if m.active {
				curscreen.Items()[*curscreen.Current()].Interact(m)
				return m, nil
			} else {
				m.active = true
				screens.UpdateFocus(curscreen, 0)
				return m, nil
			}
		case "ctrl+c", "esc":
			if m.active {
				m.active = false
				return m, nil
			} else {
				return m, tea.Quit
			}
		}
	}
	for _, v := range m.screens {
		cmds = append(cmds, v.Update(msg, m))
	}

	return m, tea.Batch(cmds...)
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

const colInactive = gss.Color(base.CWhite)
const colActive = gss.Color(base.CBrightBlack)

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
	sv := gss.Place(
		m.contentPort.Width,
		m.contentPort.Height,
		gss.Left,
		gss.Top,
		gss.JoinHorizontal(gss.Top, it...),
	)
	var popups = ""
	if len(screens.Popups) > 0 {
		popups = gss.Place(
			m.width,
			m.height,
			gss.Center,
			gss.Center,
			m.styles.ContentArea.BorderStyle(gss.ThickBorder()).
				Width(int(math.Min(80, math.Max(20, float64(m.width)*0.6)))).
				Render(screens.Popups[0].Render()),
		)
	}
	return gss.JoinHorizontal(
		gss.Top,
		m.styles.TabArea.Render(m.renderTabs(tab)),
		m.styles.ContentArea.Foreground(cont).BorderForeground(cont).Render(sv),
	) + popups
}

// /home/nodo/variables/firstboot
func initModel() model {
	var dump *os.File
	var err error
	dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		// spew.Fdump(base.Dump, err)
		os.Exit(1)
	}
	base.Dump = dump
	m := model{
		screens:   []screens.Screen{},
		current:   0,
		active:    false,
		firstBoot: base.IsFirstBoot(),
		styles:    base.InitStyles(),
	}
	if false && m.firstBoot {
		m.screens = append(m.screens,
			screens.NewFirstBoot(),
		)
	} else {
		m.screens = append(m.screens,
			screens.NewDashboard(),
			screens.NewNode(),
			screens.NewSystem(),
			// TODO ... the rest
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
	go backend.DbusSignals(prog)
	go screens.UpdateBody(prog)
	if _, err := prog.Run(); err != nil {
		log.Fatal(err)
	}
}
