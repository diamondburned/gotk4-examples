package main

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const appID = "com.github.diamondburned.gotk4-examples.gtk4.mouse"

func main() {
	app := gtk.NewApplication(appID, gio.ApplicationFlags(gio.ApplicationFlagsNone))
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

var times = [...]string{
	"once",
	"twice",
	"thrice",
}

func activate(app *gtk.Application) {
	gesture := gtk.NewGestureClick()
	// Listen to all mouse buttons instead of just the left-click.
	gesture.SetButton(0)
	// Only handle pointer and pointer-like touch events. This isn't a good
	// idea for touch-friendly applications.
	gesture.SetExclusive(true)

	label := gtk.NewLabel("Use the mouse buttons.")
	label.GrabFocus()

	var click struct {
		NPress int
		X, Y   float64
	}

	draw := gtk.NewDrawingArea()
	draw.SetDrawFunc(func(draw *gtk.DrawingArea, cr *cairo.Context, w, h int) {
		// Draw a red rectagle at the X and Y location.
		cr.SetSourceRGB(255, 0, 0)
		cr.Rectangle(click.X, click.Y, 10, 10)
		cr.Fill()
	})

	// Pressed is a signal that's emitted everytime the button is pressed down
	// (but not released - it's not a full click).
	gesture.ConnectPressed(func(n int, x, y float64) {
		click.NPress = n
		click.X = x
		click.Y = y
		// Queue the rectangle to be drawn.
		draw.QueueDraw()

		// Format the message for the label and update it.
		msg := fmt.Sprintf("Button %d pressed", gesture.CurrentButton())
		if click.NPress > 0 && click.NPress < 4 {
			msg += " " + times[click.NPress-1] + "."
		} else {
			msg += fmt.Sprintf(" %d times.", click.NPress)
		}

		label.SetLabel(msg)
	})

	overlay := gtk.NewOverlay()
	overlay.SetVExpand(true)
	overlay.SetHExpand(true)
	overlay.AddController(gesture)
	overlay.SetChild(label)
	overlay.AddOverlay(draw)

	window := gtk.NewApplicationWindow(app)
	window.SetDefaultSize(200, 200)
	window.SetTitle("mouse - gotk4 Example")
	window.SetChild(overlay)
	window.Show()
}
