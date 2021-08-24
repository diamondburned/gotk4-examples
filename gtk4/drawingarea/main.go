package main

import (
	"fmt"
	"log"
	"os"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/cairo"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

const appID = "com.github.diamondburned.gotk4-examples.gtk4.drawingarea"

//go:embed trollface.png
var trollfacePNG []byte

func main() {
	app := gtk.NewApplication(appID, 0)
	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

// State describes the cursor state.
type State struct {
	X float64
	Y float64
}

func activate(app *gtk.Application) {
	trollface, err := loadPNG(trollfacePNG)
	if err != nil {
		log.Fatalln("failed to load trollface.png:", err)
	}

	const rectangleSize = 20
	var state State

	drawArea := gtk.NewDrawingArea()
	drawArea.SetVExpand(true)
	drawArea.SetDrawFunc(func(area *gtk.DrawingArea, cr *cairo.Context, w, h int) {
		gdk.CairoSetSourcePixbuf(cr, trollface, state.X, state.Y)
		cr.Paint()
	})

	motionCtrl := gtk.NewEventControllerMotion()
	motionCtrl.Connect("motion", func(_ *gtk.EventControllerMotion, x, y float64) {
		state.X = x
		state.Y = y
		drawArea.QueueDraw()
	})
	drawArea.AddController(motionCtrl)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("drawingarea - gotk4 Example")
	window.SetChild(drawArea)
	window.SetDefaultSize(640, 480)
	window.Show()
}

func loadPNG(data []byte) (*gdkpixbuf.Pixbuf, error) {
	l, err := gdkpixbuf.NewPixbufLoaderWithType("png")
	if err != nil {
		return nil, fmt.Errorf("NewLoaderWithType png: %w", err)
	}
	defer l.Close()

	if err := l.Write(data); err != nil {
		return nil, fmt.Errorf("PixbufLoader.Write: %w", err)
	}

	if err := l.Close(); err != nil {
		return nil, fmt.Errorf("PixbufLoader.Close: %w", err)
	}

	return l.Pixbuf(), nil
}
