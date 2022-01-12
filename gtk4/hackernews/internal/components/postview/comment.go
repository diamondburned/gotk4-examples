package postview

import (
	"fmt"
	"html"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/hackernews"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/hackernews/hnhtml"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
	"github.com/dustin/go-humanize"
)

// Comment describes a single comment. A comment can contain more comments.
type Comment struct {
	*adaptive.LoadablePage
	item  *hackernews.Item
	child *childComments

	box *gtk.Box
	top struct {
		*gtk.Box
		info *gtk.Label
		rev  *gtk.ToggleButton
	}
	text *gtk.Label
}

const (
	commentRevealedIcon  = "pan-down-symbolic"
	commentCollapsedIcon = "pan-end-symbolic"
)

// NewComment creates a new Comment.
func NewComment() *Comment {
	c := Comment{}

	c.top.info = gtk.NewLabel("")
	c.top.info.SetXAlign(0)
	c.top.info.SetEllipsize(pango.EllipsizeEnd)
	c.top.info.SetSingleLineMode(true)

	c.top.Box = gtk.NewBox(gtk.OrientationHorizontal, 0)
	c.top.AddCSSClass("postview-comment-top")
	c.top.Append(c.top.info)

	c.text = gtk.NewLabel("")
	c.text.AddCSSClass("postview-comment-text")
	c.text.SetSelectable(true)
	c.text.SetXAlign(0)
	c.text.SetWrap(true)
	c.text.SetWrapMode(pango.WrapWordChar)

	c.box = gtk.NewBox(gtk.OrientationVertical, 0)
	c.box.AddCSSClass("postview-comment")
	c.box.Append(c.top)
	c.box.Append(c.text)

	c.LoadablePage = adaptive.NewLoadablePage()
	c.LoadablePage.ErrorPage.SetIconName("")
	c.AddCSSClass("postview-comment-load")
	c.SetTransitionDuration(65)
	c.SetChild(c.box)
	c.SetLoading()

	return &c
}

// SetError marks the comment as erroneous and shows the error to the user.
func (c *Comment) SetError(err error) {
	label := gtk.NewLabel(fmt.Sprintf(
		`<span color="red">Error:</span> %s`,
		html.EscapeString(err.Error()),
	))
	label.SetUseMarkup(true)

	c.LoadablePage.SetErrorWidget(label)
}

// SetItem sets the fetched item and displays it.
func (c *Comment) SetItem(item *hackernews.Item) {
	c.LoadablePage.SetChild(c.box)
	c.item = item
	c.update()
}

func (c *Comment) update() {
	c.text.SetText(hnhtml.ToMarkup(c.item.Text))
	c.top.info.SetMarkup(html.EscapeString(c.item.By) + " " + humanize.Time(c.item.Time.Time()))
}

// Append adds the given comment as the comment's children.
func (c *Comment) Append(comment *Comment) {
	if c.child == nil {
		c.child = newChildComments()
		c.child.Revealer.Connect("notify::reveal-child", func() {
			if c.child.RevealChild() {
				c.top.rev.SetIconName(commentRevealedIcon)
			} else {
				c.top.rev.SetIconName(commentCollapsedIcon)
			}
		})
		c.box.Append(c.child)

		c.top.rev = gtk.NewToggleButton()
		c.top.rev.AddCSSClass("postview-comment-reveal")
		c.top.rev.SetHasFrame(false)
		c.top.rev.SetIconName(commentRevealedIcon)
		c.top.rev.ConnectClicked(func() {
			c.child.SetRevealChild(!c.child.RevealChild())
		})
		c.top.Append(c.top.rev)
	}

	c.child.Append(comment)
	c.update()
}

type childComments struct {
	*gtk.Revealer
	box *gtk.Box
}

func newChildComments() *childComments {
	c := childComments{}
	c.box = gtk.NewBox(gtk.OrientationVertical, 0)
	c.box.AddCSSClass("postview-comment-children")

	c.Revealer = gtk.NewRevealer()
	c.Revealer.SetTransitionType(gtk.RevealerTransitionTypeSlideDown)
	c.Revealer.SetChild(c.box)
	c.Revealer.SetRevealChild(true)

	return &c
}

func (c *childComments) Append(comment *Comment) {
	c.box.Append(comment)
}
