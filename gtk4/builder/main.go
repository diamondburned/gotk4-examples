package main

import (
	"os"
	"strconv"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

//go:embed main.ui
var uiXML string

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.builder", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	// You can build UIs using Cambalache (https://flathub.org/apps/details/ar.xjuan.Cambalache)
	builder := gtk.NewBuilderFromString(uiXML)

	// MainWindow and Button are object IDs from the UI file
	window := builder.GetObject("MainWindow").Cast().(*gtk.Window)
	button := builder.GetObject("Button").Cast().(*gtk.Button)

	counter := 0
	button.Connect("clicked", func() {
		button.SetLabel(strconv.Itoa(counter))
		counter++
	})

	app.AddWindow(window)
	window.Show()
}
