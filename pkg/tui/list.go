package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/iagocanalejas/buzz-client/internal/api"
	"github.com/rivo/tview"
)

func (app *Application) initListView() *tview.List {
	app.list = tview.NewList()
	app.list.Box.SetBorder(true)

	app.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		// disable default behavior fot TAB key (next list item) as we use that to change focus
		case tcell.KeyTab:
			return nil
		case tcell.KeyEnter:
			if !app.showingDetails {
				selectedItem := app.list.GetCurrentItem()
				app.loadFolderLinks(app.links[selectedItem])
			}
		}
		return event
	})

	app.loadFolders()

	return app.list
}

func (app *Application) populateList(links []*api.Link) {
	app.list.Clear()
	for _, link := range links {
		app.list.AddItem(fmt.Sprintf("%s (%s)", link.Name, link.Href), link.Name, 0, nil)
	}
}

func (app *Application) loadFolderLinks(link *api.Link) {
	links, err := app.api.ListFolder(link)
	if err != nil {
		app.errorModal(err)
		return
	}

	app.link = link
	app.links = links
	app.showingDetails = true
	app.populateList(links)
}

func (app *Application) loadFolders() {
	links, err := app.api.List()
	if err != nil {
		app.errorModal(err)
		return
	}

	app.link = nil
	app.links = links
	app.showingDetails = false
	app.populateList(links)
}
