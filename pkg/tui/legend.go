package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *Application) initBottomLegend() *tview.TextView {
	legend := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetTextAlign(tview.AlignCenter).
		SetText("Press ESC to quit")

	legend.Box.SetBorder(true)

	return legend
}
