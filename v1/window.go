package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var mainApp fyne.App = nil
var mainWindow *Window = nil

func init() {
	mainApp = app.New()
}

type Window struct {
	w        fyne.Window
	menuRefs map[string]*fyne.MenuItem
}

func NewWindow(title string, width float32, height float32) *Window {
	window := &Window{}
	window.w = mainApp.NewWindow(title)
	window.w.Resize(fyne.NewSize(width, height))
	return window
}

func (window *Window) Show(d *DOM) {
	window.w.SetContent(d.GetRoot())

	if mainWindow == nil {
		mainWindow = window
		window.w.ShowAndRun()
	} else {
		window.w.Show()
	}
}

func (window *Window) OnFileDropped(callback func([]string)) {
	window.w.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		var paths []string
		for _, uri := range uris {
			paths = append(paths, uri.Path())
		}
		callback(paths)
	})
}

func (window *Window) GetTitle() string {
	return window.w.Title()
}

func (window *Window) SetTitle(value string) {
	window.w.SetTitle(value)
}

func (window *Window) IsFullScreen() bool {
	return window.w.FullScreen()
}

func (window *Window) SetFullScreen(value bool) {
	window.w.SetFullScreen(value)
}

func (window *Window) Resize(width float32, height float32) {
	window.w.Resize(fyne.NewSize(width, height))
}

func (window *Window) RequestFocus() {
	window.w.RequestFocus()
}

func (window *Window) IsFixedSize() bool {
	return window.w.FixedSize()
}

func (window *Window) SetFixedSize(value bool) {
	window.w.SetFixedSize(value)
}

func (window *Window) CenterOnScreen() {
	window.w.CenterOnScreen()
}

func (window *Window) IsPadded() bool {
	return window.w.Padded()
}

func (window *Window) SetPadded(value bool) {
	window.w.SetPadded(value)
}

func (window *Window) OnClosed(callback func()) {
	window.w.SetOnClosed(callback)
}

func (window *Window) OnBeforeClose(callback func()) {
	window.w.SetCloseIntercept(func() {
		callback()
		window.Close()
	})
}

func (window *Window) Close() {
	window.w.Close()
}

func (window *Window) GetClipboard() string {
	return window.w.Clipboard().Content()
}

func (window *Window) SetClipboard(value string) {
	window.w.Clipboard().SetContent(value)
}

func (window *Window) SetMainMenu(menus ...*fyne.Menu) {
	window.menuRefs = make(map[string]*fyne.MenuItem)

	for _, menu := range menus {
		if menu != nil {
			for _, item := range menu.Items {
				if item != nil {
					window.menuRefs[menu.Label+"/"+item.Label] = item
				}
			}
		}
	}

	window.w.SetMainMenu(fyne.NewMainMenu(menus...))
}

func (window *Window) GetMenuItem(id string) *fyne.MenuItem {
	return window.menuRefs[id]
}

func (window *Window) MainMenu() *fyne.MainMenu {
	return window.w.MainMenu()
}
