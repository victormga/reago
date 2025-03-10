package reago

import (
	"fyne.io/fyne/v2"
)

type DOM struct {
	root  fyne.CanvasObject
	refs  map[string]fyne.CanvasObject
	state *State
}

func NewDOM() *DOM {
	dom := &DOM{
		refs:  make(map[string]fyne.CanvasObject),
		state: NewState(),
	}
	return dom
}

func (dom *DOM) UseState() *State {
	return dom.state
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
