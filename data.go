package main

func (a *appData) NewShoppingList(name string) (*shoppingList, error) {
	newShoppingList := &shoppingList{Name: name}
	a.shoppingLists = append(a.shoppingLists, newShoppingList)
	return newShoppingList, nil
}

func (a *appData) DeleteShoppingList(index int, sl *shoppingList) error { // error in right place?
	if index < len(a.shoppingLists)-1 {
		a.shoppingLists[index] = a.shoppingLists[len(a.shoppingLists)-1]
	}
	a.shoppingLists = a.shoppingLists[:len(a.shoppingLists)-1]

	return nil
}
