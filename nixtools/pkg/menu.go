package menu

import (
	"github.com/rivo/tview"
)

func generateBox() *tview.Box {
	return tview.NewBox().SetBorder(true).SetTitle("My Box")
}

func generateList(app *tview.Application) *tview.List {
	list := tview.NewList().
		AddItem("Run Tools", "", '1', nil).
		AddItem("Quit", "", '2', func() { app.Stop() })
	return list
}

func GenerateApp() *tview.Application {
	app := tview.NewApplication()
	list := generateList(app)
	app.SetRoot(list, true).SetFocus(list).Run()
	return app
}
