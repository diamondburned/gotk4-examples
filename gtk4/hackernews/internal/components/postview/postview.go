package postview

import (
	"context"
	"runtime"

	_ "embed"

	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/components/frontpage"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/gtkutil"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/hackernews"
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"golang.org/x/sync/semaphore"
)

//go:embed postview.css
var css string

func init() { gtkutil.AddCSS(css) }

// CommentAppender describes any widget that we can add comments to. It is
// implemented by View and Comment.
type CommentAppender interface {
	Append(*Comment)
}

var (
	_ CommentAppender = (*Comment)(nil)
	_ CommentAppender = (*View)(nil)
)

// View describes a view of a single post with all its comments.
type View struct {
	*gtk.ScrolledWindow
	item *hackernews.Item
	box  *gtk.Box
	post *frontpage.Post
}

// NewView creates a new View.
func NewView(item *hackernews.Item) *View {
	v := View{item: item}
	v.post = frontpage.NewPost(item.ID)
	v.post.SetItem(item)

	v.box = gtk.NewBox(gtk.OrientationVertical, 0)
	v.box.AddCSSClass("postview-post")
	v.box.Append(v.post)

	v.ScrolledWindow = gtk.NewScrolledWindow()
	v.ScrolledWindow.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	v.ScrolledWindow.SetChild(v.box)

	// Make a context that's cancelled when the view is swapped out with
	// something else.
	ctx, cancel := context.WithCancel(context.Background())
	v.ConnectUnmap(cancel)

	fetcher := commentsFetcher{
		ctx:  ctx,
		sema: semaphore.NewWeighted(int64(runtime.GOMAXPROCS(-1))),
	}
	go fetcher.fetch(&v, item.Kids)

	return &v
}

type commentsFetcher struct {
	ctx  context.Context
	sema *semaphore.Weighted
}

func (f *commentsFetcher) fetch(parent CommentAppender, ids []hackernews.ItemID) {
	if err := f.sema.Acquire(f.ctx, 1); err != nil {
		return
	}
	defer f.sema.Release(1)

	for _, id := range ids {
		select {
		case <-f.ctx.Done():
			return
		default:
		}

		item, err := hackernews.DefaultClient.Item(f.ctx, id)
		if f.ctx.Err() != nil {
			// Context expired; don't do anything.
			return
		}

		if item != nil && item.Text == "" {
			// The comment's content is empty for some reason. Ignore it.
			return
		}

		glib.IdleAdd(func() {
			comment := NewComment()
			parent.Append(comment)

			if err != nil {
				comment.SetError(err)
				return
			}

			comment.SetItem(item)

			if len(item.Kids) > 0 {
				go f.fetch(comment, item.Kids)
			}
		})
	}
}

// type commentFetchJob struct {
// 	parent CommentAppender
// 	kids   []hackernews.ItemID
// }

// func (v *View) fetch(ctx context.Context, ch chan commentFetchJob) {
// 	var jobs []commentFetchJob
// 	var job commentFetchJob

// 	for {
// 		if len(jobs) > 0 {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case ch <- jobs[0]:
// 				jobs = append(jobs[:0], jobs[1:]...) // pop first
// 			case job = <-ch:
// 			}
// 		} else {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case job = <-ch:
// 			}
// 		}

// 		for _, id := range job.kids {
// 			item, err := hackernews.DefaultClient.Item(ctx, id)
// 		}
// 	}
// }

// func fetchCommentsCh(ch <-chan commentFetchJob) {
// }

// Append adds a comment to the toplevel.
func (v *View) Append(comment *Comment) {
	v.box.Append(comment)
}
