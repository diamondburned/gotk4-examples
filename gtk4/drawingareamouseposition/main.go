package main

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const appID = "com.github.diamondburned.gotk4-examples.gtk4.drawingareamouseposition"

func main() {
	app := gtk.NewApplication(appID, 0)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	labelX := gtk.NewLabel("X")
	labelX.SetXAlign(1)
	labelY := gtk.NewLabel("Y")
	labelY.SetXAlign(0)

	motionCtrl := gtk.NewEventControllerMotion()
	motionCtrl.ConnectMotion(func(x, y float64) {
		labelX.SetLabel(fmt.Sprintf("X: %f", x))
		labelY.SetLabel(fmt.Sprintf("Y: %f", y))
	})

	drawArea := gtk.NewDrawingArea()
	drawArea.SetVExpand(true)
	drawArea.AddController(motionCtrl)

	grid := gtk.NewGrid()
	grid.SetColumnHomogeneous(true)
	grid.SetColumnSpacing(5)
	grid.SetRowSpacing(5)
	grid.Attach(labelX, 0, 0, 1, 1)
	grid.Attach(labelY, 1, 0, 1, 1)
	grid.Attach(drawArea, 0, 1, 2, 1)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("drawingareamouseposition - gotk4 Example")
	window.SetChild(grid)
	window.SetDefaultSize(640, 480)
	window.Show()
}
