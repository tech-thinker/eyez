package cmd

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/tech-thinker/eyez/consts"
)

// createTempPNG is a helper to create a small PNG image for testing
func createTempPNG(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test_image.png")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	defer f.Close()

	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	err = png.Encode(f, img)
	if err != nil {
		t.Fatalf("failed to encode png: %v", err)
	}
	return path
}

func TestNewCommandsValid(t *testing.T) {
	// Let's test all valid permutations
	graphicsOpts := []string{consts.GRAPHICS_UNICODE, consts.GRAPHICS_ASCII, consts.GRAPHICS_KITTY}
	algoOpts := []string{consts.ALGO_CATMULL_ROM, consts.ALGO_LANCZOS}

	for _, g := range graphicsOpts {
		for _, a := range algoOpts {
			c := NewCommands(g, a)
			if c == nil {
				t.Errorf("NewCommands(%s, %s) returned nil", g, a)
			}
		}
	}
}

func TestByArgs(t *testing.T) {
	// Temporarily hijack stdout to prevent printing garbage since graphics algorithms dump to terminal
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	c := NewCommands(consts.GRAPHICS_ASCII, consts.ALGO_LANCZOS)
	path := createTempPNG(t)

	err := c.ByArgs(path, 10)

	os.Stdout = oldStdout
	w.Close()

	if err != nil {
		t.Errorf("ByArgs expected nil error, got %v", err)
	}

	// Test invalid file (but correct extension) to test if parsing files catches missing IO errors
	err = c.ByArgs("non_existent_file.png", 10)
	if err == nil {
		t.Errorf("ByArgs expected error for missing file, got nil")
	}

	// Test unsupported file format extension checks from validator
	err = c.ByArgs("file.txt", 10)
	if err == nil {
		t.Errorf("ByArgs expected error for unsupported format, got nil")
	}
}

func TestByStdin(t *testing.T) {
	oldStdout := os.Stdout
	_, wStd, _ := os.Pipe()
	os.Stdout = wStd

	c := NewCommands(consts.GRAPHICS_ASCII, consts.ALGO_LANCZOS)
	path := createTempPNG(t)

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("could not open test image: %v", err)
	}
	defer f.Close()

	err = c.ByStdin(f, 10)

	os.Stdout = oldStdout
	wStd.Close()

	if err != nil {
		t.Errorf("ByStdin expected nil error, got %v", err)
	}

	// Pass invalid image data to deduce image decode error
	r, w, _ := os.Pipe()
	w.Write([]byte("not an image"))
	w.Close()

	err = c.ByStdin(r, 10)
	r.Close()

	if err == nil {
		t.Errorf("ByStdin expected error for invalid image data, got nil")
	}
}
