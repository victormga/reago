package reago

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func init() {
	/** <row> */
	Parser.RegisterTag("row", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.NewHBox(children...)

		top, bottom, left, right := node.GetPadding()
		if (top > 0) || (bottom > 0) || (left > 0) || (right > 0) {
			obj = container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				obj,
			)
		}

		color, _ := Parser.ParseColor(node.GetAttr("background-color"))
		if color != nil {
			bg := canvas.NewRectangle(color)
			obj = container.NewStack(bg, obj)
		}

		return obj
	})

	/** <col> */
	Parser.RegisterTag("col", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.NewVBox(children...)

		top, bottom, left, right := node.GetPadding()
		if (top > 0) || (bottom > 0) || (left > 0) || (right > 0) {
			obj = container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				obj,
			)
		}

		color, _ := Parser.ParseColor(node.GetAttr("background-color"))
		if color != nil {
			bg := canvas.NewRectangle(color)
			obj = container.NewStack(bg, obj)
		}

		return obj
	})

	/** <flex-row> */
	Parser.RegisterTag("flex-row", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.New(layout.NewGridLayoutWithColumns(len(children)), children...)

		top, bottom, left, right := node.GetPadding()
		if (top > 0) || (bottom > 0) || (left > 0) || (right > 0) {
			obj = container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				obj,
			)
		}

		color, _ := Parser.ParseColor(node.GetAttr("background-color"))
		if color != nil {
			bg := canvas.NewRectangle(color)
			obj = container.NewStack(bg, obj)
		}

		return obj
	})

	/** <flex-col> */
	Parser.RegisterTag("flex-col", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.New(layout.NewGridLayoutWithRows(len(children)), children...)

		top, bottom, left, right := node.GetPadding()
		if (top > 0) || (bottom > 0) || (left > 0) || (right > 0) {
			obj = container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				obj,
			)
		}

		color, _ := Parser.ParseColor(node.GetAttr("background-color"))
		if color != nil {
			bg := canvas.NewRectangle(color)
			obj = container.NewStack(bg, obj)
		}

		return obj
	})

	Parser.RegisterTag("stack", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.NewStack(children...)
		return obj
	})

	/** <center> */
	Parser.RegisterTag("center", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		obj := container.NewCenter(children...)

		top, bottom, left, right := node.GetPadding()
		if (top > 0) || (bottom > 0) || (left > 0) || (right > 0) {
			obj = container.New(
				layout.NewCustomPaddedLayout(top, bottom, left, right),
				obj,
			)
		}

		return obj
	})

	/** <form> */
	Parser.RegisterTag("form", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		// TODO: circle back to this
		//children := Parser.ParseChildren(node, dom)
		//return widget.NewForm(children...)
		return nil
	})

	/** <grid> */
	Parser.RegisterTag("grid", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)

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
	})

	/** <layout> */
	Parser.RegisterTag("layout", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var top fyne.CanvasObject
		var bottom fyne.CanvasObject
		var left fyne.CanvasObject
		var right fyne.CanvasObject
		var content []fyne.CanvasObject

		for _, child := range node.Nodes {
			switch child.GetTag() {
			case "top":
				children := Parser.ParseChildren(&child, dom)
				top = container.NewHBox(children...)
			case "bottom":
				children := Parser.ParseChildren(&child, dom)
				bottom = container.NewHBox(children...)
			case "left":
				children := Parser.ParseChildren(&child, dom)
				left = container.NewHBox(children...)
			case "right":
				children := Parser.ParseChildren(&child, dom)
				right = container.NewHBox(children...)
			default:
				content = append(content, Parser.ParseNode(&child, dom))
			}
		}

		return container.NewBorder(top, bottom, left, right, content...)
	})

	/** <label> */
	Parser.RegisterTag("label", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewLabel("")

		node.BindContent(dom, func(value string) {
			obj.SetText(value)
		})

		return obj
	})

	/** <text> */
	Parser.RegisterTag("text", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := canvas.NewText("", theme.Color(theme.ColorNameForeground))

		node.BindContent(dom, func(value string) {
			obj.Text = value
			obj.Refresh()
		})

		node.BindString("color", dom, func(value string) {
			color, _ := Parser.ParseColor(value)
			if color != nil {
				obj.Color = color
				obj.Refresh()
			}
		})

		node.BindFloat("size", dom, func(value float64) {
			obj.TextSize = float32(value)
			obj.Refresh()
		})

		node.BindString("style", dom, func(value string) {
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

		node.BindString("align", dom, func(value string) {
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

		node.BindContent(dom, func(value string) {
			obj.Text = value
			obj.Refresh()
		})

		return obj
	})

	/** <button> */
	Parser.RegisterTag("button", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewButtonWithIcon("", nil, func() {})

		obj.OnTapped = node.BindCallback("click", dom)

		node.BindString("icon", dom, func(value string) {
			obj.Icon = theme.Icon(fyne.ThemeIconName(value))
			obj.Refresh()
		})

		node.BindContent(dom, func(value string) {
			obj.SetText(value)
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				obj.Disable()
			} else {
				obj.Enable()
			}
		})

		return obj
	})

	/** <input> */
	Parser.RegisterTag("input", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
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

		node.BindString("placeholder", dom, func(value string) {
			entry.SetPlaceHolder(value)
		})

		entry.OnChanged = node.BindString("value", dom, func(value string) {
			entry.SetText(value)
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				entry.Disable()
			} else {
				entry.Enable()
			}
		})

		//entry.OnSubmitted = node.BindCallback("submit", dom)

		if node.HasAttr("validation") {
			msg := node.GetAttr("validation-message")
			if msg == "" {
				msg = "invalid"
			}
			entry.Validator = validation.NewRegexp(node.GetAttr("validation"), msg)
		}

		return entry
	})

	/** <textarea> */
	Parser.RegisterTag("textarea", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		entry := widget.NewMultiLineEntry()

		node.BindString("placeholder", dom, func(value string) {
			entry.SetPlaceHolder(value)
		})

		entry.OnChanged = node.BindString("value", dom, func(value string) {
			entry.SetText(value)
		})

		node.BindBool("disabled", dom, func(value bool) {
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
	})

	/** <checkbox> */
	Parser.RegisterTag("checkbox", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewCheck(node.GetAttr("label"), func(bool) {})

		node.BindString("label", dom, func(value string) {
			obj.Text = value
			obj.Refresh()
		})

		obj.OnChanged = node.BindBool("value", dom, func(value bool) {
			obj.SetChecked(value)
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				obj.Disable()
			} else {
				obj.Enable()
			}
		})

		return obj
	})

	/** <radio> */
	Parser.RegisterTag("radio", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
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
	})

	/** <br> */
	Parser.RegisterTag("br", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		return widget.NewLabel("")
	})

	/** <accordion> */
	Parser.RegisterTag("accordion", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var children []*widget.AccordionItem
		for _, child := range node.Nodes {
			title := child.GetAttr("title")
			content := Parser.ParseNode(&child, dom)
			children = append(children, widget.NewAccordionItem(title, content))
		}

		return widget.NewAccordion(children...)
	})

	/** <activity> */
	Parser.RegisterTag("activity", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		return widget.NewActivity()
	})

	/** <card> */
	Parser.RegisterTag("card", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)
		wrapper := container.NewHBox(children...)

		obj := widget.NewCard("", "", wrapper)

		node.BindString("title", dom, func(value string) {
			obj.Title = value
			obj.Refresh()
		})

		node.BindString("subtitle", dom, func(value string) {
			obj.Subtitle = value
			obj.Refresh()
		})

		return obj
	})

	/** <a> */
	Parser.RegisterTag("a", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewHyperlink(node.GetContent(), nil)

		node.BindContent(dom, func(value string) {
			obj.SetText(value)
		})

		node.BindString("href", dom, func(value string) {
			obj.SetURLFromString(value)
		})

		return obj
	})

	/** <icon> */
	Parser.RegisterTag("icon", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		// TODO: this is probably wrong
		name := node.GetContent()
		icon := fyne.NewStaticResource(name, nil)
		return widget.NewIcon(icon)
	})

	/** <progress> */
	Parser.RegisterTag("progress", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewProgressBar()

		node.BindFloat("value", dom, func(value float64) {
			obj.SetValue(value)
		})

		node.BindFloat("min", dom, func(value float64) {
			obj.Min = value
			obj.Refresh()
		})

		node.BindFloat("max", dom, func(value float64) {
			obj.Max = value
			obj.Refresh()
		})

		return obj
	})

	/** <loader> */
	Parser.RegisterTag("loader", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		return widget.NewProgressBarInfinite()
	})

	/** <markdown> */
	Parser.RegisterTag("markdown", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		// TODO: ???
		//rich := widget.NewRichText();
		//rich2 := widget.NewRichTextWithText(node.Content);

		// bind content?

		text := node.GetContent()
		return widget.NewRichTextFromMarkdown(text)
	})

	/** <select> */
	Parser.RegisterTag("select", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewSelect([]string{}, func(string) {})

		if node.HasBind("options") {
			node.BindList("options", dom, func(value []any) {
				var options []string
				for _, v := range value {
					options = append(options, v.(string))
				}
				obj.Options = options
				obj.Refresh()
			})
		} else {
			var options []string
			for _, child := range node.Nodes {
				if child.GetTag() == "option" {
					label := child.GetContent()
					options = append(options, label)
				}
			}
			obj.Options = options
		}

		obj.OnChanged = node.BindString("value", dom, func(value string) {
			for _, option := range obj.Options {
				if option == value {
					obj.SetSelected(value)
					break
				}
			}
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				obj.Disable()
			} else {
				obj.Enable()
			}
		})

		return obj
	})

	/** <combobox> */
	Parser.RegisterTag("combobox", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewSelectEntry([]string{})

		if node.HasBind("options") {
			node.BindList("options", dom, func(value []any) {
				var options []string
				for _, v := range value {
					options = append(options, v.(string))
				}
				obj.SetOptions(options)
			})
		} else {
			var options []string
			for _, child := range node.Nodes {
				if child.GetTag() == "option" {
					label := child.GetContent()
					options = append(options, label)
				}
			}
			obj.SetOptions(options)
		}

		obj.OnChanged = node.BindString("value", dom, func(value string) {
			obj.SetText(value)
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				obj.Disable()
			} else {
				obj.Enable()
			}
		})

		return obj
	})

	/** <hr> */
	Parser.RegisterTag("hr", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		return widget.NewSeparator()
	})

	/** <slider> */
	Parser.RegisterTag("slider", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := widget.NewSlider(0, 100)

		node.BindFloat("value", dom, func(value float64) {
			obj.Value = value
			obj.Refresh()
		})

		node.BindFloat("min", dom, func(value float64) {
			obj.Min = value
			obj.Refresh()
		})

		node.BindFloat("max", dom, func(value float64) {
			obj.Max = value
			obj.Refresh()
		})

		node.BindFloat("step", dom, func(value float64) {
			obj.Step = value
			obj.Refresh()
		})

		obj.OnChanged = node.BindFloat("value", dom, func(value float64) {
			obj.Value = value
			obj.Refresh()
		})

		node.BindBool("disabled", dom, func(value bool) {
			if value {
				obj.Disable()
			} else {
				obj.Enable()
			}
		})

		return obj
	})

	/** <code> */
	Parser.RegisterTag("code", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		// TODO: test this
		obj := widget.NewTextGridFromString("")

		node.BindContent(dom, func(value string) {
			obj.SetText(value)
		})

		node.BindBool("line-numbers", dom, func(value bool) {
			obj.ShowLineNumbers = value
			obj.Refresh()
		})

		return obj
	})

	/** <toolbar> */
	Parser.RegisterTag("toolbar", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var items = []widget.ToolbarItem{}

		for _, child := range node.Nodes {
			var item widget.ToolbarItem

			tag := child.GetTag()
			switch tag {
			case "action":
				icon := child.GetAttr("icon")
				item = widget.NewToolbarAction(fyne.NewStaticResource(icon, nil), func() {})
			case "spacer":
				item = widget.NewToolbarSpacer()
			case "separator":
				item = widget.NewToolbarSeparator()
			}

			if item != nil {
				items = append(items, item)
			}
		}

		return widget.NewToolbar(items...)
	})

	/** <list> */
	Parser.RegisterTag("list", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		if !node.HasBind("items") {
			return widget.NewLabel("<missing bind property in list>")
		}

		bind := node.GetBind("items")

		obj := widget.NewList(
			func() int {
				list := dom.UseState().GetList(bind)
				return list.container.Length()
			},
			func() fyne.CanvasObject {
				fragment := dom.Clone()
				children := Parser.ParseChildren(node, fragment)
				return container.NewHBox(children...)
			},
			func(idx widget.ListItemID, obj fyne.CanvasObject) {
				fragment := dom.Clone()

				item := dom.UseState().GetList(bind).GetValue(idx)
				for key, value := range ParseStruct(item) {
					fragment.UseState().String(key, value)
				}

				wrapper := obj.(*fyne.Container)
				wrapper.Objects = Parser.ParseChildren(node, fragment)
				wrapper.Refresh()
			},
		)

		dom.UseState().GetList(bind).OnChange(func(_ []any) {
			obj.Refresh()
		})

		return obj
	})

	/** <table> */
	Parser.RegisterTag("table", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
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
	})

	/** <tree> */
	Parser.RegisterTag("tree", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		// TODO: circle back to this
		/*
			data := map[string][]string{
				Parser.Register("":           {"Fruits", "Vegetables"}, // Root nodes
				Parser.Register("Fruits":     {"Apple", "Banana", "Cherry"},
				Parser.Register("Vegetables": {"Carrot", "Broccoli", "Potato"},
			}
			tree := widget.NewTree(
				func(id string) []string {
					// Return child nodes of the given parent
					return data[id]
				})
				func(id string) bool {
					// Return true if the node has children (i.e., is expandable)
					_, hasChildren := data[id]
					return hasChildren
				})
				func(id string) fyne.CanvasObject {
					// Template for each row
					return widget.NewLabel("Node")
				})
				func(id string, obj fyne.CanvasObject) {
					// Set the text dynamically
					label := obj.(*widget.Label)
					label.SetText(id)
				})
			)
		*/
		return nil
	})

	/*
		// TODO: circle back to this
		/**
		* <grid-wrap>
	*/ /*
		Parser.Register("grid-wrap", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
			children := Parser.ParseChildren(node, dom)

			widthStr := node.GetAttr("width")
			heightStr := node.GetAttr("height")
			width, _ := strconv.ParseFloat(widthStr, 32)
			height, _ := strconv.ParseFloat(heightStr, 32)
			size := fyne.NewSize(float32(width), float32(height))

			return widget.NewGridWrap(size, children...)
		})
	*/

	/** <tabs> */
	Parser.RegisterTag("tabs", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var items []*container.TabItem

		for _, child := range node.Nodes {
			if child.GetTag() == "tab" {
				title := child.GetAttr("title")
				content := Parser.ParseNode(&child, dom)

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
	})

	/** <scroll> */
	Parser.RegisterTag("scroll", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		children := Parser.ParseChildren(node, dom)

		if node.GetAttr("dir") == "horizontal" {
			return container.NewHScroll(container.NewHBox(children...))
		} else {
			return container.NewVScroll(container.NewHBox(children...))
		}
	})

	/** <spacer> */
	Parser.RegisterTag("spacer", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		return layout.NewSpacer()
	})

	/** <split> */
	Parser.RegisterTag("split", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var left fyne.CanvasObject
		var right fyne.CanvasObject

		if len(node.Nodes) > 0 {
			left = Parser.ParseNode(&node.Nodes[0], dom)
		}
		if len(node.Nodes) > 1 {
			right = Parser.ParseNode(&node.Nodes[1], dom)
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
	})

	Parser.RegisterTag("circle", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := canvas.NewCircle(theme.Color(theme.ColorNameForeground))

		node.BindString("background-color", dom, func(value string) {
			color, _ := Parser.ParseColor(value)
			if color != nil {
				obj.FillColor = color
			}
		})

		node.BindString("border-color", dom, func(value string) {
			color, _ := Parser.ParseColor(value)
			if color != nil {
				obj.StrokeColor = color
			}
		})

		node.BindFloat("border-size", dom, func(value float64) {
			obj.StrokeWidth = float32(value)
		})

		node.BindFloat("size", dom, func(value float64) {
			obj.Resize(fyne.NewSize(float32(value), float32(value)))
		})

		return obj
	})

	Parser.RegisterTag("img", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		obj := canvas.NewImageFromResource(nil)

		node.BindString("src", dom, func(value string) {
			obj.Resource = fyne.NewStaticResource(value, nil)
			obj.Refresh()
		})

		node.BindFloat("width", dom, func(value float64) {
			obj.Resize(fyne.NewSize(float32(value), obj.Size().Height))
		})

		node.BindFloat("height", dom, func(value float64) {
			obj.Resize(fyne.NewSize(obj.Size().Width, float32(value)))
		})

		node.BindString("fill", dom, func(value string) {
			switch value {
			case "contain":
				obj.FillMode = canvas.ImageFillContain
			case "stretch":
				obj.FillMode = canvas.ImageFillStretch
			default:
				obj.FillMode = canvas.ImageFillOriginal
			}
			obj.Refresh()
		})

		return obj
	})

	Parser.RegisterTag("gradient", func(node *XMLNode, dom *DOM) fyne.CanvasObject {
		var obj *canvas.LinearGradient

		dir := node.GetAttr("direction")
		switch dir {
		case "vertical":
			obj = canvas.NewHorizontalGradient(theme.Color(theme.ColorNameForeground), theme.Color(theme.ColorNameBackground))
		default:
			obj = canvas.NewVerticalGradient(theme.Color(theme.ColorNameForeground), theme.Color(theme.ColorNameBackground))
		}

		node.BindString("start", dom, func(value string) {
			color, _ := Parser.ParseColor(value)
			if color != nil {
				obj.StartColor = color
			}
		})

		node.BindString("end", dom, func(value string) {
			color, _ := Parser.ParseColor(value)
			if color != nil {
				obj.EndColor = color
			}
		})

		node.BindFloat("angle", dom, func(value float64) {
			obj.Angle = value
		})

		return obj
	})
}
