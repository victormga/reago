package reago

import (
	"errors"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

var renderers map[string]func(*XMLNode, *DOM) fyne.CanvasObject

func init() {
	renderers = map[string]func(*XMLNode, *DOM) fyne.CanvasObject{
		/**
		 * <row>
		 */
		"row": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)
			return container.NewHBox(children...)
		},

		/**
		 * <col>
		 */
		"col": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)
			return container.NewVBox(children...)
		},

		/**
		 * <center>
		 */
		"center": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)
			return container.NewCenter(children...)
		},

		/**
		 * <form>
		 */
		"form": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			// TODO: circle back to this
			//children := dom.parseChildren(node)
			//return widget.NewForm(children...)
			return nil
		},

		/**
		 * <grid>
		 */
		"grid": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)

			columns := node.GetAttrInt("cols")
			if columns > 0 {
				return container.NewGridWithColumns(columns, children...)
			}

			rows := node.GetAttrInt("rows")
			if rows > 0 {
				return container.NewGridWithRows(rows, children...)
			}

			return container.NewGridWrap(
				fyne.NewSize(node.GetAttrFloat32("width"), node.GetAttrFloat32("height")),
				children...,
			)
		},

		/**
		 * <layout>
		 */
		"layout": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var top fyne.CanvasObject
			var bottom fyne.CanvasObject
			var left fyne.CanvasObject
			var right fyne.CanvasObject
			var content []fyne.CanvasObject

			for _, child := range node.Nodes {
				switch child.GetTag() {
				case "layout-top":
					top = *dom.parseNode(&child)
				case "layout-bottom":
					bottom = *dom.parseNode(&child)
				case "layout-left":
					left = *dom.parseNode(&child)
				case "layout-right":
					right = *dom.parseNode(&child)
				default:
					content = append(content, *dom.parseNode(&child))
				}
			}

			return container.NewBorder(top, bottom, left, right, content...)
		},

		/**
		 * <padding>
		 */
		"padding": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)

			var top, bottom, left, right float32
			if node.HasAttr("all") {
				all := node.GetAttrFloat32("all")
				top, bottom, left, right = all, all, all, all
			}
			if node.HasAttr("vertical") {
				vertical := node.GetAttrFloat32("vertical")
				top, bottom = vertical, vertical
			}
			if node.HasAttr("horizontal") {
				horizontal := node.GetAttrFloat32("horizontal")
				left, right = horizontal, horizontal
			}
			if node.HasAttr("top") {
				top = node.GetAttrFloat32("top")
			}
			if node.HasAttr("bottom") {
				bottom = node.GetAttrFloat32("bottom")
			}
			if node.HasAttr("left") {
				left = node.GetAttrFloat32("left")
			}
			if node.HasAttr("right") {
				right = node.GetAttrFloat32("right")
			}

			return container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				container.NewVBox(children...),
			)
		},

		/**
		 * <label>
		 */
		"label": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewLabel("")

			bindContent(node, dom, func(value string) {
				obj.SetText(value)
			})

			return obj
		},

		/**
		 * <text>
		 */
		"text": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := canvas.NewText("", theme.Color(theme.ColorNameForeground))

			bindContent(node, dom, func(value string) {
				obj.Text = value
				obj.Refresh()
			})

			bindString("color", node, dom, func(value string) {
				color, _ := parseColor(value)
				if color != nil {
					obj.Color = color
					obj.Refresh()
				}
			})

			bindFloat("size", node, dom, func(value float64) {
				obj.TextSize = float32(value)
				obj.Refresh()
			})

			bindString("style", node, dom, func(value string) {
				switch value {
				case "bold":
					obj.TextStyle = fyne.TextStyle{Bold: true}
				case "italic":
					obj.TextStyle = fyne.TextStyle{Italic: true}
				case "monospace":
					obj.TextStyle = fyne.TextStyle{Monospace: true}
				case "underline":
					obj.TextStyle = fyne.TextStyle{Underline: true}
				default:
					obj.TextStyle = fyne.TextStyle{}
				}
				obj.Refresh()
			})

			bindString("align", node, dom, func(value string) {
				switch value {
				case "center":
					obj.Alignment = fyne.TextAlignCenter
				case "right":
					obj.Alignment = fyne.TextAlignTrailing
				default:
					obj.Alignment = fyne.TextAlignLeading
				}
				obj.Refresh()
			})

			bindContent(node, dom, func(value string) {
				obj.Text = value
				obj.Refresh()
			})

			return obj
		},

		/**
		 * <button>
		 */
		"button": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			text := node.GetText()

			var obj *widget.Button

			icon := node.GetAttr("icon")
			if icon != "" {
				// TODO: this is probably wrong
				obj = widget.NewButtonWithIcon(text, fyne.NewStaticResource(icon, nil), func() {})
			} else {
				obj = widget.NewButton(text, func() {})
			}

			bindContent(node, dom, func(value string) {
				obj.SetText(value)
			})

			obj.OnTapped = bindCallback("click", node, dom)

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					obj.Disable()
				} else {
					obj.Enable()
				}
			})

			return obj
		},

		/**
		 * <input>
		 */
		"input": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var entry *widget.Entry

			switch node.GetAttr("type") {
			case "password":
				entry = widget.NewPasswordEntry()
			case "number":
				entry = widget.NewEntry()
				entry.Validator = validation.NewRegexp(`^\d+$`, "Only numbers are allowed")
			case "email":
				entry = widget.NewEntry()
				entry.Validator = validation.NewRegexp(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "Invalid email address")
			case "url":
				entry = widget.NewEntry()
				entry.Validator = validation.NewRegexp(`^https?://.+$`, "Invalid URL")
			default:
				entry = widget.NewEntry()
			}

			bindString("placeholder", node, dom, func(value string) {
				entry.SetPlaceHolder(value)
			})

			entry.OnChanged = bindString("value", node, dom, func(value string) {
				entry.SetText(value)
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					entry.Disable()
				} else {
					entry.Enable()
				}
			})

			// TODO:
			//entry.OnSubmitted = func(string) {}

			if node.HasAttr("validation") {
				msg := node.GetAttr("validation-message")
				if msg == "" {
					msg = "invalid"
				}
				entry.Validator = validation.NewRegexp(node.GetAttr("validation"), msg)
			}

			return entry
		},

		/**
		 * <textarea>
		 */
		"textarea": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			entry := widget.NewMultiLineEntry()

			bindString("placeholder", node, dom, func(value string) {
				entry.SetPlaceHolder(value)
			})

			entry.OnChanged = bindString("value", node, dom, func(value string) {
				entry.SetText(value)
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					entry.Disable()
				} else {
					entry.Enable()
				}
			})

			// validation
			if node.HasAttr("validation") {
				msg := node.GetAttr("validation-message")
				if msg == "" {
					msg = "invalid"
				}

				entry.Validator = validation.NewRegexp(node.GetAttr("validation"), msg)
			}

			return entry
		},

		/**
		 * <checkbox>
		 */
		"checkbox": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewCheck(node.GetAttr("label"), func(bool) {})

			bindString("label", node, dom, func(value string) {
				obj.Text = value
				obj.Refresh()
			})

			obj.OnChanged = bindBool("value", node, dom, func(value bool) {
				obj.SetChecked(value)
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					obj.Disable()
				} else {
					obj.Enable()
				}
			})

			return obj
		},

		/**
		 * <radio>
		 */
		"radio": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			// TODO: Implement radio group

			var children []string
			for _, child := range node.Nodes {
				if child.GetTag() == "option" {
					children = append(children, child.GetAttr("value"))
				}
			}

			obj := widget.NewRadioGroup(children, func(string) {})

			if node.GetAttrBool("readonly") || node.GetAttrBool("disabled") {
				obj.Disable()
			}

			return obj
		},

		/**
		 * <br>
		 */
		"br": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			return widget.NewLabel("")
		},

		/**
		 * <accordion>
		 */
		"accordion": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var children []*widget.AccordionItem
			for _, child := range node.Nodes {
				title := child.GetAttr("title")
				content := *dom.parseNode(&child)
				children = append(children, widget.NewAccordionItem(title, content))
			}

			return widget.NewAccordion(children...)
		},

		/**
		 * <activity>
		 */
		"activity": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			return widget.NewActivity()
		},

		/**
		 * <card>
		 */
		"card": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)
			wrapper := container.NewVBox(children...)

			obj := widget.NewCard("", "", wrapper)

			bindString("title", node, dom, func(value string) {
				obj.Title = value
				obj.Refresh()
			})

			bindString("subtitle", node, dom, func(value string) {
				obj.Subtitle = value
				obj.Refresh()
			})

			return obj
		},

		/**
		 * <a>
		 */
		"a": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewHyperlink(node.GetText(), nil)

			bindContent(node, dom, func(value string) {
				obj.SetText(value)
			})

			bindString("href", node, dom, func(value string) {
				obj.SetURLFromString(value)
			})

			return obj
		},

		/**
		 * <icon>
		 */
		"icon": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			// TODO: this is probably wrong
			name := node.GetText()
			icon := fyne.NewStaticResource(name, nil)
			return widget.NewIcon(icon)
		},

		/**
		 * <progress>
		 */
		"progress": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewProgressBar()

			bindFloat("value", node, dom, func(value float64) {
				obj.SetValue(value)
			})

			bindFloat("min", node, dom, func(value float64) {
				obj.Min = value
				obj.Refresh()
			})

			bindFloat("max", node, dom, func(value float64) {
				obj.Max = value
				obj.Refresh()
			})

			return obj
		},

		/**
		 * <loader>
		 */
		"loader": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			return widget.NewProgressBarInfinite()
		},

		/**
		 * <markdown>
		 */
		"markdown": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			// TODO: ???
			//rich := widget.NewRichText();
			//rich2 := widget.NewRichTextWithText(node.Content);

			// bind content?

			text := node.GetText()
			return widget.NewRichTextFromMarkdown(text)
		},

		/**
		 * <select>
		 */
		"select": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var options []string
			for _, child := range node.Nodes {
				if child.GetTag() == "option" {
					label := child.GetText()
					options = append(options, label)
				}
			}

			obj := widget.NewSelect(options, func(string) {})

			// TODO: bind options?

			obj.OnChanged = bindString("value", node, dom, func(value string) {
				for _, option := range obj.Options {
					if option == value {
						obj.SetSelected(value)
						break
					}
				}
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					obj.Disable()
				} else {
					obj.Enable()
				}
			})

			return obj
		},

		/**
		 * <combobox>
		 */
		"combobox": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var options []string
			for _, child := range node.Nodes {
				if child.GetTag() == "option" {
					value := child.GetText()
					options = append(options, value)
					continue
				}
			}

			obj := widget.NewSelectEntry(options)

			// TODO: bind options?

			obj.OnChanged = bindString("value", node, dom, func(value string) {
				obj.SetText(value)
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					obj.Disable()
				} else {
					obj.Enable()
				}
			})

			return obj
		},

		/**
		 * <hr>
		 */
		"hr": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			return widget.NewSeparator()
		},

		/**
		 * <slider>
		 */
		"slider": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewSlider(0, 100)

			bindFloat("value", node, dom, func(value float64) {
				obj.Value = value
				obj.Refresh()
			})

			bindFloat("min", node, dom, func(value float64) {
				obj.Min = value
				obj.Refresh()
			})

			bindFloat("max", node, dom, func(value float64) {
				obj.Max = value
				obj.Refresh()
			})

			bindFloat("step", node, dom, func(value float64) {
				obj.Step = value
				obj.Refresh()
			})

			obj.OnChanged = bindFloat("value", node, dom, func(value float64) {
				obj.Value = value
				obj.Refresh()
			})

			bindBool("disabled", node, dom, func(value bool) {
				if value {
					obj.Disable()
				} else {
					obj.Enable()
				}
			})

			return obj
		},

		/**
		 * <code>
		 */
		// TODO: test this
		"code": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			obj := widget.NewTextGridFromString("")

			bindContent(node, dom, func(value string) {
				obj.SetText(value)
			})

			bindBool("line-numbers", node, dom, func(value bool) {
				obj.ShowLineNumbers = value
				obj.Refresh()
			})

			return obj
		},

		/**
		 * <toolbar>
		 */
		"toolbar": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var items = []widget.ToolbarItem{}

			for _, child := range node.Nodes {
				var item widget.ToolbarItem

				tag := child.GetTag()
				switch tag {
				case "toolbar-action":
					icon := child.GetAttr("icon")
					item = widget.NewToolbarAction(fyne.NewStaticResource(icon, nil), func() {})
				case "toolbar-spacer":
					item = widget.NewToolbarSpacer()
				case "toolbar-separator":
					item = widget.NewToolbarSeparator()
				}

				if item != nil {
					// TODO: make this work?
					/*id := child.GetAttr("id")
					if id != "" {
						d.refs[id] = item
					}*/

					items = append(items, item)
				}
			}

			return widget.NewToolbar(items...)
		},

		/**
		 * <list>
		 */
		"list": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var tpl fyne.CanvasObject
			if len(node.Nodes) == 1 {
				tpl = *dom.parseNode(&node.Nodes[0])
			} else {
				children := dom.parseChildren(node)
				tpl = container.NewVBox(children...)
			}

			return widget.NewList(
				func() int { return 0 },
				func() fyne.CanvasObject { return tpl },
				func(widget.ListItemID, fyne.CanvasObject) {},
				// TODO: /\ this is weird. this is supposed to change the item if the data changes?
			)
		},

		/**
		 * <table>
		 */
		"table": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			data := [][]string{}
			for _, child := range node.Nodes {
				if child.GetTag() == "tr" {
					row := []string{}
					for _, cell := range child.Nodes {
						if cell.GetTag() == "td" {
							row = append(row, cell.Content)
						}
					}
					data = append(data, row)
				}
			}

			table := widget.NewTable(
				func() (int, int) {
					return len(data), len(data[0])
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("")
				},
				func(i widget.TableCellID, obj fyne.CanvasObject) {
					label := obj.(*widget.Label)
					label.SetText(data[i.Row][i.Col]) // Set data dynamically
				},
			)

			// TODO: Adjust column widths
			table.SetColumnWidth(0, 50)  // ID column
			table.SetColumnWidth(1, 100) // Name column
			table.SetColumnWidth(2, 50)  // Age column

			return table
		},

		/**
		 * <tree>
		 */
		"tree": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			// TODO: circle back to this
			/*
				data := map[string][]string{
					"":           {"Fruits", "Vegetables"}, // Root nodes
					"Fruits":     {"Apple", "Banana", "Cherry"},
					"Vegetables": {"Carrot", "Broccoli", "Potato"},
				}
				tree := widget.NewTree(
					func(id string) []string {
						// Return child nodes of the given parent
						return data[id]
					},
					func(id string) bool {
						// Return true if the node has children (i.e., is expandable)
						_, hasChildren := data[id]
						return hasChildren
					},
					func(id string) fyne.CanvasObject {
						// Template for each row
						return widget.NewLabel("Node")
					},
					func(id string, obj fyne.CanvasObject) {
						// Set the text dynamically
						label := obj.(*widget.Label)
						label.SetText(id)
					},
				)
			*/
			return nil
		},

		/*
			// TODO: circle back to this
			/**
			* <grid-wrap>
		*/ /*
			"grid-wrap": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
				children := dom.parseChildren(node)

				widthStr := node.GetAttr("width")
				heightStr := node.GetAttr("height")
				width, _ := strconv.ParseFloat(widthStr, 32)
				height, _ := strconv.ParseFloat(heightStr, 32)
				size := fyne.NewSize(float32(width), float32(height))

				return widget.NewGridWrap(size, children...)
			},
		*/

		/**
		 * <tabs>
		 */
		"tabs": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var items []*container.TabItem

			for _, child := range node.Nodes {
				if child.GetTag() == "tab" {
					title := child.GetAttr("title")
					content := *dom.parseNode(&child)

					// TODO: disabled tab

					tab := container.NewTabItem(title, content)
					items = append(items, tab)
				}
			}

			tabs := container.NewAppTabs(items...)

			location := node.GetAttr("location")
			switch location {
			case "top":
				tabs.SetTabLocation(container.TabLocationTop)
			case "bottom":
				tabs.SetTabLocation(container.TabLocationBottom)
			case "left":
				tabs.SetTabLocation(container.TabLocationLeading)
			case "right":
				tabs.SetTabLocation(container.TabLocationTrailing)
			}

			return tabs
		},

		/**
		 * <scroller>
		 */
		"scroller": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := dom.parseChildren(node)

			if node.GetAttr("dir") == "horizontal" {
				return container.NewHScroll(container.NewHBox(children...))
			} else {
				return container.NewVScroll(container.NewVBox(children...))
			}
		},

		/**
		 * <spacer>
		 */
		"spacer": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			return layout.NewSpacer()
		},

		/**
		 * <split>
		 */
		"split": func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			var left fyne.CanvasObject
			var right fyne.CanvasObject

			if len(node.Nodes) > 0 {
				left = *dom.parseNode(&node.Nodes[0])
			}
			if len(node.Nodes) > 1 {
				right = *dom.parseNode(&node.Nodes[1])
			}

			var split *container.Split
			direction := node.GetAttr("dir")
			if direction == "vertical" {
				split = container.NewVSplit(left, right)
			} else {
				split = container.NewHSplit(left, right)
			}

			leftSize := node.Nodes[0].GetAttrFloat("size")
			rightSize := node.Nodes[1].GetAttrFloat("size")
			if leftSize > 0 && rightSize > 0 {
				split.SetOffset(leftSize / (leftSize + rightSize))
			}

			return split
		},
	}
}

