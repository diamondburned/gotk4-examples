package main

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type CustomButton struct {
	gtk.Button
	Number glib.Property[int] `glib:"number"`

	initCalled        bool
	constructedCalled bool
}

var customButtonType = glib.RegisterSubclass[*CustomButton](
	glib.WithParamSpecs([]*glib.ParamSpec{
		glib.NewInt64Param("number", "Number", "A number property", 0, 1000, 0, glib.ParamReadWrite|glib.ParamConstruct),
	}),
	glib.WithOverrides(func(b *CustomButton) glib.ObjectOverrides {
		return glib.ObjectOverrides{
			Init: func() {
				b.initCalled = true
			},
			Constructed: func() {
				b.constructedCalled = true
				b.Number.Notify(func(val int) {
					b.updateLabel()
				})
				b.updateLabel()
			},
		}
	}),
	glib.WithOverrides(func(b *CustomButton) gtk.ButtonOverrides {
		return gtk.ButtonOverrides{
			Clicked: func() {
				b.Number.Set(b.Number.Get() + 1)
			},
		}
	}),
)

func NewCustomButton() *CustomButton {
	return customButtonType.New()
}

func NewCustomButtonWithNumber(num int) *CustomButton {
	return customButtonType.NewWithProperties(map[string]any{
		"number": num,
	})
}

func (b *CustomButton) updateLabel() {
	b.SetLabel(fmt.Sprintf("Clicked %d times", b.Number.Get()))
}

func main() {
	app := gtk.NewApplication("com.example.custombutton", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() {
		win := gtk.NewApplicationWindow(app)
		win.SetTitle("Custom Button Example")
		win.SetDefaultSize(300, 120)

		btn := NewCustomButtonWithNumber(0)
		btn.Number.Notify(func(val int) {
			btn.updateLabel()
		})
		btn.updateLabel()

		btn.SetMarginTop(24)
		btn.SetMarginBottom(24)
		btn.SetMarginStart(24)
		btn.SetMarginEnd(24)

		win.SetChild(btn)
		win.Present()
	})

	os.Exit(app.Run(os.Args))
}
