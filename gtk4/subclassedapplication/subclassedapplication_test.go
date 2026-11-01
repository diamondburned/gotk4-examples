package main

import (
	"testing"

	"github.com/diamondburned/gotk4/pkg/core/glib"
)

func TestExampleApplicationLifecycle(t *testing.T) {
	app := NewExampleApplication("com.example.SubclassAppTest", "MyApp")
	defer glib.Destroy(app)

	if app == nil {
		t.Fatal("expected NewExampleApplication to return non-nil")
	}
	if !app.initCalled {
		t.Fatal("expected Init override to be called")
	}
	if !app.constructedCalled {
		t.Fatal("expected Constructed override to be called")
	}
	if app.TitlePrefix.Get() != "MyApp" {
		t.Fatalf("expected TitlePrefix 'MyApp', got %q", app.TitlePrefix.Get())
	}

	app.QuitOnActivate.Set(true)

	code := app.Run([]string{"subclassedapplication"})
	if code != 0 {
		t.Fatalf("expected app.Run to return 0, got %d", code)
	}

	if !app.startupCalled {
		t.Fatal("expected Startup override to be called")
	}
	if !app.activateCalled {
		t.Fatal("expected Activate override to be called")
	}
	if app.window == nil {
		t.Fatal("expected window to be created in Activate")
	}
	if app.window.Title() != "MyApp: Subclassed GtkApplication" {
		t.Fatalf("expected window title 'MyApp: Subclassed GtkApplication', got %q", app.window.Title())
	}
}
