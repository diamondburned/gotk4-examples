package main

import "github.com/diamondburned/gotk4/pkg/gtk/v3"

func init() {
	gtk.Init()
}

func main() {
	label := gtk.NewLabel("Hello from Go!")
	label.Show()

	window := gtk.NewWindow(gtk.WindowTypeToplevel)
	window.InitiallyUnowned.Connect("destroy", gtk.MainQuit)
	window.SetTitle("gotk4 Example")
	window.Add(label)
	window.SetDefaultSize(400, 300)
	window.Show()

	gtk.Main()
}
