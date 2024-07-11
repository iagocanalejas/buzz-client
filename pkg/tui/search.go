package tui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/iagocanalejas/buzz-client/internal/api"
	"github.com/rivo/tview"
)

func (app *Application) initSearchView() *tview.Flex {
	app.searchInput = tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	app.searchInput.SetChangedFunc(func(text string) {
		app.currentSearch = strings.ToLower(text)
	})

	legend := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetTextAlign(tview.AlignLeft).
		SetText("Filters -> ")

	searchBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.searchInput, 1, 0, true).
		AddItem(legend, 1, 0, false)
	searchBox.Box.SetBorder(true)

	return searchBox
}

func (app *Application) filterLinks() {
	var links []*api.Link

	if app.currentSearch == "" {
		app.populateList(app.links)
		app.filteredLinks = nil
		return
	}

	for _, link := range app.links {
		if strings.Contains(strings.ToLower(link.Name), app.currentSearch) {
			links = append(links, link)
		}
	}
	app.filteredLinks = links
	app.populateList(links)
}
