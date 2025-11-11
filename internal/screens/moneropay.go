package screens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/moneronodo/sshui/internal/backend/i_moneropay"
	"github.com/moneronodo/sshui/internal/base"
	"github.com/moneronodo/sshui/internal/model/moneropay"
)

var mpay *Moneropay = &Moneropay{}

const (
	mpayUrl = "http://127.0.0.1:5000"
)

var (
	moneropayPane    *ScreenPane
	transactionsPane *ScreenPane

	changeAddrButton *ScreenButton
	clearAddrButton  *ScreenButton

	addrLabel         *ScreenLabel
	statusLabel       *ScreenLabel
	transactionsLabel *ScreenLabel

	transactions []i_moneropay.Transaction
)

type Moneropay struct {
	init    bool
	items   []ScreenItem
	current int
}

func NewMoneropay() *Moneropay {
	return mpay
}

func GetTxDetails(address string) (i_moneropay.Transaction, error) {
	tx := i_moneropay.Transaction{}
	c := &http.Client{Timeout: 5 * time.Second}
	resp, err := c.Get(fmt.Sprintf("%s/receive/%s", mpayUrl, address))
	if err != nil {
		return i_moneropay.Transaction{}, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	j := moneropay.MoneropayReceive{}
	if err := dec.Decode(&j); err != nil {
		spew.Fprintf(base.Dump, "Decode: %v\n", err)
		return i_moneropay.Transaction{}, err
	}
	tx.Covered = j.Amount.Covered
	tx.Complete = j.Complete
	tx.TxIds = j.Transactions
	return tx, nil
}

func _updateMpay(prog *tea.Program) {
	hlt := *i_moneropay.GetHealth(mpayUrl + "/health")
	prog.Send(&moneropay.MpayHealthMsg{
		Health: hlt,
	})
	txs := i_moneropay.GetTxList()
	prog.Send(&i_moneropay.MpayTxListMsg{
		Transactions: txs,
	})
	for i := range txs {
		if i >= len(transactions) || transactions[i].Covered.Total == 0 {
			t, err := GetTxDetails(txs[i].Subaddress)
			if err == nil {
				txs[i].Covered = t.Covered
				txs[i].Queried = true
				prog.Send(&i_moneropay.MpayTxUpdateMsg{
					Index:       i,
					Transaction: txs[i],
				})
			} else {
				spew.Fdump(base.Dump, err)
			}
		}
	}
}

func UpdateMpay(prog *tea.Program) {
	rpcTick := time.NewTicker(5 * time.Second)
	defer rpcTick.Stop()
	_updateMpay(prog)
	for range rpcTick.C {
		_updateMpay(prog)
	}
}

func (s *Moneropay) Init() tea.Msg {
	addr, _ := base.GetVal("moneropay", "deposit_address").(string)
	addrLabel = NewScreenLabel(addr, gss.Color(base.CGray))
	statusLabel = NewScreenLabel("status pending...", gss.Color(base.CBrightYellow))
	moneropayPane = NewScreenPane("MoneroPay", gss.Color(base.CYellow))
	clearAddrButton = NewScreenButton("Clear Address (disable MoneroPay)", gss.Color(base.CBrightYellow),
		func(sb *ScreenButton) tea.Cmd {
			AddPopup(NewDefaultPopupYesNo("Restart", "Are you sure?", gss.Color(base.CBrightRed),
				func(sb *ScreenButton) tea.Cmd {
					base.SetMpayConfig("deposit_address", "")
					base.SetMpayConfig("enabled", base.Bool(false))
					return nil
				},
				nil,
			))
			return nil
		})
	changeAddrButton = NewScreenButton("Update Address", gss.Color(base.CBrightYellow),
		func(sb *ScreenButton) tea.Cmd {
			var inp *ScreenInputField
			inp = NewScreenInputField(addrLabel.label, "address", gss.Color(base.CWhite))
			AddPopup(NewDefaultPopupYesNo("Restart", "Are you sure?", gss.Color(base.CBrightRed),
				func(sb *ScreenButton) tea.Cmd {
					if base.ValidateAddr(inp.Delegate.Value()) {
						base.SetMpayConfig("deposit_address", inp.Delegate.Value)
						base.SetMpayConfig("enabled", base.Bool(true))
					}
					return nil
				},
				nil,
			))
			return nil
		})
	moneropayPane.Items = append(moneropayPane.Items,
		statusLabel,
		NewScreenHr(90, gss.Color(base.CBrightBlack)),
		addrLabel,
		changeAddrButton,
		clearAddrButton,
	)
	transactionsLabel = NewScreenLabel("", gss.Color(base.CBlue))
	transactionsPane = NewScreenPane("Recent Transactions", gss.Color(base.CBrightAqua),
		transactionsLabel,
	)
	s.items = append(s.items, transactionsPane, moneropayPane)
	s.init = true
	return nil
}

var mpayStatusStyle = gss.NewStyle().Bold(true)
var stringed = map[bool]string{ // lol
	true:  "ready",
	false: "dead",
}

var (
	mpayTxAddrStyle        = gss.NewStyle().Foreground(gss.Color(base.CPurple))
	mpayTxAddrComplStyle   = gss.NewStyle().Foreground(gss.Color(base.CBrightGreen))
	mpayTxAddrIncomplStyle = gss.NewStyle().Foreground(gss.Color(base.CBrightPurple))
	mpayTxAmountStyle      = gss.NewStyle().Foreground(gss.Color(base.CBrightAqua))
)

func (s *Moneropay) Update(msg tea.Msg, m tea.Model) tea.Cmd {
	switch m := msg.(type) {
	case *i_moneropay.MpayTxUpdateMsg:
		transactions[m.Index] = m.Transaction
	case *i_moneropay.MpayTxListMsg:
		transactions = m.Transactions
	case *moneropay.MpayHealthMsg:
		if statusLabel != nil {
			str := fmt.Sprintf("Wallet %s, SQLite %s",
				stringed[m.Health.Services.Walletrpc],
				stringed[m.Health.Services.Sqlite],
			)
			switch m.Health.Status {
			case 200:
				statusLabel.label = mpayStatusStyle.Foreground(gss.Color(base.CBrightGreen)).
					Render(str)
			case 503:
				statusLabel.label = mpayStatusStyle.Foreground(gss.Color(base.CBrightYellow)).
					Render(str)
			case 0:
				statusLabel.label = mpayStatusStyle.Foreground(gss.Color(base.CRed)).
					Render("dead")

			}
		}
	}
	return nil
}

func (s *Moneropay) View() {
	if !s.init {
		return
	}
	var sb strings.Builder
	for _, t := range transactions {
		var st = gss.Style{}
		if t.Complete {
			st = mpayTxAddrComplStyle
		} else {
			st = mpayTxAddrIncomplStyle
		}
		if t.Queried {
			sb.WriteString(fmt.Sprintf("%s\n %s/%s XMR   %s\n",
				st.Render(shorthandAddress(t.Subaddress, 4, 4)),
				mpayTxAmountStyle.Render(fmt.Sprintf("%.4f", float64(t.Covered.Unlocked)*float64(10e-12))),
				mpayTxAmountStyle.Render(fmt.Sprintf("%.4f", float64(t.Expected)*float64(10e-12))),
				base.UnixTime(t.CreatedAt.Unix()),
			))
			for _, t := range t.TxIds {
				sb.WriteString(fmt.Sprintf("%s %s  %s XMR\n",
					t.TxHash,
					base.UnixTime(t.Timestamp.Unix()),
					mpayTxAmountStyle.Render(fmt.Sprintf("%.4f", float64(t.Amount)*float64(10e-12))),
				))
			}
			sb.WriteString("\n")
		} else {
			sb.WriteString(fmt.Sprintf("%s\n ?/%s XMR   %s\n\n",
				mpayTxAddrStyle.Render(shorthandAddress(t.Subaddress, 4, 4)),
				mpayTxAmountStyle.Render(fmt.Sprintf("%.4f", float64(t.Expected)*float64(10e-12))),
				base.UnixTime(t.CreatedAt.Unix()),
			))
		}
	}
	transactionsLabel.label = sb.String()
}

func (s *Moneropay) Label() string {
	return "MoneroPay"
}

func (s *Moneropay) Items() []ScreenItem {
	return s.items
}

func (s *Moneropay) Current() *int {
	return &s.current
}

func (s *Moneropay) Next() tea.Msg {
	return UpdateFocus(s, 1)
}

func (s *Moneropay) Prev() tea.Msg {
	return UpdateFocus(s, -1)
}

func (s *Moneropay) Interact(m tea.Model) tea.Cmd {
	return s.items[s.current].Interact(m)
}

func (s *Moneropay) PosVertical() gss.Position {
	return gss.Position(0.8)
}

func (s *Moneropay) PosHorizontal() gss.Position {
	return gss.Center
}

func (s *Moneropay) ItemWidth() int {
	return 8
}

func (s *Moneropay) Vertical() bool {
	return false
}
