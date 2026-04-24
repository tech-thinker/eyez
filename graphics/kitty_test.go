package graphics

import (
	"image"
	"math/rand"
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

		// Must close pipe and restore stdout quickly to avoid deadlocking tests!
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

	t.Run("Passes When Large Image Requires Multiple Chunks", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "123")
		os.Setenv("TERM", "")

		// Create a large enough image to ensure it requires multiple chunks (base64 > 4096 bytes)
		largeImg := image.NewRGBA(image.Rect(0, 0, 500, 500))
		for i := 0; i < len(largeImg.Pix); i++ {
			// Add some pseudo-randomness to prevent high PNG compression
			largeImg.Pix[i] = uint8((i * 137) % 256)
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Read continuously to prevent deadlock if pipe buffer fills up
		go func() {
			var discard [1024]byte
			for {
				_, err := r.Read(discard[:])
				if err != nil {
					break
				}
			}
		}()

		err := k.Draw(largeImg)

		w.Close()
		r.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Errorf("expected no error when writing large image, got: %v", err)
		}
	})

	t.Run("Fails on Encode Error", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "123")
		os.Setenv("TERM", "")

		// 0x0 image will cause png.Encode to fail
		err := k.Draw(image.NewRGBA(image.Rect(0, 0, 0, 0)))
		if err == nil {
			t.Errorf("expected error on zero sized image, got nil")
		}
	})

	t.Run("Fails When First Write Errors", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "123")
		os.Setenv("TERM", "")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Close the read end immediately
		r.Close()

		// Large image with noise forces Fprintf to overflow its bufio buffer and flush immediately on the first chunk
		largeImg := image.NewRGBA(image.Rect(0, 0, 500, 500))
		for i := 0; i < len(largeImg.Pix); i++ {
			largeImg.Pix[i] = uint8((i * 137) % 256)
		}
		err := k.Draw(largeImg)

		w.Close()
		os.Stdout = oldStdout

		if err == nil {
			t.Errorf("expected error when writing to closed pipe on first chunk, got nil")
		}
	})

	t.Run("Fails When Second Write Errors", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "123")
		os.Setenv("TERM", "")

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Close the read end asynchronously to allow the first chunk to succeed
		go func() {
			buf := make([]byte, 4096)
			r.Read(buf)
			r.Close()
		}()

		// Large noisy image (> 64KB base64) to fill the pipe buffer and block until closed
		largeImg := image.NewRGBA(image.Rect(0, 0, 500, 500))
		for i := 0; i < len(largeImg.Pix); i++ {
			largeImg.Pix[i] = uint8(rand.Intn(256))
		}

		err := k.Draw(largeImg)

		w.Close()
		os.Stdout = oldStdout

		if err == nil {
			t.Errorf("expected error when writing to closed pipe on subsequent chunks, got nil")
		}
	})
}
