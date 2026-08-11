package main

import (
	"os"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:generate blueprint-compiler compile ./main.blp --output ./main.ui
//go:embed main.ui
var uiXML string

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.blueprint", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	//You can build UIs using Blueprint (https://jwestman.pages.gitlab.gnome.org/blueprint-compiler/)
	builder := gtk.NewBuilderFromString(uiXML, len(uiXML))

	//MyMainWindow can be found inside the .blp file
	window := builder.GetObject("MyMainWindow").Cast().(*gtk.Window)

	app.AddWindow(window)
	window.Show()
}
