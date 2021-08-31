package main

import (
	"os"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// IDs to access the tree view columns by
const (
	COLUMN_NAME = iota
	COLUMN_COMMENT
)

// ItemList is a thin wrapper around gtk.ListStore
type ItemList struct {
	*gtk.ListStore
}

// NewItemList create a new list of items
func NewItemList() *ItemList {
	listStore := gtk.NewListStore([]glib.Type{glib.TypeString, glib.TypeString})
	return &ItemList{listStore}
}

// Add adds a new item to the list
func (i *ItemList) Add(name, comment string) {
	iter := i.Append()

	i.Set(&iter,
		[]int32{COLUMN_NAME, COLUMN_COMMENT},
		[]glib.Value{*glib.NewValue(name), *glib.NewValue(comment)})
}

func createColumn(title string, id int32) *gtk.TreeViewColumn {
	cellRenderer := gtk.NewCellRendererText()
	column := gtk.NewTreeViewColumn()
	column.SetTitle(title)

	column.PackEnd(cellRenderer, false)
	column.AddAttribute(cellRenderer, "text", id)
	column.SetResizable(true)

	return column
}

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.treeview", 0)
	app.Connect("activate", activate)

	if code := app.Run(os.Args); code > 0 {
		os.Exit(int(code))
	}
}

func activate(app *gtk.Application) {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("Simple Treeview")
	win.SetDefaultSize(600, 300)

	treeView := gtk.NewTreeView()

	treeView.AppendColumn(createColumn("Name", COLUMN_NAME))
	treeView.AppendColumn(createColumn("Comment", COLUMN_COMMENT))

	items := NewItemList()
	treeView.SetModel(items)
	win.SetChild(&treeView.Widget)

	// Add some rows to the list store
	items.Add("hello", "Gtk4")

	win.Show()
}
