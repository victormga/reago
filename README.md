## ReaGO
ReaGO is a Go library for building truly native, reactive GUIs. Unlike other solutions, it does not rely on WebViews or workarounds—everything is rendered natively.
This is currently under development and a lot of things might change. It's not production ready.

#
---
#### XML-Based Syntax
ReaGO uses an XML-based syntax to define UI templates, making it intuitive and familiar for developers experienced with HTML. Many standard HTML-like tags are supported, and in most cases, they behave similarly.
``` go
func main() {
	dom := reago.NewDOM()

	dom.Template(`
		<padding all="15">
			<row>
				<col>
					<label>This will be displayed on the left.</label>
					<label>This will be displayed below /\.</label>
				</col>
				<col>
					<label>This will be displayed on the right.</label>
					<label>This will be displayed below /\.</label>
				</col>
			</row>
		</padding>
	`)

	window := reago.NewWindow("My App", 720, 480)
	window.Show(dom)
}
```

#
---
#### Reactive Data Binding
You can sync the data between GO and the UI by using `UseState`.
``` go
func main() {
	dom := reago.NewDOM()

	clock := dom.UseState().String("clock", "...")
	go func() {
		for {
			clock.Set(time.Now().Format(time.RFC1123Z))
			time.Sleep(time.Second)
		}
	}()

	dom.Template(`
		<center>
			<text size="20" bind:content="clock"></text>
		</center>
	`)

	window := reago.NewWindow("My App", 400, 600)
	window.Show(dom)
}
```

#
---
#### Two-Way Binding
Two-way binding is supported. You can bind data between GO and the UI by using `UseState`.
``` go
func main() {
	dom := reago.NewDOM()

	username := dom.UseState().String("username", "")
	username.OnChange(func(value string) {
		println("The username was changed to:", value)
	})

	dom.Template(`
		<center>
			<label>Username</label>
			<input type="text" bind:value="username" />
		</center>
	`)

	window := reago.NewWindow("My App", 400, 600)
	window.Show(dom)
}
```

#
---
#### Callbacks
You can also assign callbacks for events.
``` go
func main() {
	dom := reago.NewDOM()
	
	dom.UseCallback("clicked", func() {
		println("Button clicked!")
	})
	
	dom.Template(`
		<padding all="15">
			<col>
				<label>Click the button below</label>
				<button bind:click="clicked">Click me!</button>
			</col>
		</padding>
	`)

	window := reago.NewWindow("My App", 400, 600)
	window.Show(dom)
}
```

#
#
#
---
#### Loading From Files
You can also load the XML string from a file instead of passing the string.
```

```

#
#
#
---
#### Watching For Changes
This will rerender the window every time the xml file changes (works on dev).
```

```

#
#
#
---
#### Requirements
- GO 1.24 (Older versions might be supported, but no test was done).
- ReaGO uses `fyne v2` under the hood for renderization, so the `fyne v2` environment prerequisites are required for it to work. You can prepare your environment by following the instructions in `https://docs.fyne.io/started/`.

#
#
#
---
## LICENSE
Copyright [2025] Victor Sabiar

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.