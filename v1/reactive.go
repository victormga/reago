package reago

import (
	"fmt"
	"log"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type IReactive interface{}

type Reactive[T any] struct {
	IReactive
	container binding.DataItem
	listeners []binding.DataListener
	getter    func() (T, error)
	setter    func(T) error
}

func NewReactive[T any](container binding.DataItem) *Reactive[T] {
	r := &Reactive[T]{container: container}

	switch v := container.(type) {
	case binding.Bool:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			boolVal, ok := any(val).(bool)
			if !ok {
				return fmt.Errorf("type assertion to bool failed")
			}
			return v.Set(boolVal)
		}
	case binding.Bytes:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			bytesVal, ok := any(val).([]byte)
			if !ok {
				return fmt.Errorf("type assertion to []byte failed")
			}
			return v.Set(bytesVal)
		}
	case binding.Float:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			floatVal, ok := any(val).(float64)
			if !ok {
				return fmt.Errorf("type assertion to float64 failed")
			}
			return v.Set(floatVal)
		}
	case binding.Int:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			intVal, ok := any(val).(int)
			if !ok {
				return fmt.Errorf("type assertion to int failed")
			}
			return v.Set(intVal)
		}
	case binding.String:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			strVal, ok := any(val).(string)
			if !ok {
				return fmt.Errorf("type assertion to string failed")
			}
			return v.Set(strVal)
		}
	case binding.URI:
		r.getter = func() (T, error) {
			val, err := v.Get()
			return any(val).(T), err
		}
		r.setter = func(val T) error {
			uriVal, ok := any(val).(fyne.URI)
			if !ok {
				return fmt.Errorf("type assertion to fyne.URI failed")
			}
			return v.Set(uriVal)
		}
	default:
		panic("unsupported binding type")
	}

	return r
}

func (r *Reactive[T]) Set(value T) {
	current, err := r.getter()
	if err != nil {
		fmt.Println("get error:", err)
		return
	}

	if reflect.DeepEqual(current, value) {
		return
	}

	if err := r.setter(value); err != nil {
		fmt.Println("set error:", err)
	}
}

func (r *Reactive[T]) Get() T {
	val, err := r.getter()
	if err != nil {
		fmt.Println("get error:", err)
	}
	return val
}

func (r *Reactive[T]) OnChange(callback func(T)) {
	listener := binding.NewDataListener(func() {
		callback(r.Get())
	})
	r.listeners = append(r.listeners, listener)
	r.container.AddListener(listener)
}

func (r *Reactive[T]) ClearListeners() {
	for _, listener := range r.listeners {
		r.container.RemoveListener(listener)
	}
	r.listeners = nil
}

type ReactiveList[T any] struct {
	IReactive
	container binding.UntypedList
	listeners []binding.DataListener
}

func NewReactiveList[T any](container binding.UntypedList) *ReactiveList[T] {
	return &ReactiveList[T]{container: container}
}

func (rl *ReactiveList[T]) Append(value T) {
	if err := rl.container.Append(value); err != nil {
		log.Println("ReactiveList Append error:", err)
	}
}

func (rl *ReactiveList[T]) Get() []T {
	items, err := rl.container.Get()
	if err != nil {
		log.Println("ReactiveList Get error:", err)
		return []T{}
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		typed, ok := item.(T)
		if !ok {
			log.Println("ReactiveList Get type assertion failed for item:", item)
			continue
		}
		result = append(result, typed)
	}
	return result
}

func (rl *ReactiveList[T]) GetValue(index int) T {
	item, err := rl.container.GetValue(index)
	if err != nil {
		log.Println("ReactiveList GetValue error:", err)
		var zero T
		return zero
	}
	typed, ok := item.(T)
	if !ok {
		log.Println("ReactiveList GetValue type assertion failed for item:", item)
		var zero T
		return zero
	}
	return typed
}

func (rl *ReactiveList[T]) Prepend(value T) {
	if err := rl.container.Prepend(value); err != nil {
		log.Println("ReactiveList Prepend error:", err)
	}
}

func (rl *ReactiveList[T]) Remove(value T) {
	if err := rl.container.Remove(value); err != nil {
		log.Println("ReactiveList Remove error:", err)
	}
}

func (rl *ReactiveList[T]) Set(list []T) {
	anyList := make([]any, len(list))
	for i, v := range list {
		anyList[i] = v
	}
	if err := rl.container.Set(anyList); err != nil {
		log.Println("ReactiveList Set error:", err)
	}
}

func (rl *ReactiveList[T]) SetValue(index int, value T) {
	if err := rl.container.SetValue(index, value); err != nil {
		log.Println("ReactiveList SetValue error:", err)
	}
}

func (r *ReactiveList[T]) OnChange(callback func([]T)) {
	listener := binding.NewDataListener(func() {
		callback(r.Get())
	})
	r.listeners = append(r.listeners, listener)
	r.container.AddListener(listener)
}

func (r *ReactiveList[T]) ClearListeners() {
	for _, listener := range r.listeners {
		r.container.RemoveListener(listener)
	}
	r.listeners = nil
}
