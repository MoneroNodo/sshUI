package screens

import (
	// "fmt"

	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	i_lws "github.com/moneronodo/sshui/internal/backend/lws"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/model/lws"
)

var lightWallet *LightWallet = &LightWallet{}

const (
	lwsPort int = 18089
)

var (
	lwsClearnetAddr *ScreenLabel
	lwsI2pAddr      *ScreenLabel
	lwsOnionAddr    *ScreenLabel

	lwsAddWalletButton *ScreenButton

	lwsAccounts lws.LwsListAccounts

	lwsInfoPane     *ScreenPane
	lwsAccountsPane *ScreenPane
)

type LightWallet struct {
	init    bool
	items   []ScreenItem
	current int
}

func NewLightWallet() *LightWallet {
	return lightWallet
}

var inputAcct, inputKey *ScreenInputField

func (s *LightWallet) Init() tea.Msg {
	lwsClearnetAddr = NewScreenLabel(fmt.Sprintf(" http://%s:%d", clearnetAddr, lwsPort), gss.Color(base.CBlue))
	lwsOnionAddr = NewScreenLabel(fmt.Sprintf(" http://%s:%d", onionAddr, lwsPort), gss.Color(base.CBlue))
	lwsI2pAddr = NewScreenLabel(fmt.Sprintf(" http://%s:%d", i2pAddr, lwsPort), gss.Color(base.CBlue))

	inputAcct = NewScreenInputField("", "Primary address", gss.Color(base.CWhite))
	inputKey = NewScreenInputField("", "Private view key", gss.Color(base.CWhite))

	lwsAddWalletButton = NewScreenButton("Add Wallet", gss.Color(base.CBrightGreen),
		func(sb *ScreenButton) tea.Cmd {
			e := i_lws.AddAccount(inputAcct.Delegate.Value(), inputKey.Delegate.Value())
			if e != nil {
				AddPopup(
					NewDefaultPopupOK("Couldn't add wallet", e.Error(), gss.Color(base.CBrightRed), nil),
				)
			} else {
				inputAcct.Delegate.SetValue("")
				inputKey.Delegate.SetValue("")
			}
			return nil
		})

	lwsInfoPane = NewScreenPane(
		"Connection Info",
		gss.Color(base.CPurple),
		NewScreenLabel("Clearnet", gss.Color(base.CWhite)),
		lwsClearnetAddr,
		NewScreenLabel("Tor", gss.Color(base.CWhite)),
		lwsOnionAddr,
		NewScreenLabel("I2P", gss.Color(base.CWhite)),
		lwsI2pAddr,
		NewScreenHr(90, gss.Color(base.CBrightBlack)),
		inputAcct,
		inputKey,
		lwsAddWalletButton,
	)

	lwsAccountsPane = NewScreenPane(
		"Active Accounts",
		gss.Color(base.CGreen),
	)

	UpdateAccounts()
	s.items = append(s.items, lwsInfoPane, lwsAccountsPane)
	s.init = true
	return nil
}

func newLwsAccountButton(l lws.LwsAccount) *ScreenButton {
	var b *ScreenButton
	if l.Status == lws.LwsAccountActive {
		b = NewScreenButton("Deactivate", gss.Color(base.CBlue),
			func(sb *ScreenButton) tea.Cmd {
				i_lws.DeactivateAccount(l.Address)
				UpdateAccounts()
				return nil
			})
	} else {
		b = NewScreenButton("Activate", gss.Color(base.CBlue),
			func(sb *ScreenButton) tea.Cmd {
				i_lws.ReactivateAccount(l.Address)
				UpdateAccounts()
				return nil
			})
	}
	sb := NewScreenButton(shorthandAddress(l.Address, 3, 4), gss.Color(base.CWhite),
		func(sb *ScreenButton) tea.Cmd {
			p := newDefaultPopup(
				l.Address,
				fmt.Sprintf(
					"Last accessed: %s (%s)\nHeight: %s",
					base.UnixTime(l.AccessTime),
					base.UnixTimeRelative(l.AccessTime),
					strconv.Itoa(int(l.ScanHeight)),
				),
				gss.Color(base.CGray),
			)
			p.items = append(p.items,
				b,
				NewScreenButton("Delete", gss.Color(base.CRed),
					func(sb *ScreenButton) tea.Cmd {
						AddPopup(
							NewDefaultPopupYesNo(
								"Delete Account",
								fmt.Sprintf(
									"%s\nAre you sure you want to delete this account? This action cannot be undone.",
									l.Address,
								),
								gss.Color(base.CRed),
								func(sb *ScreenButton) tea.Cmd {
									i_lws.DeleteAccount(l.Address)
									UpdateAccounts()
									return nil
								}, nil),
						)
						return nil
					}),
				NewScreenButton("Rescan", gss.Color(base.CYellow),
					func(sb *ScreenButton) tea.Cmd {
						AddPopup(
							NewDefaultPopupYesNo("Rescan", "This may take a while, are you sure?", gss.Color(base.CRed),
								func(sb *ScreenButton) tea.Cmd {
									i_lws.Rescan(l.Address, int(l.ScanHeight))
									UpdateAccounts()
									return nil
								}, nil),
						)
						return nil
					}),
				NewScreenButton("Close", gss.Color(base.CGreen), nil),
			)
			AddPopup(p)
			return nil
		})
	return sb
}

func shorthandAddress(addr string, groups int, size int) string {
	if len(addr) < 95 || groups < 1 {
		return addr
	}
	var sb strings.Builder
	for i := range groups {
		sb.WriteString(addr[i*size : i*size+size])
		sb.WriteByte(' ')
	}
	sb.WriteString("... ")
	for x := range groups {
		i := groups - x - 1 // reverse range
		sb.WriteString(addr[95-i*size-size : 95-i*size])
		sb.WriteByte(' ')
	}

	return sb.String()
}

func UpdateAccounts() {
	var err error
	lwsAccounts, err = i_lws.ListAccounts()
	if err != nil {
		spew.Fdump(base.Dump, lwsAccounts, err)
		return
	}
	lwsAccountsPane.Items = nil
	lwsAccountsPane.Items = append(lwsAccountsPane.Items,
		NewScreenLabel("Active", gss.Color(base.CBrightGreen)),
	)
	for _, a := range lwsAccounts.Active {
		lwsAccountsPane.Items = append(lwsAccountsPane.Items, newLwsAccountButton(a))
	}
	lwsAccountsPane.Items = append(lwsAccountsPane.Items,
		NewScreenLabel("Inactive", gss.Color(base.CBrightRed)),
	)
	for _, a := range lwsAccounts.Inactive {
		lwsAccountsPane.Items = append(lwsAccountsPane.Items, newLwsAccountButton(a))
	}
}

func (s *LightWallet) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	return nil
}

func (s *LightWallet) View() {
	if !s.init {
		return
	}

}

func (s *LightWallet) Label() string {
	return "LightWallet"
}

func (s *LightWallet) Items() []ScreenItem {
	return s.items
}

func (s *LightWallet) Current() *int {
	return &s.current
}

func (s *LightWallet) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *LightWallet) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *LightWallet) Interact(m tea.Model) tea.Cmd {
	return s.items[s.current].Interact(m)
}

func (s *LightWallet) PosVertical() gss.Position {
	return gss.Position(0.8)
}

func (s *LightWallet) PosHorizontal() gss.Position {
	return gss.Center
}

func (s *LightWallet) ItemWidth() int {
	return 4
}

func (s *LightWallet) Vertical() bool {
	return true
}
