package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type State struct {
	binds map[string]binding.DataItem
	lists map[string]binding.UntypedList
}

func NewState() *State {
	return &State{
		binds: make(map[string]binding.DataItem),
		lists: make(map[string]binding.UntypedList),
	}
}

func (state *State) Bool(name string, value bool) *Reactive[bool] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		container = binding.NewBool()
	}

	bind := Reactive[bool]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) Bytes(name string, value []byte) *Reactive[[]byte] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		container = binding.NewBytes()
	}

	bind := Reactive[[]byte]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) Float(name string, value float64) *Reactive[float64] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		container = binding.NewFloat()
	}

	bind := Reactive[float64]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) Int(name string, value int) *Reactive[int] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		container = binding.NewInt()
	}

	bind := Reactive[int]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) String(name string, value string) *Reactive[string] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		container = binding.NewString()
	}

	bind := Reactive[string]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) URI(name string, value fyne.URI) *Reactive[fyne.URI] {
	var container binding.DataItem
	if _, ok := state.binds[name]; ok {
		container = state.binds[name]
	} else {
		binding.NewURI()
	}

	bind := Reactive[fyne.URI]{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}

func (state *State) List(name string, value []any) *ReactiveList {
	var container binding.UntypedList
	if _, ok := state.binds[name]; ok {
		container = state.lists[name]
	} else {
		container = binding.NewUntypedList()
	}

	bind := ReactiveList{container: container}
	bind.Set(value)
	state.binds[name] = bind.container
	return &bind
}
