package reago

import (
	"fyne.io/fyne/v2"
)

type DOM struct {
	root      fyne.CanvasObject
	refs      map[string]fyne.CanvasObject
	state     *State
	callbacks map[string]func(*XMLNode)
}

func NewDOM() *DOM {
	dom := &DOM{
		refs:      make(map[string]fyne.CanvasObject),
		state:     NewState(),
		callbacks: make(map[string]func(*XMLNode)),
	}
	return dom
}

func (dom *DOM) Clone() *DOM {
	clone := NewDOM()
	for name, reactive := range dom.state.binds {
		clone.state.binds[name] = reactive // might have to clone each one?
	}
	for name, callback := range dom.callbacks {
		clone.callbacks[name] = callback
	}
	return clone
}

func (dom *DOM) UseState() *State {
	return dom.state
}

func (dom *DOM) UseCallback(name string, callback func(node *XMLNode)) {
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
