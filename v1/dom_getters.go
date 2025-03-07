package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (dom *DOM) GetRoot() fyne.CanvasObject {
	return dom.root
}

func (dom *DOM) GetRow(id string) *fyne.Container {
	return cast[*fyne.Container](dom.refs, id)
}

func (dom *DOM) GetCol(id string) *fyne.Container {
	return cast[*fyne.Container](dom.refs, id)
}

func (dom *DOM) GetCenter(id string) *fyne.Container {
	return cast[*fyne.Container](dom.refs, id)
}

func (dom *DOM) GetForm(id string) *widget.Form {
	return cast[*widget.Form](dom.refs, id)
}

func (dom *DOM) GetGrid(id string) *fyne.Container {
	return cast[*fyne.Container](dom.refs, id)
}

func (dom *DOM) GetLayout(id string) *fyne.Container {
	return cast[*fyne.Container](dom.refs, id)
}

func (dom *DOM) GetLabel(id string) *widget.Label {
	return cast[*widget.Label](dom.refs, id)
}

func (dom *DOM) GetText(id string) *canvas.Text {
	return cast[*canvas.Text](dom.refs, id)
}

func (dom *DOM) GetButton(id string) *widget.Button {
	return cast[*widget.Button](dom.refs, id)
}

func (dom *DOM) GetInput(id string) *widget.Entry {
	return cast[*widget.Entry](dom.refs, id)
}

func (dom *DOM) GetTextarea(id string) *widget.Entry {
	return cast[*widget.Entry](dom.refs, id)
}

func (dom *DOM) GetCheckbox(id string) *widget.Check {
	return cast[*widget.Check](dom.refs, id)
}

func (dom *DOM) GetRadio(id string) *widget.RadioGroup {
	return cast[*widget.RadioGroup](dom.refs, id)
}

func (dom *DOM) GetAccordion(id string) *widget.Accordion {
	return cast[*widget.Accordion](dom.refs, id)
}

func (dom *DOM) GetActivity(id string) *widget.Activity {
	return cast[*widget.Activity](dom.refs, id)
}

func (dom *DOM) GetCard(id string) *widget.Card {
	return cast[*widget.Card](dom.refs, id)
}

func (dom *DOM) GetA(id string) *widget.Hyperlink {
	return cast[*widget.Hyperlink](dom.refs, id)
}

func (dom *DOM) GetIcon(id string) *widget.Icon {
	return cast[*widget.Icon](dom.refs, id)
}

func (dom *DOM) GetProgress(id string) *widget.ProgressBar {
	return cast[*widget.ProgressBar](dom.refs, id)
}

func (dom *DOM) GetLoader(id string) *widget.ProgressBarInfinite {
	return cast[*widget.ProgressBarInfinite](dom.refs, id)
}

func (dom *DOM) GetMarkdown(id string) *widget.RichText {
	return cast[*widget.RichText](dom.refs, id)
}

func (dom *DOM) GetSelect(id string) *widget.Select {
	return cast[*widget.Select](dom.refs, id)
}

func (dom *DOM) GetCombobox(id string) *widget.SelectEntry {
	return cast[*widget.SelectEntry](dom.refs, id)
}

func (dom *DOM) GetHr(id string) *widget.Separator {
	return cast[*widget.Separator](dom.refs, id)
}

func (dom *DOM) GetSlider(id string) *widget.Slider {
	return cast[*widget.Slider](dom.refs, id)
}

func (dom *DOM) GetCode(id string) *widget.TextGrid {
	return cast[*widget.TextGrid](dom.refs, id)
}

func (dom *DOM) GetToolbar(id string) *widget.Toolbar {
	return cast[*widget.Toolbar](dom.refs, id)
}

// TODO: make this work
func (dom *DOM) GetToolbarAction(id string) *widget.ToolbarAction {
	return cast[*widget.ToolbarAction](dom.refs, id)
}

func (dom *DOM) GetVirtualList(id string) *widget.List {
	return cast[*widget.List](dom.refs, id)
}

func (dom *DOM) GetTable(id string) *widget.Table {
	return cast[*widget.Table](dom.refs, id)
}

func (dom *DOM) GetTree(id string) *widget.Tree {
	return cast[*widget.Tree](dom.refs, id)
}

func (dom *DOM) GetTabs(id string) *container.AppTabs {
	return cast[*container.AppTabs](dom.refs, id)
}

func (dom *DOM) GetScroller(id string) *container.Scroll {
	return cast[*container.Scroll](dom.refs, id)
}

func (dom *DOM) GetSpacer(id string) *fyne.CanvasObject {
	return cast[*fyne.CanvasObject](dom.refs, id)
}

func cast[T any](refs map[string]fyne.CanvasObject, id string) T {
	if obj, exists := refs[id]; exists {
		if casted, ok := obj.(T); ok {
			return casted
		}
	}
	var zero T
	return zero
}
