package reago

import (
	"encoding/xml"

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

func (dom *DOM) LoadFile(path string) {
	content := readXMLFile(path)
	dom.LoadString(content)
}

func (dom *DOM) LoadString(content string) {
	dom.refs = make(map[string]fyne.CanvasObject)

	var xmlRoot XMLNode
	if err := xml.Unmarshal([]byte(content), &xmlRoot); err != nil {
		dom.LoadString(`
			<row>
				<text>Error parsing XML</text>
				<text style="bold" color="red">` + err.Error() + `</text>
			</row>
		`)
		return
	}

	dom.root = *dom.parseNode(&xmlRoot)
}

func (dom *DOM) AppendTo(parent *fyne.Container) {
	parent.Add(dom.root)
}
