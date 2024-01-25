package main

func (a *appData) CreateTab() *container.TabItem {
	minimalPlaceIndex := a.minimalPlaceIndex()
	newShoppingList, err := a.newShoppingList(fmt.Sprintf("Unknown place %d", minimalPlaceIndex))
	if err != nil {
		return nil
	}
	newDocItem := a.buildTabItem(newShoppingList)
	newShoppingLocationEntry := widget.NewEntry()
	dialog.ShowForm(
		"New shopping place", 
		"Create", 
		"Cancel",
		[]*widget.FormItem{
			{
				Text: "Name", 
				Widget: newShoppingLocationEntry
			}
		},
		func(confirm bool) {
			if confirm {
				newShoppingList.Name = newShoppingLocationEntry.Text
				newDocItem.Text = newShoppingList.Name
				a.tabs.Refresh()
			} else {
				a.tabs.Remove(newDocItem)
				for index, value := range a.shoppingLists {
					if value == newShoppingList {
						a.deleteShoppingList(index, value)
					}
				}
			}
		}, 
		a.win
	)
	a.win.Canvas().Focus(newShoppingLocationEntry)
	
	return newDocItem
}

func (a *appData) BuildTabItem(sl *shoppingList) *container.TabItem {
	displayCheckedItem := true
	filter := ""

	sl.list = widget.NewList(
		func() int {
			if filter == "" && displayCheckedItem {
				return len(sl.Items)
			}

			count := 0

			for _, i := range sl.Items {
				if i.shouldFilter(filter, displayCheckedItem) {
					continue
				}
				count++
			}
			return count
		},

		func() fyne.CanvasObject {
			return widget.NewCheck("test", func(b bool) {})
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			a.setFilteredItem(
				sl, 
				filter, 
				displayCheckedItem, 
				lii,
				co)
		}
	)

	func (a *appData) SetFilteredItem(
		sl *shoppingList, 
		filter strong, 
		displayCheckedItem bool, 
		index widget.ListItemID,
		co fyne.CanvasObject
	) {
		var pos widget.ListItemID
		for realIndex, i := range sl.Items {
			if i.shouldFilter(filter, displayCheckedItem) {
				continue
			}
			pos++

			if pos-1 == index {
				c := co.(*widget.Check)
				c.Text = i.What
				c.Checked = i.Checked
				c.OnChanged = func(b bool) {
					sl.Items[realIndex].Checked = b
				}
				c.Refresh()
				return
			}
		}
	}

	var toolbar *widget.Toolbar
	var visibilityAction *widget.ToolbarAction

	visibilityAction = widget.NewToolbarAction(
		theme.VisibilityIcon(), 
		func() {
			if displayCheckedItem {
				visibilityAction.SetIcon(theme.VisibilityOffIcon())
				displayCheckedItem = false
			} else {
				visibilityAction.SetIcon(theme.VisibilityIcon())
				displayChekckedItem = true
			}
			toolbar.Refresh()
			sl.list.Refresh()
	})

	toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), a.addItem(sl)),
		visibilityAction,
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(
			theme.ContentClearIcon(), 
			func() {
				keepItem := []item{}
				for _, item := range sl.Items {
					if item.Checked {
						keepItem = append(keepItem, item)	
					}
				}
				sl.Items = keepItem
				sl.list.Refresh()
		}),
	)

	func (a *appData) AddItem(sl *shoppingList) func() {
		return func() {
			newItemEntry := widget.NewEntry()
			dialog.ShowForm(
				"New shopping item", 
				"Create", 
				"Cancel",
				[]*widget.FormItem{{Text: "Name", Widget: newItemEntry}}, func(confirm bool) {
					if confirm {
						sl.Items = append(sl.Items, item{What: newItemEntry.Text})
						sl.list.Refresh()
					}
				}, 
				a.win)
			a.win.Canvas().Focus(newItemEntry)
		}
	}

	sl.filterEntry = widget.NewEntry()
	sl.filterEntry.OnChanged = func(s string) {
		filter = s
		sl.list.Refresh()
	}

	return container.NewTabItem(
		sl.Name, 
		container.NewBorder(
			container.NewBorder(
				nil, 
				nil, 
				widget.NewLabel("Filter"), 
				nil, 
				sl.filterEntry
			),
		toolbar,
		nil, 
		nil,
		sl.list),
	)
}