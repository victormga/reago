package reago

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/fsnotify/fsnotify"
)

type DOM struct {
	root      *fyne.Container
	refs      map[string]fyne.CanvasObject
	state     *State
	callbacks map[string]func(*XMLNode)
}

func NewDOM() *DOM {
	dom := &DOM{
		root:      container.NewStack(nil),
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

func (dom *DOM) FileTemplate(path string, watch bool) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	dom.Template(readXMLFile(absPath))

	if watch {
		go watchFile(absPath, func(event fsnotify.Event) {
			log.Println("File " + path + " changed")
			dom.Template(readXMLFile(absPath))
		})
	}
}

func (dom *DOM) Template(content string) {
	dom.refs = make(map[string]fyne.CanvasObject)
	dom.root.Objects = []fyne.CanvasObject{Parser.ParseXML(content, dom)}
	dom.root.Refresh()
}

func (dom *DOM) AppendTo(parent *fyne.Container) {
	parent.Add(dom.root)
}

func readXMLFile(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func watchFile(filename string, callback func(event fsnotify.Event)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch the directory that contains the file.
	dir := filepath.Dir(filename)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			// Compare absolute paths to ensure a match.
			eventPath, err := filepath.Abs(event.Name)
			if err != nil {
				continue
			}
			if eventPath == filename {
				callback(event)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}
