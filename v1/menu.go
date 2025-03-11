package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func MenuGroup(label string, items ...*fyne.MenuItem) *fyne.Menu {
	return fyne.NewMenu(label, items...)
}

func MenuItem(label string, action func()) *fyne.MenuItem {
	return fyne.NewMenuItem(label, action)
}

func MenuItemWithIcon(label string, icon string, action func()) *fyne.MenuItem {
	res := theme.Icon(fyne.ThemeIconName(icon))
	if res == nil {
		return MenuItem(label, action)
	}

	item := fyne.NewMenuItem(label, action)
	item.Icon = res

	return item
}

func MenuItemSeparator() *fyne.MenuItem {
	return fyne.NewMenuItemSeparator()
}
