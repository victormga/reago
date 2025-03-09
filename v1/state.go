package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type State struct {
	binds map[string]IReactive
}

func NewState() *State {
	return &State{
		binds: make(map[string]IReactive),
	}
}

func (state *State) Has(name string) bool {
	_, ok := state.binds[name]
	return ok
}

func (state *State) GetBool(name string) *Reactive[bool] {
	var bind *Reactive[bool]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[bool])
	} else {
		bind = &Reactive[bool]{container: binding.NewBool()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) Bool(name string, value bool) *Reactive[bool] {
	bind := state.GetBool(name)
	bind.Set(value)
	return bind
}

func (state *State) GetBytes(name string) *Reactive[[]byte] {
	var bind *Reactive[[]byte]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[[]byte])
	} else {
		bind = &Reactive[[]byte]{container: binding.NewBytes()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) Bytes(name string, value []byte) *Reactive[[]byte] {
	bind := state.GetBytes(name)
	bind.Set(value)
	return bind
}

func (state *State) GetFloat(name string) *Reactive[float64] {
	var bind *Reactive[float64]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[float64])
	} else {
		bind = &Reactive[float64]{container: binding.NewFloat()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) Float(name string, value float64) *Reactive[float64] {
	bind := state.GetFloat(name)
	bind.Set(value)
	return bind
}

func (state *State) GetInt(name string) *Reactive[int] {
	var bind *Reactive[int]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[int])
	} else {
		bind = &Reactive[int]{container: binding.NewInt()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) Int(name string, value int) *Reactive[int] {
	bind := state.GetInt(name)
	bind.Set(value)
	return bind
}

func (state *State) GetString(name string) *Reactive[string] {
	var bind *Reactive[string]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[string])
	} else {
		bind = &Reactive[string]{container: binding.NewString()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) String(name string, value string) *Reactive[string] {
	bind := state.GetString(name)
	bind.Set(value)
	return bind
}

func (state *State) GetURI(name string) *Reactive[fyne.URI] {
	var bind *Reactive[fyne.URI]
	if reactive, ok := state.binds[name]; ok {
		bind = reactive.(*Reactive[fyne.URI])
	} else {
		bind = &Reactive[fyne.URI]{container: binding.NewURI()}
		state.binds[name] = bind
	}
	return bind
}

func (state *State) URI(name string, value fyne.URI) *Reactive[fyne.URI] {
	bind := state.GetURI(name)
	bind.Set(value)
	return bind
}
