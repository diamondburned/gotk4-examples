package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", 0)

	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	window := gtk.NewApplicationWindow(app)
	grid := gtk.NewGrid()

	window.SetChild(grid)

	number1 := gtk.NewEntry()
	number1.SetPlaceholderText("Enter Number")
	grid.Attach(number1, 0, 0, 1, 1)

	plusTitle := gtk.NewLabel(" + ")
	grid.Attach(plusTitle, 1, 0, 1, 1)

	number2 := gtk.NewEntry()
	number2.SetPlaceholderText("Enter Number")
	grid.Attach(number2, 2, 0, 1, 1)

	calculate := gtk.NewButtonWithLabel("Caluclate")
	grid.Attach(calculate, 3, 0, 1, 1)

	resultTitle := gtk.NewLabel("Result:")
	grid.Attach(resultTitle, 0, 1, 1, 1)

	resultValue := gtk.NewLabel("")
	grid.Attach(resultValue, 1, 1, 3, 1)

	calculate.Connect("clicked", func() {
		result := "Invalid Number"

		defer func() {
			resultValue.SetText(result)
		}()

		value1, err := strconv.Atoi(number1.Text())
		if err != nil {
			return
		}
		value2, err := strconv.Atoi(number2.Text())
		if err != nil {
			return
		}

		result = fmt.Sprintf("%d", value1+value2)
	})

	window.SetTitle("Addition App")
	window.Show()
}
