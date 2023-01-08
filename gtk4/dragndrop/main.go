package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	coreglib "github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Simple window with a label where you can drop things
func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.dragndrop", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Drag & Drop Example")
	lbl := gtk.NewLabel("Drop something here!")
	window.SetChild(lbl)
	window.SetDefaultSize(400, 300)

	drop := gtk.NewDropTarget(glib.TypeString, gdk.ActionCopy)
	drop.Connect("drop", func(drop *gtk.DropTarget, src *coreglib.Value, x, y float64) {
		str, ok := src.GoValue().(string)
		if !ok {
			lbl.SetText("invalid drag source")
			return
		}

		u, err := url.Parse(strings.TrimSpace(str))
		if err != nil {
			lbl.SetText(fmt.Sprintf("invalid URL %s: %s", str, err))
			return
		}

		switch u.Scheme {
		case "file":
			lbl.SetText(fmt.Sprintf("File dropped:\n%s", u.Path))
		case "http", "https":
			lbl.SetMarkup(fmt.Sprintf("URL dropped:\n<a href=\"%s\">%s</a>", str, str))
		default:
			lbl.SetText(fmt.Sprintf("Something else dropped:\n%s", str))
		}

	})

	lbl.AddController(drop)
	window.Show()
}
