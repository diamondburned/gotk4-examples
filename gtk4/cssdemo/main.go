package main

import (
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
)

//go:embed style.css
var styleCSS string

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.cssdemo", 0)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	// Load the CSS and apply it globally.
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(), loadCSS(styleCSS),
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	title := gtk.NewLabel("Woah!")
	title.AddCSSClass("title")
	title.SetWrap(true)
	title.SetWrapMode(pango.WrapWordChar)
	title.SetXAlign(0)
	title.SetYAlign(0)

	subtitle := gtk.NewLabel("This looks like an unstyled HTML page! How?")
	subtitle.AddCSSClass("subtitle")
	subtitle.SetWrap(true)
	subtitle.SetWrapMode(pango.WrapWordChar)
	subtitle.SetXAlign(0)
	subtitle.SetYAlign(0)

	close := gtk.NewButtonWithLabel("Close")
	close.SetHAlign(gtk.AlignStart)

	box := gtk.NewBox(gtk.OrientationVertical, 0)
	box.Append(title)
	box.Append(subtitle)
	box.Append(close)

	window := gtk.NewApplicationWindow(app)
	window.SetTitle("Anti-browser")
	window.SetChild(box)
	window.SetTitlebar(gtk.NewBox(gtk.OrientationVertical, 0)) // hide headerbar
	window.SetDefaultSize(400, 300)
	window.Show()

	close.ConnectClicked(window.Close)
}

func loadCSS(content string) *gtk.CSSProvider {
	prov := gtk.NewCSSProvider()
	prov.ConnectParsingError(func(sec *gtk.CSSSection, err error) {
		// Optional line parsing routine.
		loc := sec.StartLocation()
		lines := strings.Split(content, "\n")
		log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
	})
	prov.LoadFromData(content)
	return prov
}
