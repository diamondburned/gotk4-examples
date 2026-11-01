package main

import (
	"os"
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func TestMain(m *testing.M) {
	// Initialize GTK in headless mode
	gtk.Init()
	os.Exit(m.Run())
}

func TestCustomButtonLifecycleAndConstruct(t *testing.T) {
	btn := NewCustomButtonWithNumber(42)

	if !btn.initCalled {
		t.Fatal("expected Init to be called")
	}
	if !btn.constructedCalled {
		t.Fatal("expected Constructed to be called")
	}
	if btn.Number.Get() != 42 {
		t.Fatalf("expected construct property number 42, got %d", btn.Number.Get())
	}
	if btn.Label() != "Clicked 42 times" {
		t.Fatalf("expected label 'Clicked 42 times' from Constructed, got %q", btn.Label())
	}
}

func TestCustomButtonOverridesAndNotify(t *testing.T) {
	btn := NewCustomButton()
	if btn.Number.Get() != 0 {
		t.Fatalf("expected initial number 0, got %d", btn.Number.Get())
	}

	var notifiedVal int
	btn.Number.Notify(func(val int) {
		notifiedVal = val
	})

	// Trigger clicked vfunc
	btn.Emit("clicked")

	if btn.Number.Get() != 1 {
		t.Fatalf("expected number 1, got %d", btn.Number.Get())
	}
	if btn.Label() != "Clicked 1 times" {
		t.Fatalf("expected label 'Clicked 1 times', got %q", btn.Label())
	}
	if notifiedVal != 1 {
		t.Fatalf("expected notify value 1, got %d", notifiedVal)
	}
}
