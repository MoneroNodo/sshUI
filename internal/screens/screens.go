package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gss "github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	"github.com/moneronodo/sshui/internal/base"
)

var (
	Popups []Popup
)

type ScreenToggleAction func(*ScreenToggle, bool) tea.Msg
type ScreenButtonAction func(*ScreenButton) tea.Msg

type ScreenLeaveMsg struct {
	Current Screen
}

type ScreenEnterMsg struct {
	Current Screen
}

type FocusChangeMsg struct {
	Current ScreenItem
}

type Popup interface {
	Title() string
	Body() string
	Items() []ScreenItem
	Prev() tea.Msg
	Next() tea.Msg
	Interact(model tea.Model) tea.Msg
	Render() string
	Current() int
}

type DefaultPopup struct {
	title   string
	body    string
	items   []ScreenItem
	color   gss.Color
	current int
}

type Screen interface {
	Init() tea.Msg
	View()
	Update(msg tea.Msg, model tea.Model) tea.Cmd
	Label() string
	Items() []ScreenItem
	Current() *int
	Next() tea.Msg
	Prev() tea.Msg
	Interact(model tea.Model) tea.Msg
}

type ScreenItem interface {
	Render() string
	IsEnabled() bool
	SetFocus(bool)
	SetColor(gss.Color)
	GetColor() gss.Color
	IsFocus() bool
	Interact(model tea.Model) tea.Msg
}

type ScreenPane struct {
	Title   string
	Current int
	Items   []ScreenItem
	Focus   bool
	Color   gss.Color
	Style   gss.Style
}

func NewScreenPane(title string, color gss.Color, items ...ScreenItem) *ScreenPane {
	sp := new(ScreenPane)
	sp.Title = title
	sp.Color = color
	sp.Items = items
	return sp
}

func (sp *ScreenPane) Interact(m tea.Model) tea.Msg {
	if sp.Current < len(sp.Items) && sp.Items[sp.Current].IsEnabled() {
		return sp.Items[sp.Current].Interact(m)
	}
	return nil
}

func (sp *ScreenPane) Render() string {
	var s []string
	s = append(s, labelInner.Bold(true).Render(sp.Title))
	for _, i := range sp.Items {
		s = append(s, i.Render())
	}
	return sp.Style.Render(gss.JoinVertical(
		gss.Left,
		s...,
	))
}

func (sp *ScreenPane) IsEnabled() bool {
	for _, i := range sp.Items {
		if i.IsEnabled() {
			return true
		}
	}
	return false
}

func (sp *ScreenPane) SetFocus(focus bool) {
	for _, i := range sp.Items {
		i.SetFocus(false)
	}
	if len(sp.Items) > sp.Current {
		sp.Items[sp.Current].SetFocus(focus)
	}
	sp.Focus = focus
}

func (sb *ScreenPane) IsFocus() bool {
	return sb.Focus
}

func (sp *ScreenPane) SetColor(color gss.Color) {
	sp.Color = color
}

func (sp *ScreenPane) GetColor() gss.Color {
	return sp.Color
}

type ScreenTable struct {
	Delegate table.Model
	color    gss.Color
	Style    gss.Style
}

func NewScreenTable(columns []table.Column, rows []table.Row, color gss.Color) *ScreenTable {
	si := new(ScreenTable)
	si.Delegate = table.New(table.WithColumns(columns), table.WithRows(rows))
	// TODO
	si.SetColor(color)
	return si
}

func (st *ScreenTable) Interact(m tea.Model) tea.Msg { return nil }
func (st *ScreenTable) Render() string {
	return st.Style.Render(st.Delegate.View())
}

func (st *ScreenTable) SetFocus(focus bool) {}

func (st *ScreenTable) IsFocus() bool {
	return false
}

func (st *ScreenTable) IsEnabled() bool {
	return false
}

func (st *ScreenTable) SetColor(color gss.Color) {
	st.color = color
	st.Style = gss.NewStyle().Foreground(st.color)
}

func (st *ScreenTable) GetColor() gss.Color {
	return st.color
}

type ScreenLabel struct {
	label string
	color gss.Color
	Style gss.Style
	focus bool
}

func NewScreenLabel(text string, color gss.Color) *ScreenLabel {
	si := new(ScreenLabel)
	si.label = text
	si.SetColor(color)
	return si
}

