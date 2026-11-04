package main

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

type CustomAppWindow struct {
	gtk.ApplicationWindow
}

func CustomAppWindow_GetWidgetClass(class *gtk.ApplicationWindowClass) *gtk.WidgetClass {
	window := class.ParentClass()
	widget := window.ParentClass()
	return widget
}

var CustomAppWindowSubclass = NewTemplate[*CustomAppWindow, *gtk.ApplicationWindowClass](Embedded_AppUI, CustomAppWindow_GetWidgetClass, nil)

func NewCustomAppWindow() *CustomAppWindow {
	return CustomAppWindowSubclass.New()
}
