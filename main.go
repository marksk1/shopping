package main

import (
	"log"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	
	"go.etcd.io/bbolt"
)

type item struct {
	What    string
	Checked bool
}

type shoppingList struct {
	Name  string
	Items []item

	list        *widget.list
	filterEntry *widget.Entry
}

type appData struct {
	shoppingLists []*shoppingList

	app fyne.app
	win fyne.Window
	tabs *container.DocTabs
}

func main() {
	a := app.NewWithID("github.com.marksk1.shopping")

	myApp := &appData{
		shoppingLists:[]*shoppingList{}, 
		app: a, 
		win: a.NewWindow("Shopping List")}

	myApp.tabs = container.NewDocTabs(nil)
	myApp.tabs.CreateTab = myApp.CreateTab
	myApp.tabs.OnClosed = funct(item*container.TabItem) {
		for index, value := range myApp.shoppingLists {
			if value.Name == item.Text {
				myApp.deleteShoppingList(index, value)
				return
			}
		}
	}

	myApp.tabs.SetTabLocation(container.TabLocationLeading)

	myApp.win.SetContent(container.NewMax(myApp.tabs))
	
	myApp.win.Resize(fyne.NewSize(800,600))
	myApp.win.ShowAndRun()
}

