package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	glibv2 "github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// ImagePaintable is a custom GdkPaintable implemented purely in Go.
// It wraps a Go standard library image.Image and adapts it for GTK4.
type ImagePaintable struct {
	glib.Object
	*gdk.Paintable

	image   image.Image
	texture *gdk.Texture
}

var imagePaintableType = glib.RegisterSubclass[*ImagePaintable](
	glib.WithOverrides(func(p *ImagePaintable) glib.ObjectOverrides {
		return glib.ObjectOverrides{
			Init: func() {
				p.Paintable = &gdk.Paintable{Object: &p.Object}
			},
			Dispose: func() {
				p.image = nil
				p.texture = nil
			},
		}
	}),
	glib.WithOverrides(func(p *ImagePaintable) gdk.PaintableOverrides {
		return gdk.PaintableOverrides{
			IntrinsicWidth: func() int {
				if p.image == nil {
					return 0
				}
				return p.image.Bounds().Dx()
			},
			IntrinsicHeight: func() int {
				if p.image == nil {
					return 0
				}
				return p.image.Bounds().Dy()
			},
			Snapshot: func(snapshot gdk.Snapshotter, width, height float64) {
				if p.texture != nil {
					p.texture.Snapshot(snapshot, width, height)
				}
			},
		}
	}),
)

// NewImagePaintable creates a new ImagePaintable wrapping a Go image.Image.
func NewImagePaintable(img image.Image) *ImagePaintable {
	p := imagePaintableType.New()
	p.SetImage(img)
	return p
}

// Image returns the current Go image.Image.
func (p *ImagePaintable) Image() image.Image {
	return p.image
}

// Texture returns the underlying GdkTexture.
func (p *ImagePaintable) Texture() *gdk.Texture {
	return p.texture
}

// SetImage updates the Go image.Image and notifies GTK of size and content changes.
func (p *ImagePaintable) SetImage(img image.Image) {
	hadImage := p.image != nil
	p.image = img

	if img != nil {
		rgba := imageToRGBA(img)
		bytes := glibv2.NewBytes(rgba.Pix)
		p.texture = &gdk.NewMemoryTexture(
			rgba.Bounds().Dx(),
			rgba.Bounds().Dy(),
			gdk.MemoryR8G8B8A8Premultiplied,
			bytes,
			uint(rgba.Stride),
		).Texture
	} else {
		p.texture = nil
	}

	if hadImage || img != nil {
		p.InvalidateSize()
		p.InvalidateContents()
	}
}

// imageToRGBA converts an arbitrary image.Image to an *image.RGBA with origin at (0, 0).
func imageToRGBA(src image.Image) *image.RGBA {
	if rgba, ok := src.(*image.RGBA); ok && rgba.Rect.Min == (image.Point{}) && len(rgba.Pix) == rgba.Stride*rgba.Rect.Dy() {
		return rgba
	}
	bounds := src.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Rect, src, bounds.Min, draw.Src)
	return rgba
}

// generateGradient creates a procedural Go *image.RGBA image.
func generateGradient(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8(float64(x) / float64(width) * 255)
			g := uint8(float64(y) / float64(height) * 255)
			b := uint8(180)
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	return img
}

// generateCheckerboard creates a procedural checkerboard *image.RGBA.
func generateCheckerboard(width, height, cellSize int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	c1 := color.RGBA{R: 50, G: 130, B: 220, A: 255}
	c2 := color.RGBA{R: 240, G: 240, B: 245, A: 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if ((x/cellSize)+(y/cellSize))%2 == 0 {
				img.SetRGBA(x, y, c1)
			} else {
				img.SetRGBA(x, y, c2)
			}
		}
	}
	return img
}

func main() {
	app := gtk.NewApplication("com.example.custompaintable", gio.ApplicationDefaultFlags)
	app.ConnectActivate(func() {
		win := gtk.NewApplicationWindow(app)
		win.SetTitle("Custom Paintable (Go image.Image -> GTK)")
		win.SetDefaultSize(400, 480)

		box := gtk.NewBox(gtk.OrientationVertical, 16)
		box.SetMarginTop(20)
		box.SetMarginBottom(20)
		box.SetMarginStart(20)
		box.SetMarginEnd(20)

		// 1. Generate an initial procedural Go image
		initialImg := generateGradient(256, 256)

		// 2. Wrap Go image.Image inside our custom GdkPaintable subclass
		paintable := NewImagePaintable(initialImg)

		// 3. Display in a standard GTK4 Picture widget
		picture := gtk.NewPictureForPaintable(paintable)
		picture.SetCanShrink(true)
		picture.SetContentFit(gtk.ContentFitContain)
		picture.SetVExpand(true)

		infoLabel := gtk.NewLabel(fmt.Sprintf(
			"Image: %dx%d (%T)",
			initialImg.Bounds().Dx(), initialImg.Bounds().Dy(), initialImg,
		))
		infoLabel.AddCSSClass("title-4")

		btnBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
		btnBox.SetHAlign(gtk.AlignCenter)

		gradientBtn := gtk.NewButtonWithLabel("Gradient")
		gradientBtn.ConnectClicked(func() {
			grad := generateGradient(256, 256)
			paintable.SetImage(grad)
			infoLabel.SetText(fmt.Sprintf("Image: %dx%d (%T)", grad.Bounds().Dx(), grad.Bounds().Dy(), grad))
		})

		checkerBtn := gtk.NewButtonWithLabel("Checkerboard")
		checkerBtn.ConnectClicked(func() {
			checker := generateCheckerboard(256, 256, 32)
			paintable.SetImage(checker)
			infoLabel.SetText(fmt.Sprintf("Image: %dx%d (%T)", checker.Bounds().Dx(), checker.Bounds().Dy(), checker))
		})

		btnBox.Append(gradientBtn)
		btnBox.Append(checkerBtn)

		box.Append(picture)
		box.Append(infoLabel)
		box.Append(btnBox)

		win.SetChild(box)
		win.Present()
	})

	os.Exit(app.Run(os.Args))
}