func (sl *ScreenLabel) Interact(m tea.Model) tea.Msg { return nil }
func (sl *ScreenLabel) Render() string {
	return sl.Style.Render(sl.label)
}

func (sl *ScreenLabel) SetFocus(focus bool) {
	sl.focus = focus
}

func (sl *ScreenLabel) IsFocus() bool {
	return sl.focus
}

func (sl *ScreenLabel) IsEnabled() bool {
	return false
}

func (sl *ScreenLabel) SetColor(color gss.Color) {
	sl.color = color
	sl.Style = gss.NewStyle().Foreground(sl.color)
}

func (sl *ScreenLabel) GetColor() gss.Color {
	return sl.color
}

type ScreenInputField struct {
	Delegate      textinput.Model
	focus         bool
	enabled       bool
	color         gss.Color
	Style         gss.Style
	StyleLabel    gss.Style
	StyleFocus    gss.Style
	StyleDisabled gss.Style
}

func NewScreenInputField(text string, placeholder string, color gss.Color) *ScreenInputField {
	si := new(ScreenInputField)
	si.Delegate.Placeholder = placeholder
	si.SetColor(color)
	si.Delegate.TextStyle = gss.NewStyle().Foreground(gss.Color(base.CGray))
	si.Delegate.TextStyle = gss.NewStyle().Foreground(gss.Color(base.CWhite))
	si.StyleDisabled = gss.NewStyle().Foreground(gss.Color(base.CBrightBlack))
	si.enabled = true
	return si
}

func (si *ScreenInputField) Interact(m tea.Model) tea.Msg {
	return si.Delegate.Focus() // FIX
}

func (si *ScreenInputField) Render() string {
	return ""
}

func (si *ScreenInputField) SetFocus(focus bool) {
	si.focus = focus
	if focus {
		si.Delegate.TextStyle = gss.NewStyle().Foreground(gss.Color(base.CGray))
	} else {
		si.Delegate.TextStyle = gss.NewStyle().Foreground(gss.Color(base.CWhite))
	}
}

func (sl *ScreenInputField) IsFocus() bool {
	return sl.focus
}

func (si *ScreenInputField) IsEnabled() bool {
	return si.enabled
}

func (si *ScreenInputField) SetColor(color gss.Color) {
	si.color = color
	si.Style = gss.NewStyle().Foreground(si.color)
}

func (si *ScreenInputField) GetColor() gss.Color {
	return si.color
}

type ScreenToggle struct {
	label         string
	focus         bool
	enabled       bool
	toggled       bool
	Action        ScreenToggleAction
	color         gss.Color
	Style         gss.Style
	StyleLabel    gss.Style
	StyleFocus    gss.Style
	StyleDisabled gss.Style
}

func NewScreenToggle(label string, color gss.Color, action ScreenToggleAction) *ScreenToggle {
	sb := new(ScreenToggle)
	sb.label = label
	sb.SetColor(color)
	sb.StyleFocus = gss.NewStyle().Foreground(gss.Color(base.CWhite))
	sb.StyleLabel = gss.NewStyle().Foreground(gss.Color(base.CGray))
	sb.StyleDisabled = gss.NewStyle().Foreground(gss.Color(base.CBrightBlack))
	sb.Action = action
	sb.enabled = true
	return sb
}

func (st *ScreenToggle) Interact(m tea.Model) tea.Msg {
	st.toggled = !st.toggled
	spew.Fdump(base.Dump, st)
	if st.Action == nil {
		return nil
	}
	return st.Action(st, st.toggled)
}

func (st *ScreenToggle) Render() string {
	l := st.label
	sel := " "
	if st.toggled {
		sel = "X"
	}
	if !st.enabled {
		r := st.StyleDisabled.Render(fmt.Sprintf("(%s) %s", sel, l))
		return r
	} else if st.focus {
		r := st.Style.Render("[") +
			st.StyleLabel.Render(fmt.Sprintf("%s", sel)) +
			st.Style.Render("] ") +
			st.StyleLabel.Render(l)
		return r
	} else {
		r := st.Style.Render(fmt.Sprintf("(%s) %s", sel, l))
		return r
	}
}

func (st *ScreenToggle) SetFocus(focus bool) {
	st.focus = focus
}

