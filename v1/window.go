package reago

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/fsnotify/fsnotify"
)

var mainApp fyne.App = nil
var mainWindow *Window = nil

func init() {
	mainApp = app.New()
}

type Window struct {
	w fyne.Window
}

func NewWindow(title string, width float32, height float32) *Window {
	window := &Window{}
	window.w = mainApp.NewWindow(title)
	window.w.Resize(fyne.NewSize(width, height))
	return window
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

func (window *Window) Show(d *DOM) {
	window.w.SetContent(d.GetRoot())

	if mainWindow == nil {
		mainWindow = window
		window.w.ShowAndRun()
	} else {
		window.w.Show()
	}
}

func (window *Window) Watch(path string, callback func(*DOM)) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	d := NewDOM()
	d.FileTemplate(absPath)
	callback(d)

	go watchFile(absPath, func(event fsnotify.Event) {
		println("File " + path + " changed")

		d = NewDOM()
		d.FileTemplate(absPath)
		callback(d)

		window.w.SetContent(d.GetRoot())
	})

	window.Show(d)
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
