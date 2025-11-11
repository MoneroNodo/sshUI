package base

import (
	"math"

	gss "github.com/charmbracelet/lipgloss"
)

const (
	CBlack        = "0"
	CBrightBlack  = "8"
	CRed          = "1"
	CBrightRed    = "9"
	CGreen        = "2"
	CBrightGreen  = "10"
	CYellow       = "3"
	CBrightYellow = "11"
	CBlue         = "4"
	CBrightBlue   = "12"
	CPurple       = "5"
	CBrightPurple = "13"
	CAqua         = "6"
	CBrightAqua   = "14"
	CGray         = "7"
	CWhite        = "15"
)

type Styles struct {
	Border      gss.Color
	TabArea     gss.Style // tabs area
	Tabs        gss.Style // tab
	TabsHg      gss.Style // highlighted tab
	ContentArea gss.Style
	ContentItem gss.Style
	Label       gss.Style
}

func InitStyles(size, w, h, items float64) *Styles {
	var st = new(Styles)
	st.Border = gss.Color(CBrightBlack)
	st.Tabs = gss.NewStyle().Foreground(gss.Color(CWhite))
	st.TabsHg = gss.NewStyle().Foreground(gss.Color(CBlack)).Background(gss.Color(CGray))
	st.TabArea = gss.NewStyle().PaddingRight(2)
	st.ContentArea = gss.NewStyle().
		Padding(1).BorderStyle(gss.NormalBorder())
	st.ContentItem = gss.NewStyle().
		Padding(
			int(math.Max(1, size*h*0.01/items)),
			int(math.Max(2, size*w*0.02/items)),
		).
		Margin(
			int(size*h*0.001/items),
			int(size*w*0.002/items),
		).BorderStyle(gss.NormalBorder())
	st.Label = gss.NewStyle().Margin(1)
	return st
}
