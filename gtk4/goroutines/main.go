package main

import (
	"fmt"
	"os"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.goroutines", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	topLabel := gtk.NewLabel("Text set by initializer")
	topLabel.SetVExpand(true)
	topLabel.SetHExpand(true)

	bottomLabel := gtk.NewLabel("Text set by initializer")
	bottomLabel.SetVExpand(true)
	bottomLabel.SetHExpand(true)

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(topLabel)
	box.Append(bottomLabel)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("gotk4 Example")
	window.SetChild(box)
	window.SetDefaultSize(400, 300)
	window.Show()

	go func() {
		var ix int
		for t := range time.Tick(time.Second) {
			// Make a copy of the state so we can reference it in the closure.
			currentTime := t
			currentIx := ix

			ix++

			glib.IdleAdd(func() {
				topLabel.SetLabel(fmt.Sprintf("Set a label %d time(s)!", currentIx))
				bottomLabel.SetLabel(fmt.Sprintf(
					"Last updated at %s.",
					currentTime.Format(time.StampMilli),
				))
			})
		}
	}()
}
