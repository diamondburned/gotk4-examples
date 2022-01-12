package gtkutil

import (
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func init() {
	go func() {
		for range time.Tick(30 * time.Second) {
			runtime.GC()
		}
	}()
}

var globalCSS strings.Builder
var cssOnce sync.Once

// AddCSS adds the given CSS into the global writer. This function must only be
// called during Go initialization.
func AddCSS(css string) {
	globalCSS.WriteString(css)
}

// LoadCSS loads the global CSS.
func LoadCSS() {
	cssOnce.Do(func() {
		content := globalCSS.String()

		prov := gtk.NewCSSProvider()
		prov.Connect("parsing-error", func(sec *gtk.CSSSection, err error) {
			// Optional line parsing routine.
			loc := sec.StartLocation()
			lines := strings.Split(content, "\n")
			log.Printf("CSS error (%v) at line: %q", err, lines[loc.Lines()])
		})
		prov.LoadFromData(content)

		// Load the CSS and apply it globally.
		gtk.StyleContextAddProviderForDisplay(
			gdk.DisplayGetDefault(), prov,
			gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
		)
	})
}
