package reago

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type ReactiveList struct {
	container binding.UntypedList
}

func (r *ReactiveList) Set(value []any) {
	r.container.Set(value)
}

type Reactive[T any] struct {
	container binding.DataItem
	listeners []binding.DataListener
}

func (reactive *Reactive[T]) Set(value T) {
	if reflect.DeepEqual(reactive.Get(), value) {
		return
	}

	switch v := reactive.container.(type) {
	case binding.Bool:
		if val, ok := any(value).(bool); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to bool failed")
		}
	case binding.Bytes:
		if val, ok := any(value).([]byte); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to []byte failed")
		}
	case binding.Float:
		if val, ok := any(value).(float64); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to float64 failed")
		}
	case binding.Int:
		if val, ok := any(value).(int); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to int failed")
		}
	case binding.String:
		if val, ok := any(value).(string); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to string failed")
		}
	case binding.URI:
		if val, ok := any(value).(fyne.URI); ok {
			v.Set(val)
		} else {
			fmt.Println("type assertion to fyne.URI failed")
		}
	default:
		fmt.Println("unsupported binding type")
	}
}

func (reactive *Reactive[T]) Get() T {
	var value T
	var err error

	switch v := reactive.container.(type) {
	case binding.Bool:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	case binding.Bytes:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	case binding.Float:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	case binding.Int:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	case binding.String:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	case binding.URI:
		val, e := v.Get()
		value = any(val).(T)
		err = e
	default:
		err = fmt.Errorf("unsupported binding type")
	}

	if err != nil {
		var zero T
		return zero
	}

	return value
}

func (reactive *Reactive[T]) AddListener(callback func(T)) {
	listener := binding.NewDataListener(func() {
		callback(reactive.Get())
	})
	reactive.listeners = append(reactive.listeners, listener)
	reactive.container.AddListener(listener)
}

func (reactive *Reactive[T]) RemoveAllListeners(callback func(T)) {
	for _, listener := range reactive.listeners {
		reactive.container.RemoveListener(listener)
	}
	reactive.listeners = nil
}
