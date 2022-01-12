package frontpage

import (
	"fmt"
	"html"
	"net/url"

	_ "embed"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/hackernews"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
	"github.com/dustin/go-humanize"
)

// PostParent is the parent widget of the post.
type PostParent interface {
	gtk.Widgetter
	ViewComments(*Post)
}

// Post is a HackerNews post widget.
type Post struct {
	*adaptive.LoadablePage
	body *postBody
}

type postBody struct {
	*gtk.Box
	item *hackernews.Item

	score struct {
		*gtk.Box
		icon  *gtk.Image
		count *gtk.Label
	}

	info struct {
		*gtk.Box
		title *gtk.Label
		desc  *gtk.Label
	}

	right *commentButton
}

// NewPost creates a new post.
func NewPost(id hackernews.ItemID) *Post {
	p := Post{}
	p.body = newPostBody()

	p.LoadablePage = adaptive.NewLoadablePage()
	p.LoadablePage.ErrorPage.SetIconName("")
	p.SetTransitionDuration(65)
	p.AddCSSClass("frontpage-post")
	p.SetChild(p.body) // keep the height
	p.SetLoading()

	return &p
}

func newPostBody() *postBody {
	body := postBody{}
	body.score.icon = gtk.NewImageFromIconName("pan-up-symbolic")
	body.score.count = gtk.NewLabel("0")

	body.score.Box = gtk.NewBox(gtk.OrientationVertical, 0)
	body.score.AddCSSClass("post-score")
	body.score.SetHAlign(gtk.AlignCenter)
	body.score.SetVAlign(gtk.AlignCenter)
	body.score.Append(body.score.icon)
	body.score.Append(body.score.count)

	body.info.title = gtk.NewLabel("")
	body.info.title.AddCSSClass("post-title")
	body.info.title.SetHAlign(gtk.AlignStart)
	body.info.title.SetEllipsize(pango.EllipsizeEnd)
	body.info.title.SetWrap(true)
	body.info.title.SetWrapMode(pango.WrapWordChar)
	body.info.title.SetLines(2)
	body.info.title.SetXAlign(0)

	body.info.desc = gtk.NewLabel("")
	body.info.desc.AddCSSClass("post-description")
	body.info.desc.SetHAlign(gtk.AlignStart)
	body.info.desc.SetWrap(true)
	body.info.desc.SetWrapMode(pango.WrapWordChar)
	body.info.desc.SetXAlign(0)

	body.info.Box = gtk.NewBox(gtk.OrientationVertical, 0)
	body.info.Box.AddCSSClass("post-info")
	body.info.Box.SetHExpand(true)
	body.info.Append(body.info.title)
	body.info.Append(body.info.desc)

	body.right = newCommentButton()
	body.right.Hide()

	body.Box = gtk.NewBox(gtk.OrientationHorizontal, 0)
	body.Append(body.score)
	body.Append(body.info)
	body.Append(body.right)

	return &body
}

// ConnectShowComments connects f to be called when the user clicks the comments
// button. Without calling this, the comments button is hidden.
func (p *Post) ConnectShowComments(f func()) {
	p.body.right.Show()
	p.body.right.ConnectClicked(f)
}

// ID returns the Post's item ID.
func (p *Post) ID() hackernews.ItemID {
	if p.body.item == nil {
		return 0
	}
	return p.body.item.ID
}

// Item returns the post's HackerNews Item, or nil if the Post doesn't yet have
// one.
func (p *Post) Item() *hackernews.Item {
	return p.body.item
}

func (p *Post) SetItem(item *hackernews.Item) {
	p.SetChild(p.body)
	p.body.item = item
	p.body.update()
}

func (p *postBody) update() {
	p.info.title.SetTooltipText(p.item.Title)
	if p.item.URL != "" {
		p.info.title.SetMarkup(fmt.Sprintf(
			`<a href="%s">%s</a>`,
			html.EscapeString(p.item.URL),
			html.EscapeString(p.item.Title),
		))
	} else {
		p.info.title.SetText(p.item.Title)
	}

	head := ""
	if u, err := url.Parse(p.item.URL); err == nil && u.Host != "" {
		head = u.Host + " | "
	}

	p.info.desc.SetMarkup(head + fmt.Sprintf(
		`by %s %s`,
		html.EscapeString(p.item.By),
		humanize.Time(p.item.Time.Time()),
	))

	p.score.count.SetText(fmt.Sprint(p.item.Score))
	p.right.SetCount(p.item.Descendants)
}

// SetError sets the error for the post while it's fetched.
func (p *Post) SetError(err error) {
	label := gtk.NewLabel(fmt.Sprintf(
		`<span color="red">Error:</span> %s`,
		html.EscapeString(err.Error()),
	))
	label.SetUseMarkup(true)

	p.SetErrorWidget(label)
}

type commentButton struct {
	*gtk.Button
	box   *gtk.Box
	icon  *gtk.Image
	count *gtk.Label
}

func newCommentButton() *commentButton {
	b := commentButton{}
	b.icon = gtk.NewImageFromIconName("mail-unread-symbolic")
	b.count = gtk.NewLabel("0")

	b.box = gtk.NewBox(gtk.OrientationVertical, 0)
	b.box.SetHAlign(gtk.AlignCenter)
	b.box.SetVAlign(gtk.AlignCenter)
	b.box.Append(b.icon)
	b.box.Append(b.count)

	b.Button = gtk.NewButton()
	b.Button.AddCSSClass("post-comment-button")
	b.Button.SetHasFrame(false)
	b.Button.SetChild(b.box)

	return &b
}

func (b *commentButton) SetCount(count int) {
	b.count.SetText(fmt.Sprint(count))
}
