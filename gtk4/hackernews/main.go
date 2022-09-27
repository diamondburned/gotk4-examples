package main

import (
	"context"
	"os"
	"os/signal"

	_ "embed"

	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/components/frontpage"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/components/postview"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/gtkutil"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.hackernews", gio.ApplicationFlags(gio.ApplicationFlagsNone))
	app.ConnectActivate(func() {
		gtkutil.LoadCSS()
		hn := newHackerNews(ctx, app)
		hn.Show()
	})

	go func() {
		<-ctx.Done()
		glib.IdleAdd(app.Quit)
	}()

	if code := app.Run(os.Args); code > 0 {
		cancel()
		os.Exit(code)
	}
}

type hackerNews struct {
	*gtk.Application
	window *gtk.ApplicationWindow
	header *gtk.HeaderBar

	back    *gtk.Button
	refresh *gtk.Button

	view struct {
		*gtk.Stack
		front *frontpage.View
		post  *postview.View
	}
}

func newHackerNews(ctx context.Context, app *gtk.Application) *hackerNews {
	hn := hackerNews{Application: app}
	hn.view.front = frontpage.NewView(ctx, &hn)
	hn.view.front.Refresh(nil)

	hn.view.Stack = gtk.NewStack()
	hn.view.SetTransitionType(gtk.StackTransitionTypeSlideLeftRight)
	hn.view.AddChild(hn.view.front)

	hn.back = gtk.NewButtonFromIconName("go-previous-symbolic")
	hn.back.Hide()
	hn.back.SetTooltipText("Back")
	hn.back.ConnectClicked(hn.ViewPosts)

	hn.refresh = gtk.NewButtonFromIconName("view-refresh-symbolic")
	hn.refresh.SetTooltipText("Refresh")
	hn.refresh.ConnectClicked(func() {
		hn.refresh.SetSensitive(false)
		hn.view.front.Refresh(func() { hn.refresh.SetSensitive(true) })
	})

	hn.header = gtk.NewHeaderBar()
	hn.header.PackStart(hn.back)
	hn.header.PackEnd(hn.refresh)

	hn.window = gtk.NewApplicationWindow(app)
	hn.window.SetDefaultSize(700, 550)
	hn.window.SetChild(hn.view)
	hn.window.SetTitle("HackerNews")
	hn.window.SetTitlebar(hn.header)

	return &hn
}

// Show shows the HackerNews window.
func (hn *hackerNews) Show() {
	hn.window.Show()
}

// ViewComments shows the given post and its comments.
func (hn *hackerNews) ViewComments(post *frontpage.Post) {
	if hn.view.post != nil {
		hn.view.post.Unparent()
		hn.view.post = nil
	}

	hn.view.post = postview.NewView(post.Item())
	hn.view.AddChild(hn.view.post)
	hn.view.SetVisibleChild(hn.view.post)

	hn.window.SetTitle(post.Item().Title + " — HackerNews")

	hn.back.Show()
	hn.refresh.Hide()
}

// ViewPosts switches to the frontpage.
func (hn *hackerNews) ViewPosts() {
	hn.view.SetVisibleChild(hn.view.front)

	hn.window.SetTitle("HackerNews")

	hn.back.Hide()
	hn.refresh.Show()

	if hn.view.post != nil {
		hn.view.Remove(hn.view.post)
		hn.view.post = nil
	}
}
