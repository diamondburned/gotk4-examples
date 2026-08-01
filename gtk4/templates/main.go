package main

import (
	_ "embed"
	"os"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed app.ui
var Embedded_AppUI []byte

//go:embed style.css
var Embedded_CSS string

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.templates", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() {
		window := NewCustomAppWindow()
		app.AddWindow(&window.Window)

		display := window.Widget.Display()
		css := gtk.NewCSSProvider()
		css.LoadFromData(Embedded_CSS)
		gtk.StyleContextAddProviderForDisplay(display, css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

		window.SetVisible(true)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
