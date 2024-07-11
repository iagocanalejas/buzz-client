package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/iagocanalejas/buzz-client/internal/api"
	"github.com/rivo/tview"
)

type Application struct {
	App *tview.Application

	api *api.API

	link           *api.Link
	links          []*api.Link
	filteredLinks  []*api.Link
	currentSearch  string // current search keywords
	hasError       bool   // if the error modal is showing or not
	showingDetails bool   // if the details view is in display

	flex        *tview.Flex
	searchInput *tview.InputField
	list        *tview.List
}

func BuildApp() *Application {
	app := &Application{
		App: tview.NewApplication().EnableMouse(true),
		api: api.Init(),
	}

	app.setupListeners()
	app.initFlex()

	app.App.SetRoot(app.flex, true)

	return app
}

func (app *Application) initFlex() {
	app.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.initSearchView(), 4, 0, false).
		AddItem(app.initListView(), 0, 1, true).
		AddItem(app.initBottomLegend(), 3, 0, false)
}

func (app *Application) setupListeners() {
	app.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.hasError {
			return event
		}
		switch event.Key() {
		case tcell.KeyEnter:
			if app.searchInput.HasFocus() {
				// configure search on <CR> press
				app.filterLinks()
			}
		case tcell.KeyTab:
			app.nextFocus()
		case tcell.KeyEsc:
			if app.showingDetails {
				app.App.SetRoot(app.flex, true)
				app.loadFolders()
			} else {
				app.App.Stop()
			}
		}
		return event
	})
}

func (app *Application) nextFocus() {
	if app.searchInput.HasFocus() {
		app.App.SetFocus(app.list)
	} else {
		app.App.SetFocus(app.searchInput)
	}
}

func (app *Application) errorModal(err error) {
	if app.hasError {
		return
	}

	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Continue"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.App.SetRoot(app.flex, true).SetFocus(app.flex)
			app.hasError = false
		})

	modal.
		SetBackgroundColor(tcell.ColorDarkRed).
		SetTextColor(tcell.ColorYellow).
		SetBorder(true).
		SetBorderColor(tcell.ColorWhite).
		SetBorderPadding(2, 2, 2, 2)

	app.App.SetRoot(modal, true).SetFocus(modal)
	app.hasError = true
}
