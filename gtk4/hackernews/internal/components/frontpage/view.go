package frontpage

import (
	"context"
	"runtime"

	_ "embed"

	"github.com/diamondburned/adaptive"
	"github.com/diamondburned/gotk4-examples/gtk4/hackernews/internal/hackernews"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"golang.org/x/sync/semaphore"
)

// CommentsViewer describes the parent that can view posts.
type CommentsViewer interface {
	ViewComments(post *Post)
}

// View is the frontpage view.
type View struct {
	*gtk.ScrolledWindow
	ctx    context.Context
	parent CommentsViewer

	main *adaptive.LoadablePage
	view *gtk.Box

	// Keep track of the ongoing context so the refresh button can interrupt the
	// previous job.
	cancel context.CancelFunc
}

// NewView creates a new View instance.
func NewView(ctx context.Context, parent CommentsViewer) *View {
	v := View{
		ctx:    ctx,
		parent: parent,
	}

	v.main = adaptive.NewLoadablePage()

	viewport := gtk.NewViewport(nil, nil)
	viewport.SetScrollToFocus(true)
	viewport.SetChild(v.main)

	v.ScrolledWindow = gtk.NewScrolledWindow()
	v.ScrolledWindow.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	v.ScrolledWindow.SetChild(viewport)
	v.ScrolledWindow.SetPropagateNaturalHeight(true)

	return &v
}

// Refresh refreshes the whole view.
func (v *View) Refresh(done func()) {
	v.view = gtk.NewBox(gtk.OrientationVertical, 0)
	v.view.AddCSSClass("frontpage-posts")

	// Cancel the previous job.
	if v.cancel != nil {
		v.cancel()
	}

	// Create a cancellable context so we can cancel this fetching job if
	// needed by hitting the refresh button.
	ctx, cancel := context.WithCancel(v.ctx)
	v.cancel = cancel

	// loadingCtx is specifically used for loading up until we change the page
	// away from the spinning circle, since it's only used for the Stop button
	// there.
	loadingCtx := v.main.SetCancellableLoading(v.ctx)

	go func() {
		ids, err := hackernews.DefaultClient.Stories(loadingCtx, hackernews.TopStories)

		glib.IdleAdd(func() {
			if done != nil {
				done()
			}

			if err != nil {
				v.main.SetError(err)
			} else {
				v.addThenFetchPosts(ctx, ids)
			}
		})
	}()
}

func (v *View) addThenFetchPosts(ctx context.Context, ids []hackernews.ItemID) {
	v.main.SetChild(v.view)

	// Create a placeholder for the first 20 posts. Anything else is going to be
	// dynamically created later. This prevents UI stuttering on refresh.
	placeholders := 20
	if placeholders > len(ids) {
		placeholders = len(ids)
	}

	posts := make([]*Post, placeholders)
	for i := range posts {
		posts[i] = v.newPost(ids[i])
		v.view.Append(posts[i])
	}

	go fetchItems(ctx, ids, func(ix int, item *hackernews.Item, err error) {
		var post *Post
		if ix < placeholders {
			post = posts[ix]
		} else {
			post = v.newPost(item.ID)
			v.view.Append(post)
		}

		if err != nil {
			post.SetError(err)
		} else {
			post.SetItem(item)
		}
	})
}

func (v *View) newPost(id hackernews.ItemID) *Post {
	post := NewPost(id)
	post.ConnectShowComments(func() { v.parent.ViewComments(post) })
	return post
}

var httpConcurrency = int64(runtime.GOMAXPROCS(-1))

func fetchItems(
	ctx context.Context, ids []hackernews.ItemID, done func(int, *hackernews.Item, error)) {

	// There's no rate limiting so just do whatever we want.
	sema := semaphore.NewWeighted(httpConcurrency)

	for i, id := range ids {
		i := i
		id := id

		if err := sema.Acquire(ctx, 1); err != nil {
			return
		}
		go func() {
			defer sema.Release(1)

			item, err := hackernews.DefaultClient.Item(ctx, id)

			glib.IdleAdd(func() {
				select {
				case <-ctx.Done():
					// Do not add this item if the context is cancelled; we
					// might be inserting more items into a new list.
				default:
					done(i, item, err)
				}
			})
		}()
	}
}