func (st *ScreenToggle) IsFocus() bool {
	return st.focus
}

func (st *ScreenToggle) IsEnabled() bool {
	return st.enabled
}

func (st *ScreenToggle) SetColor(color gss.Color) {
	st.color = color
	st.Style = gss.NewStyle().Foreground(st.color)
}

func (st *ScreenToggle) GetColor() gss.Color {
	return st.color
}

type ScreenButton struct {
	label         string
	focus         bool
	enabled       bool
	Action        ScreenButtonAction
	color         gss.Color
	Style         gss.Style
	StyleLabel    gss.Style
	StyleFocus    gss.Style
	StyleDisabled gss.Style
}

func NewScreenButton(label string, color gss.Color, action ScreenButtonAction) *ScreenButton {
	sb := new(ScreenButton)
	sb.label = label
	sb.SetColor(color)
	sb.StyleFocus = gss.NewStyle().Foreground(gss.Color(base.CWhite))
	sb.StyleLabel = gss.NewStyle().Foreground(gss.Color(base.CGray))
	sb.StyleDisabled = gss.NewStyle().Foreground(gss.Color(base.CBrightBlack))
	sb.Action = action
	sb.enabled = true
	return sb
}

func (sb *ScreenButton) Interact(m tea.Model) tea.Msg {
	if sb.Action == nil {
		return nil
	}
	return sb.Action(sb)
}

func (sb *ScreenButton) Render() string {
	l := sb.label
	if !sb.enabled {
		r := sb.StyleDisabled.Render(fmt.Sprintf("  %s  ", l))
		return r
	} else if sb.focus {
		r := sb.Style.Render("[ ") + sb.StyleLabel.Render(l) + sb.Style.Render(" ]")
		return r
	} else {
		r := sb.Style.Render(fmt.Sprintf("  %s  ", l))
		return r
	}
}

func (sb *ScreenButton) SetFocus(focus bool) {
	sb.focus = focus
}

func (sb *ScreenButton) IsFocus() bool {
	return sb.focus
}

func (sb *ScreenButton) IsEnabled() bool {
	return sb.enabled
}

func (sb *ScreenButton) SetColor(color gss.Color) {
	sb.color = color
	sb.Style = gss.NewStyle().Foreground(sb.color)
}

func (sb *ScreenButton) GetColor() gss.Color {
	return sb.color
}

func enabledItems(items []ScreenItem) []int {
	en := []int{}
	for i, v := range items {
		if v.IsEnabled() {
			en = append(en, i)
		}
	}
	return en
}

func indexOf(s *[]int, val int) int {
	if len((*s)) == 0 {
		return -1
	} else if len((*s)) == 1 {
		return (*s)[0]
	}
	for i, v := range *s {
		if val == v {
			return i
		}
	}
	return (*s)[0]
}

func WrapPane(s *ScreenPane, mod int) int {
	en := enabledItems(s.Items)
	l := len(en)
	switch l {
	case 0:
		s.Current = 0
		return mod
	case 1:
		s.Current = en[0]
		return mod
	}
	index := indexOf(&en, s.Current)
	if index == -1 {
		return mod
	}
	if index+mod < 0 {
		s.Current = en[0]
		return mod
	} else if index+mod >= l {
		s.Current = en[l-1]
		return mod
	} else {
		s.Current = en[index+mod]
	}
	return 0
}

func WrapScreen(s Screen, mod int) {
	en := enabledItems(s.Items())
	l := len(en)
	switch l {
	case 0:
		*s.Current() = 0
	case 1:
		*s.Current() = en[0]
		return
	}
	index := indexOf(&en, *s.Current())
	if index == -1 {
		*s.Current() = 0
		return
	}
	if index+mod < 0 {
		*s.Current() = en[l-1]
	} else if index+mod >= l {
		*s.Current() = en[0]
	} else {
		*s.Current() = en[index+mod]
	}
}

func SetFocusScreen(s Screen) {
	en := enabledItems(s.Items())
	for _, v := range en {
		s.Items()[v].SetFocus(v == *s.Current())
	}
}

func SetFocusPane(s *ScreenPane) {
	en := enabledItems(s.Items)
	for _, v := range en {
		s.Items[v].SetFocus(v == s.Current)
	}
}

func SetFocusPopup(s Popup) {
	for i, v := range s.Items() {
		v.SetFocus(i == s.Current())
	}
}

