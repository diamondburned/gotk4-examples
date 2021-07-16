package main

import (
	"log"
	"os"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const appID = "com.github.diamondburned.gotk4-examples.gtk4.mouse"

func main() {
	app := gtk.NewApplication(appID, 0)
	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	evkey := gtk.NewEventControllerKey()
	evkey.Connect("key-pressed",
		func(evkey *gtk.EventControllerKey, val, code uint, mod gdk.ModifierType) {
			log.Printf("val = %d, code = %d, mod = %v", val, code, mod)
		},
	)

	label := gtk.NewLabel("Use the mouse buttons.")
	label.AddController(evkey)
	label.SetSizeRequest(200, 200)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("mouse - gotk4 Example")
	window.SetChild(label)
	window.SetDefaultSize(640, 480)
	window.Show()
}
