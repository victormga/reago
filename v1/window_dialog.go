package reago

import (
	"errors"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (window *Window) DlgColorPicker(title string, message string, callback func(color.Color)) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.NewColorPicker("Pick a Color", "Choose a color from the palette", func(c color.Color) {
		callback(c)
	}, window.w)
}

func (window *Window) DlgFileOpen(callback func(string)) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}

		callback(f.URI().String())
	}, window.w)
}

func (window *Window) DlgFileSave(callback func(string)) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.ShowFileSave(func(f fyne.URIWriteCloser, err error) {
		if err != nil {
			fmt.Println(err)
			return
		}
		callback(f.URI().String())
	}, window.w)
}

func (window *Window) DlgFolderOpen(callback func(string)) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.ShowFolderOpen(func(f fyne.ListableURI, err error) {
		if err != nil {
			return
		}
		callback(f.String())
	}, window.w)
}

func (window *Window) DlgAlert(title string, message string) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	minSize := fyne.NewSize(200, 100)
	minRect := canvas.NewRectangle(color.Transparent)
	minRect.SetMinSize(minSize)

	messageLabel := widget.NewLabel(message)
	content := container.NewStack(minRect, container.NewCenter(messageLabel))

	dlg := dialog.NewCustom(title, "OK", content, window.w)
	dlg.Show()
}

func (window *Window) DlgError(title string, message string) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.ShowError(errors.New(message), window.w)
}

func (window *Window) DlgConfirm(title string, message string, callback func(bool)) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	dialog.ShowConfirm(title, message, func(b bool) {
		callback(b)
	}, window.w)
}

func (window *Window) DlgProgress(title string, message string, total int, callback func()) *dialog.CustomDialog {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return nil
	}

	dlg := dialog.NewCustomWithoutButtons("Loading", widget.NewProgressBarInfinite(), window.w)
	dlg.Show()

	return dlg
}

func (window *Window) DlgModal(title string, content string) {
	if window.w == nil {
		fmt.Println("Error: can't show dialog on a fragment")
		return
	}

	fragment := NewDOM()
	fragment.LoadString(content)
	dialog.ShowCustom(title, "Close", fragment.root, window.w)
}
