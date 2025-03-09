package reago

import (
	"fyne.io/fyne/v2"
)

type DOM struct {
	root      fyne.CanvasObject
	refs      map[string]fyne.CanvasObject
	state     *State
	callbacks map[string]func() // TODO: this needs to be able to receive arguments (example, when input changes)
}

func NewDOM() *DOM {
	dom := &DOM{
		refs:      make(map[string]fyne.CanvasObject),
		state:     NewState(),
		callbacks: make(map[string]func()),
	}
	return dom
}

func (dom *DOM) UseState() *State {
	return dom.state
}

func (dom *DOM) UseCallback(name string, callback func()) {
	dom.callbacks[name] = callback
}

func (dom *DOM) FileTemplate(path string) {
	content := readXMLFile(path)
	dom.Template(content)
}

func (dom *DOM) Template(content string) {
	dom.refs = make(map[string]fyne.CanvasObject)
	dom.root = Parser.ParseXML(content, dom)
}

func (dom *DOM) AppendTo(parent *fyne.Container) {
	parent.Add(dom.root)
}
