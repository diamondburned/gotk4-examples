package main

import (
	"image"
	"image/color"
	"os"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestMain(m *testing.M) {
	gtk.Init()
	os.Exit(m.Run())
}

func TestImagePaintable(t *testing.T) {
	img := generateGradient(200, 250)

	paintable := NewImagePaintable(img)
	if paintable == nil {
		t.Fatal("expected NewImagePaintable to return non-nil")
	}
	defer glib.Destroy(paintable)

	if !paintable.Type().IsA(gdk.GTypePaintable) {
		t.Fatal("expected ImagePaintable to implement GdkPaintable")
	}

	// Verify intrinsic dimensions match Go image bounds
	if w := paintable.IntrinsicWidth(); w != 200 {
		t.Fatalf("expected IntrinsicWidth 200, got %d", w)
	}
	if h := paintable.IntrinsicHeight(); h != 250 {
		t.Fatalf("expected IntrinsicHeight 250, got %d", h)
	}

	// Verify snapshotting
	snapshot := gtk.NewSnapshot()
	paintable.Snapshot(snapshot, 200, 250)

	// Verify gtk.Picture accepts ImagePaintable
	picture := gtk.NewPictureForPaintable(paintable)
	if picture == nil {
		t.Fatal("expected NewPictureForPaintable to return non-nil")
	}
	defer glib.Destroy(picture)

	// Test dynamic image update
	var sizeInvalidated, contentsInvalidated bool
	paintable.ConnectInvalidateSize(func() {
		sizeInvalidated = true
	})
	paintable.ConnectInvalidateContents(func() {
		contentsInvalidated = true
	})

	newImg := generateCheckerboard(150, 180, 15)
	paintable.SetImage(newImg)

	if !sizeInvalidated {
		t.Error("expected InvalidateSize signal to be emitted on SetImage")
	}
	if !contentsInvalidated {
		t.Error("expected InvalidateContents signal to be emitted on SetImage")
	}

	if w := paintable.IntrinsicWidth(); w != 150 {
		t.Fatalf("expected IntrinsicWidth 150 after update, got %d", w)
	}
	if h := paintable.IntrinsicHeight(); h != 180 {
		t.Fatalf("expected IntrinsicHeight 180 after update, got %d", h)
	}
}

func TestImagePaintableMemoryLeakFree(t *testing.T) {
	paintable := NewImagePaintable(nil)
	defer glib.Destroy(paintable)

	snapshot := gtk.NewSnapshot()
	defer glib.Destroy(snapshot)

	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	// Perform 1000 procedural image updates and snapshots
	for i := 0; i < 1000; i++ {
		img.SetRGBA(i%64, (i*2)%64, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		paintable.SetImage(img)
		paintable.Snapshot(snapshot, 64, 64)
	}

	runtime.GC()
	debug.FreeOSMemory()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	heapAllocMB := float64(m.HeapAlloc) / (1024 * 1024)
	if heapAllocMB > 10 {
		t.Fatalf("HeapAlloc unexpectedly high: %.2f MB", heapAllocMB)
	}
}
