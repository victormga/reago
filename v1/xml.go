package reago

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type XMLNode struct {
	XMLName xml.Name   `xml:""`
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",chardata"`
	Nodes   []XMLNode  `xml:",any"`
}

func (node *XMLNode) GetTag() string {
	return strings.ToLower(node.XMLName.Local)
}

func (node *XMLNode) HasAttr(name string) bool {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return true
		}
	}
	return false
}

func (node *XMLNode) GetAttr(name string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func (node *XMLNode) GetAttrInt(name string) int {
	value := node.GetAttr(name)
	if value != "" {
		if result, err := strconv.Atoi(value); err == nil {
			return result
		}
	}
	return 0
}

func (node *XMLNode) GetAttrInt32(name string) int32 {
	return int32(node.GetAttrInt(name))
}

func (node *XMLNode) GetAttrFloat(name string) float64 {
	value := node.GetAttr(name)
	if value != "" {
		if result, err := strconv.ParseFloat(value, 64); err == nil {
			return result
		}
	}
	return 0
}

func (node *XMLNode) GetAttrFloat32(name string) float32 {
	return float32(node.GetAttrFloat(name))
}

func (node *XMLNode) GetAttrBool(name string) bool {
	value := node.GetAttr(name)
	return value == "true" || value == "1"
}

func (node *XMLNode) GetContent() string {
	return strings.TrimSpace(node.Content)
}

func (node *XMLNode) HasBind(name string) bool {
	for _, attr := range node.Attrs {
		if attr.Name.Space == "bind" && attr.Name.Local == name {
			return true
		}
	}
	return false
}

func (node *XMLNode) GetBind(name string) string {
	for _, attr := range node.Attrs {
		if attr.Name.Space == "bind" && attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func (node *XMLNode) BindCallback(name string, target *DOM) func() {
	if node.HasBind(name) {
		bind := node.GetBind(name)
		return func() {
			if callback, ok := target.callbacks[bind]; ok {
				callback(node)
			}
		}
	}
	return nil
}

func (node *XMLNode) BindContent(target *DOM, update func(string)) func(string) {
	value := node.GetContent()
	bind := node.GetBind("content")

	if bind == "" && node.HasBind("content") {
		tpl := NewTplParser(value)
		for _, bind := range tpl.GetBinds() {
			target.state.GetString(bind).OnChange(func(_ string) {
				update(tpl.Render(target.state))
			})
		}

		update(tpl.Render(target.state))
		return nil
	}

	return bindToState(value, bind, target.state, target.state.GetString, update)
}

func (node *XMLNode) BindList(name string, target *DOM, update func([]any)) func([]any) {
	bind := node.GetBind(name)
	if bind != "" {
		reactive := target.state.GetList(bind)
		reactive.OnChange(update)
		update(reactive.Get())
		return reactive.Set
	}

	return nil
}

func (node *XMLNode) BindString(name string, target *DOM, update func(string)) func(string) {
	value := node.GetAttr(name)
	bind := node.GetBind(name)
	return bindToState(value, bind, target.state, target.state.GetString, update)
}

func (node *XMLNode) BindInt(name string, target *DOM, update func(int)) func(int) {
	value := node.GetAttrInt(name)
	bind := node.GetBind(name)
	return bindToState(value, bind, target.state, target.state.GetInt, update)
}

func (node *XMLNode) BindFloat(name string, target *DOM, update func(float64)) func(float64) {
	value := node.GetAttrFloat(name)
	bind := node.GetBind(name)
	return bindToState(value, bind, target.state, target.state.GetFloat, update)
}

func (node *XMLNode) BindBool(name string, target *DOM, update func(bool)) func(bool) {
	value := node.GetAttrBool(name)
	bind := node.GetBind(name)
	return bindToState(value, bind, target.state, target.state.GetBool, update)
}

func bindToState[T comparable](
	value T,
	bind string,
	state *State,
	getter func(string) *Reactive[T],
	update func(T),
) func(T) {
	if !isZero(value) {
		update(value)
	}

	if bind != "" {
		initialized := state.Has(bind)

		reactive := getter(bind)
		reactive.OnChange(update)

		if !initialized {
			reactive.Set(value)
		}

		return reactive.Set
	}

	return nil
}

func isZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func (node *XMLNode) GetPadding() (float32, float32, float32, float32) {
	var top, bottom, left, right float32

	if node.HasAttr("padding") {
		all := node.GetAttrFloat32("padding")
		top, bottom, left, right = all, all, all, all
	}
	if node.HasAttr("padding-vertical") {
		vertical := node.GetAttrFloat32("padding-vertical")
		top, bottom = vertical, vertical
	}
	if node.HasAttr("padding-horizontal") {
		horizontal := node.GetAttrFloat32("padding-horizontal")
		left, right = horizontal, horizontal
	}
	if node.HasAttr("padding-top") {
		top = node.GetAttrFloat32("padding-top")
	}
	if node.HasAttr("padding-bottom") {
		bottom = node.GetAttrFloat32("padding-bottom")
	}
	if node.HasAttr("padding-left") {
		left = node.GetAttrFloat32("padding-left")
	}
	if node.HasAttr("padding-right") {
		right = node.GetAttrFloat32("padding-right")
	}

	return top, bottom, left, right
}
