package main

import (
	_ "embed"
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/labstack/gommon/log"
	"os"
)

//go:embed listview.ui
var windowXML string

type demoListItem struct {
	name   string
	number int
}

type listHelper struct {
	stringList *gtk.StringList
	listItems  map[string]*demoListItem
	factory    *gtk.SignalListItemFactory
	selection  *gtk.SingleSelection
	itemSerial int
}

func newListHelper(listView *gtk.ListView) *listHelper {
	stringList := gtk.NewStringList([]string{})
	helper := &listHelper{
		stringList: stringList,
		listItems:  make(map[string]*demoListItem),
		factory:    gtk.NewSignalListItemFactory(),
		selection:  gtk.NewSingleSelection(stringList),
		itemSerial: 0,
	}
	listView.SetFactory(&helper.factory.ListItemFactory)
	helper.factory.ConnectSetup(helper.Setup)
	helper.factory.ConnectBind(helper.Bind)
	listView.SetModel(helper.selection)
	return helper
}

func (l *listHelper) AddItem(listItemMetadata *demoListItem) {
	l.itemSerial += 1
	key := fmt.Sprintf("item:%d", l.itemSerial)
	l.listItems[key] = listItemMetadata

	// Append to stringList last to avoid a race condition where listItmes[key] is missing when Bind() runs
	l.stringList.Append(key)
	log.Infof("Added item: %s", listItemMetadata.name)
}

func (l *listHelper) Bind(listItem *gtk.ListItem) {
	idx := listItem.Position()
	var listItemMetadata *demoListItem
	value := l.stringList.String(idx)
	listItemMetadata = l.listItems[value]

	labelText := "(blank)"
	if listItemMetadata != nil {
		labelText = fmt.Sprintf("%s: %d", listItemMetadata.name, listItemMetadata.number)
	}

	label := listItem.Child().(*gtk.Label)
	label.SetLabel(labelText)
}

func (l *listHelper) Setup(listItem *gtk.ListItem) {
	listItem.SetChild(gtk.NewLabel("unbound"))
}

func main() {
	app := gtk.NewApplication("com.github.gotk4examples.listview", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() {
		builder := gtk.NewBuilderFromString(windowXML, len(windowXML))
		window := builder.GetObject("appWindow").Cast().(*gtk.ApplicationWindow)
		window.SetApplication(app)
		window.Show()

		listViewInstance := builder.GetObject("myListView").Cast().(*gtk.ListView)

		myHelper := newListHelper(listViewInstance)
		myHelper.AddItem(
			&demoListItem{
				name:   "First item",
				number: 1,
			})
		myHelper.AddItem(
			&demoListItem{
				name:   "Second item",
				number: 2,
			})

	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