func UpdateFocus(s Screen, mod int) tea.Msg {
	if len(s.Items()) == 0 {
		return nil
	}
	c := *s.Current()
	switch t := s.Items()[c].(type) { // take panes into account
	case *ScreenPane:
		modret := WrapPane(t, mod)
		SetFocusPane(t)
		if modret != 0 {
			WrapScreen(s, modret)
		}
	default:
		WrapScreen(s, mod)
	}

	if *s.Current() != c { // if switching between panes or items not in a pane
		switch t := s.Items()[*s.Current()].(type) {
		case *ScreenPane:
			en := enabledItems(t.Items)
			if len(en) > 1 {
				if mod < 0 { // if moving backward and landing on another pane, set to last item
					t.Current = en[len(en)-1]
				} else if mod > 0 { // if moving forward, set to first item
					t.Current = en[0]
				}
			}
			WrapPane(t, 0)
			SetFocusPane(t)
		}
	}

	SetFocusScreen(s)

	return FocusChangeMsg{
		Current: s.Items()[*s.Current()],
	}
}

func newDefaultPopup(title string, body string, color gss.Color) *DefaultPopup {
	popup := new(DefaultPopup)
	popup.title = title
	popup.body = body
	popup.color = color
	return popup
}

func NewDefaultPopupOK(title string, body string, color gss.Color,
	ok ScreenButtonAction, items ...ScreenItem) *DefaultPopup {
	popup := newDefaultPopup(title, body, color)
	popup.items = append(items,
		NewScreenButton("OK", gss.Color(base.CBrightBlue), ok),
	)
	return popup
}

func NewDefaultPopupOKCancel(title string, body string, color gss.Color,
	ok ScreenButtonAction, cancel ScreenButtonAction, items ...ScreenItem) *DefaultPopup {
	popup := newDefaultPopup(title, body, color)
	popup.items = append(items,
		NewScreenButton("OK", gss.Color(base.CBrightBlue), ok),
		NewScreenButton("Cancel", gss.Color(base.CBrightYellow), cancel),
	)
	return popup
}

func NewDefaultPopupYesNo(title string, body string, color gss.Color,
	yes ScreenButtonAction, no ScreenButtonAction, items ...ScreenItem) *DefaultPopup {
	popup := newDefaultPopup(title, body, color)
	popup.items = append(items,
		NewScreenButton("Yes", gss.Color(base.CBrightGreen), yes),
		NewScreenButton("No", gss.Color(base.CBrightYellow), no),
	)
	return popup
}

func (dp *DefaultPopup) Title() string {
	return dp.title
}

func (dp *DefaultPopup) Body() string {
	return dp.body
}

func (dp *DefaultPopup) Items() []ScreenItem {
	return dp.items
}

func (dp *DefaultPopup) Prev() tea.Msg {
	if len(dp.items) <= 1 {
		dp.current = 0
	} else {
		dp.current--
		if dp.current < 0 {
			dp.current = len(dp.items) - 1
		}
	}
	SetFocusPopup(dp)
	return FocusChangeMsg{
		Current: dp.items[dp.current],
	}
}

func (dp *DefaultPopup) Next() tea.Msg {
	if len(dp.items) <= 1 {
		dp.current = 0
	} else {
		dp.current++
		if dp.current >= len(dp.items) {
			dp.current = 0
		}
	}
	SetFocusPopup(dp)
	return FocusChangeMsg{
		Current: dp.items[dp.current],
	}
}

func (dp *DefaultPopup) Interact(model tea.Model) tea.Msg {
	switch len(dp.items) {
	case 0:
		return nil
	case 1:
		return dp.items[0].Interact(model)
	default:
		return dp.items[dp.current].Interact(model)
	}
}

func (dp *DefaultPopup) Render() string {
	var s []string
	s = append(s, labelInner.Bold(true).Render(dp.title), labelInner.Render(dp.body))
	for _, i := range dp.items {
		s = append(s, i.Render())
	}
	return gss.JoinVertical(
		gss.Left,
		s...,
	)
}

func (dp *DefaultPopup) Current() int {
	return dp.current
}

func AddPopup(popup Popup) {
	SetFocusPopup(popup)
	Popups = append([]Popup{popup}, Popups...)
}
