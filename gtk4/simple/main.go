package main

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	gtk.Init()
}

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", 0)
	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("gotk4 Example")
	window.SetChild(gtk.NewLabel("Hello from Go!"))
	window.SetDefaultSize(400, 300)
	window.Show()
}
