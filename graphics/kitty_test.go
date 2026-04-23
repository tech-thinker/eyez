package graphics

import (
	"image"
	"os"
	"strings"
	"testing"
)

func TestKittyDraw(t *testing.T) {
	// Backup original environment variables to be a good neighbor
	origWindowID := os.Getenv("KITTY_WINDOW_ID")
	origTerm := os.Getenv("TERM")
	defer func() {
		os.Setenv("KITTY_WINDOW_ID", origWindowID)
		os.Setenv("TERM", origTerm)
	}()

	k := &Kitty{}
	// Use a 4x4 image, enough to process png encoding and kitty escape seqs
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))

	t.Run("Error When Environment Variables Blocked", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM", "xterm-color") // Something other than xterm-kitty

		err := k.Draw(img)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "kitty graphics not supported") {
			t.Errorf("expected kitty unsupport error, got: %v", err)
		}
	})

	t.Run("Passes When KITTY_WINDOW_ID is Set", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "123")
		os.Setenv("TERM", "")

		// We need to redirect standard output briefly so we don't spam the test runner
		// with kitty escape sequences.
		oldStdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		err := k.Draw(img)

		// Must close pipe and restore stout quickly to avoid deadlocking tests!
		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Errorf("expected no error when KITTY_WINDOW_ID is provided, got: %v", err)
		}
	})

	t.Run("Passes When TERM is xterm-kitty", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM", "xterm-kitty")

		oldStdout := os.Stdout
		_, w, _ := os.Pipe()
		os.Stdout = w

		err := k.Draw(img)

		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Errorf("expected no error when TERM is xterm-kitty, got: %v", err)
		}
	})
}
