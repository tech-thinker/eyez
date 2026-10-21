package cmd

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tech-thinker/eyez/consts"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

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
	r, w, _ := os.Pipe()
	os.Stdout = w
	go io.Copy(io.Discard, r)
	defer func() {
		os.Stdout = oldStdout
		w.Close()
		r.Close()
	}()

	c := NewCommands(consts.GRAPHICS_ASCII, consts.ALGO_LANCZOS)
	path := createTempPNG(t)

	err := c.ByArgs(path, 10)
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

	// Test valid extension but invalid image data
	invalidImgPath := filepath.Join(t.TempDir(), "invalid.png")
	os.WriteFile(invalidImgPath, []byte("not an image"), 0644)
	err = c.ByArgs(invalidImgPath, 10)
	if err == nil {
		t.Errorf("ByArgs expected error for invalid image data, got nil")
	}

	// Test valid SVG file
	validSVGPath := filepath.Join(t.TempDir(), "test.svg")
	os.WriteFile(validSVGPath, []byte(`<svg viewBox="0 0 100 100"><circle cx="50" cy="50" r="40" fill="red"/></svg>`), 0644)
	err = c.ByArgs(validSVGPath, 10)
	if err != nil {
		t.Errorf("ByArgs expected nil error for valid SVG, got %v", err)
	}

	// Test invalid SVG data
	invalidSVGPath := filepath.Join(t.TempDir(), "invalid.svg")
	os.WriteFile(invalidSVGPath, []byte("not an svg"), 0644)
	err = c.ByArgs(invalidSVGPath, 10)
	if err == nil {
		t.Errorf("ByArgs expected error for invalid SVG data, got nil")
	}
}

func TestByStdin(t *testing.T) {
	oldStdout := os.Stdout
	rStd, wStd, _ := os.Pipe()
	os.Stdout = wStd
	go io.Copy(io.Discard, rStd)
	defer func() {
		os.Stdout = oldStdout
		wStd.Close()
		rStd.Close()
	}()

	c := NewCommands(consts.GRAPHICS_ASCII, consts.ALGO_LANCZOS)
	path := createTempPNG(t)

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("could not open test image: %v", err)
	}
	defer f.Close()

	err = c.ByStdin(f, 10)
	if err != nil {
		t.Errorf("ByStdin expected nil error, got %v", err)
	}

	// Test valid SVG via stdin
	svgPipeR, svgPipeW, _ := os.Pipe()
	svgPipeW.Write([]byte(`<svg viewBox="0 0 100 100"><rect width="100" height="100" fill="blue"/></svg>`))
	svgPipeW.Close()

	err = c.ByStdin(svgPipeR, 10)
	svgPipeR.Close()
	if err != nil {
		t.Errorf("ByStdin expected nil error for valid SVG, got %v", err)
	}

	// Test valid SVG with leading whitespace via stdin (triggers svg.Decode fallback when image.Decode fails)
	svgWsR, svgWsW, _ := os.Pipe()
	svgWsW.Write([]byte(`   <svg viewBox="0 0 100 100"><rect width="100" height="100" fill="green"/></svg>`))
	svgWsW.Close()

	err = c.ByStdin(svgWsR, 10)
	svgWsR.Close()
	if err != nil {
		t.Errorf("ByStdin expected nil error for SVG with leading whitespace, got %v", err)
	}

	// Test ByStdin with closed file to trigger io.ReadAll error
	closedR, closedW, _ := os.Pipe()
	closedR.Close()
	closedW.Close()

	err = c.ByStdin(closedR, 10)
	if err == nil {
		t.Errorf("ByStdin expected error for closed reader, got nil")
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

func TestDecodeImage(t *testing.T) {
	t.Run("Reader error", func(t *testing.T) {
		_, err := decodeImage(errorReader{}, "image.png")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("SVG with invalid content falls back to image.Decode then fails", func(t *testing.T) {
		_, err := decodeImage(strings.NewReader("not an svg"), "file.svg")
		if err == nil {
			t.Fatal("expected error for invalid SVG, got nil")
		}
	})

	t.Run("Fallback SVG decode succeeds when image.Decode does not recognize prefix", func(t *testing.T) {
		img, err := decodeImage(strings.NewReader("   <svg viewBox=\"0 0 50 50\"><circle cx=\"25\" cy=\"25\" r=\"20\" fill=\"red\"/></svg>"), "")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})
}

func TestNewCommandsInvalidGraphics(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		NewCommands("invalid_graphics", consts.ALGO_LANCZOS)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNewCommandsInvalidGraphics")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status != 0", err)
}

func TestNewCommandsInvalidAlgorithm(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		NewCommands(consts.GRAPHICS_ASCII, "invalid_algorithm")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestNewCommandsInvalidAlgorithm")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status != 0", err)
}
