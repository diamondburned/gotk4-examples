package main

import (
	"fmt"
	"os"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type ExampleApplication struct {
	gtk.Application
	TitlePrefix    glib.Property[string] `glib:"title-prefix"`
	QuitOnActivate glib.Property[bool]   `glib:"quit-on-activate"`

	initCalled        bool
	constructedCalled bool
	startupCalled     bool
	activateCalled    bool
	window            *gtk.ApplicationWindow
}

var exampleAppType = glib.RegisterSubclass[*ExampleApplication](
	glib.WithParamSpecs([]*glib.ParamSpec{
		glib.NewStringParam("title-prefix", "Title Prefix", "Prefix for window titles", "Default", glib.ParamReadWrite|glib.ParamConstruct),
		glib.NewBoolParam("quit-on-activate", "Quit on Activate", "Whether to quit the application after activation", false, glib.ParamReadWrite|glib.ParamConstruct),
	}),
	glib.WithOverrides(func(app *ExampleApplication) glib.ObjectOverrides {
		return glib.ObjectOverrides{
			Init:        func() { app.initCalled = true },
			Constructed: func() { app.constructedCalled = true },
		}
	}),
	glib.WithOverrides(func(app *ExampleApplication) gio.ApplicationOverrides {
		return gio.ApplicationOverrides{
			Startup:  app.startup,
			Activate: app.activate,
		}
	}),
)

func NewExampleApplication(appID string, prefix string) *ExampleApplication {
	return exampleAppType.NewWithProperties(map[string]any{
		"application-id": appID,
		"flags":          gio.ApplicationDefaultFlags,
		"title-prefix":   prefix,
	})
}

func (app *ExampleApplication) startup() {
	app.startupCalled = true
	// Chain up to parent GtkApplication/GApplication startup
	parent := glib.ParentOverrides[gio.ApplicationOverrides](app)
	if parent.Startup != nil {
		parent.Startup()
	}
}

func (app *ExampleApplication) activate() {
	app.activateCalled = true
	// Chain up to parent GtkApplication/GApplication activate
	parent := glib.ParentOverrides[gio.ApplicationOverrides](app)
	if parent.Activate != nil {
		parent.Activate()
	}

	// Create and present the application window
	app.window = gtk.NewApplicationWindow(&app.Application)
	app.window.SetTitle(fmt.Sprintf("%s: Subclassed GtkApplication", app.TitlePrefix.Get()))
	app.window.SetDefaultSize(400, 250)

	box := gtk.NewBox(gtk.OrientationVertical, 16)
	box.SetMarginTop(24)
	box.SetMarginBottom(24)
	box.SetMarginStart(24)
	box.SetMarginEnd(24)

	titleLabel := gtk.NewLabel("Pure Go GtkApplication Subclass!")
	titleLabel.AddCSSClass("title-1")

	count := 0
	btn := gtk.NewButtonWithLabel("Click me: 0")
	btn.ConnectClicked(func() {
		count++
		btn.SetLabel(fmt.Sprintf("Click me: %d", count))
	})

	box.Append(titleLabel)
	box.Append(btn)
	app.window.SetChild(box)

	app.window.Present()

	if app.QuitOnActivate.Get() {
		app.Quit()
	}
}

func main() {
	app := NewExampleApplication("com.example.SubclassApp", "MyApp")
	os.Exit(app.Run(os.Args))
}