func (dom *DOM) parseNode(node *XMLNode) *fyne.CanvasObject {
	if node == nil {
		return nil
	}

	var obj fyne.CanvasObject

	tag := node.GetTag()
	if handler, ok := renderers[tag]; ok {
		obj = handler(node, dom)
	} else {
		obj = widget.NewLabel("<unsupported tag: " + tag + ">")
	}

	id := node.GetAttr("id")
	if id != "" {
		dom.refs[id] = obj
	}

	bindBool("hidden", node, dom, func(value bool) {
		if value {
			obj.Hide()
		} else {
			obj.Show()
		}
	})

	return &obj
}

func (dom *DOM) parseChildren(node *XMLNode) []fyne.CanvasObject {
	var children []fyne.CanvasObject
	for _, child := range node.Nodes {
		children = append(children, *dom.parseNode(&child))
	}
	return children
}

func parseColor(str string) (color.Color, error) {
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

func bindString(name string, node *XMLNode, dom *DOM, update func(string)) func(string) {
	value := node.GetAttr(name)
	if value != "" {
		update(value)
	}

	bind := node.GetBind(name)
	if bind != "" {
		state := dom.UseState().String(node.GetBind(name), value)
		state.AddListener(update)
		return state.Set
	}

	return nil
}

func bindInt(name string, node *XMLNode, dom *DOM, update func(int)) func(int) {
	value := node.GetAttrInt(name)
	if value != 0 {
		update(value)
	}

	bind := node.GetBind(name)
	if bind != "" {
		state := dom.UseState().Int(node.GetBind(name), value)
		state.AddListener(update)
		return state.Set
	}

	return nil
}

func bindFloat(name string, node *XMLNode, dom *DOM, update func(float64)) func(float64) {
	value := node.GetAttrFloat(name)
	if value != 0 {
		update(value)
	}

	bind := node.GetBind(name)
	if bind != "" {
		state := dom.UseState().Float(node.GetBind(name), value)
		state.AddListener(update)
		return state.Set
	}

	return nil
}

func bindBool(name string, node *XMLNode, dom *DOM, update func(bool)) func(bool) {
	value := node.GetAttrBool(name)
	if value {
		update(value)
	}

	bind := node.GetBind(name)
	if bind != "" {
		state := dom.UseState().Bool(node.GetBind(name), value)
		state.AddListener(update)
		return state.Set
	}

	return nil
}

func bindContent(node *XMLNode, dom *DOM, update func(string)) func(string) {
	value := node.GetText()
	if value != "" {
		update(value)
	}

	bind := node.GetBind("content")
	if bind != "" {
		state := dom.UseState().String(node.GetBind("content"), value)
		state.AddListener(update)
		return state.Set
	}

	return nil
}

func bindCallback(name string, node *XMLNode, dom *DOM) func() {
	if node.HasBind(name) {
		bind := node.GetBind(name)
		return func() {
			if callback, ok := dom.callbacks[bind]; ok {
				callback()
			}
		}
	}

	return nil
}
