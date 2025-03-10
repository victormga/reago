package reago

import (
	"encoding/xml"
	"errors"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type iParser struct {
	tags       map[string]func(*XMLNode, *DOM) fyne.CanvasObject
	components map[string]func(*XMLNode, *DOM) string
}

var Parser = iParser{
	tags:       make(map[string]func(*XMLNode, *DOM) fyne.CanvasObject),
	components: make(map[string]func(*XMLNode, *DOM) string),
}

func (parser *iParser) RegisterTag(tag string, handler func(*XMLNode, *DOM) fyne.CanvasObject) {
	parser.tags[tag] = handler
}

func (parser *iParser) RegisterComponent(name string, component func(*XMLNode, *DOM) string) {
	parser.components[name] = component
}

func (parser *iParser) ParseXML(content string, target *DOM) fyne.CanvasObject {
	var xmlRoot XMLNode
	if err := xml.Unmarshal([]byte(content), &xmlRoot); err != nil {
		return widget.NewLabelWithStyle(
			"component_error: "+err.Error(),
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true, Monospace: true},
		)
	}

	return parser.ParseNode(&xmlRoot, target)
}

func (parser *iParser) ParseNode(node *XMLNode, target *DOM) fyne.CanvasObject {
	if node == nil {
		return nil
	}

	var obj fyne.CanvasObject

	tag := node.GetTag()
	if handler, ok := parser.tags[tag]; ok {
		obj = handler(node, target)
	} else if component, ok := parser.components[tag]; ok {
		obj = parser.ParseXML(component(node, target), target)
	} else {
		obj = widget.NewLabel("<unknown tag: " + tag + ">")
	}

	id := node.GetAttr("id")
	if id != "" {
		target.refs[id] = obj
	}

	node.BindBool("hidden", target, func(value bool) {
		if value {
			obj.Hide()
		} else {
			obj.Show()
		}
	})

	return obj
}

func (parser *iParser) ParseChildren(node *XMLNode, target *DOM) []fyne.CanvasObject {
	var children []fyne.CanvasObject
	for _, child := range node.Nodes {
		children = append(children, Parser.ParseNode(&child, target))
	}
	return children
}

func ParseColor(str string) (color.Color, error) {
	if str == "" {
		return nil, errors.New("color is empty")
	}

	if str[0] == '#' {
		str := strings.TrimPrefix(str, "#")
		if len(str) == 3 {
			str = string([]byte{str[0], str[0], str[1], str[1], str[2], str[2]})
		}
		if len(str) != 6 {
			return nil, errors.New("invalid hex color format")
		}

		r, err := strconv.ParseUint(str[0:2], 16, 8)
		if err != nil {
			return nil, err
		}
		g, err := strconv.ParseUint(str[2:4], 16, 8)
		if err != nil {
			return nil, err
		}
		b, err := strconv.ParseUint(str[4:6], 16, 8)
		if err != nil {
			return nil, err
		}

		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
	} else {
		c, ok := colornames.Map[strings.ToLower(str)]
		if !ok {
			return nil, errors.New("unknown color name")
		}

		return c, nil
	}
}
