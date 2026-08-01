package main

import (
	core "github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type HasInitTemplate interface {
	InitTemplate()
}

func NewTemplate[T core.Objector, classT any](
	template_data []byte,
	widget_func func(classT) *gtk.WidgetClass,
	init_func func(),
	opts ...core.RegisterOptsFunc,
) *core.RegisteredSubclass[T] {
	opts1 := core.WithClassInit(func(class classT) {
		resource := glib.NewBytes(template_data)
		widget := widget_func(class)
		widget.SetTemplate(resource)
	})

	opts2 := core.WithOverrides(func(obj T) core.ObjectOverrides {
		templateable, ok := any(obj).(HasInitTemplate)

		if ok {
			return core.ObjectOverrides{
				Init: func() {
					templateable.InitTemplate()

					if init_func != nil {
						init_func()
					}
				},
			}
		} else {
			return core.ObjectOverrides{}
		}
	})

	opts = append(opts, opts1)
	opts = append(opts, opts2)
	class := core.RegisterSubclass[T](opts...)

	return class
}
