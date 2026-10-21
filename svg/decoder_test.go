package svg

import (
	"image"
	"strings"
	"testing"
)

const sampleSVG = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
  <rect width="100" height="100" fill="blue"/>
  <circle cx="50" cy="50" r="30" fill="yellow"/>
</svg>`

const minimalSVG = `<svg viewBox="0 0 50 50">
  <line x1="0" y1="0" x2="50" y2="50" stroke="black"/>
</svg>`

const dimensionSVG = `<svg width="200" height="100">
  <rect width="200" height="100" fill="green"/>
</svg>`

const noDimensionsSVG = `<svg>
  <circle cx="25" cy="25" r="20" fill="red"/>
</svg>`

const largeSVG = `<svg viewBox="0 0 4000 2000">
  <rect width="4000" height="2000" fill="purple"/>
</svg>`

func TestDecode(t *testing.T) {
	t.Run("Valid XML SVG", func(t *testing.T) {
		img, err := Decode(strings.NewReader(sampleSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
		rgba, ok := img.(*image.RGBA)
		if !ok {
			t.Fatalf("expected *image.RGBA, got %T", img)
		}
		// Scaled up because original is 100x100 (< 512)
		if rgba.Bounds().Dx() < 100 || rgba.Bounds().Dy() < 100 {
			t.Errorf("image bounds too small: %v", rgba.Bounds())
		}
	})

	t.Run("Minimal SVG without XML header", func(t *testing.T) {
		img, err := Decode(strings.NewReader(minimalSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})

	t.Run("SVG with width and height", func(t *testing.T) {
		img, err := Decode(strings.NewReader(dimensionSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})

	t.Run("SVG with no dimensions", func(t *testing.T) {
		img, err := Decode(strings.NewReader(noDimensionsSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
		if img.Bounds().Dx() != 512 || img.Bounds().Dy() != 512 {
			t.Errorf("expected 512x512 fallback, got %v", img.Bounds())
		}
	})

	t.Run("Large SVG capped dimensions", func(t *testing.T) {
		img, err := Decode(strings.NewReader(largeSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if img.Bounds().Dx() > 2048 || img.Bounds().Dy() > 2048 {
			t.Errorf("expected dimensions capped to <= 2048, got %v", img.Bounds())
		}
	})

	t.Run("Invalid SVG data", func(t *testing.T) {
		_, err := Decode(strings.NewReader("not an svg"))
		if err == nil {
			t.Fatal("expected error for invalid data, got nil")
		}
	})
}

func TestDecodeConfig(t *testing.T) {
	t.Run("Valid SVG config", func(t *testing.T) {
		cfg, err := DecodeConfig(strings.NewReader(sampleSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Width != 100 || cfg.Height != 100 {
			t.Errorf("expected 100x100, got %dx%d", cfg.Width, cfg.Height)
		}
	})

	t.Run("No dimensions fallback config", func(t *testing.T) {
		cfg, err := DecodeConfig(strings.NewReader(noDimensionsSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Width != 512 || cfg.Height != 512 {
			t.Errorf("expected 512x512, got %dx%d", cfg.Width, cfg.Height)
		}
	})

	t.Run("Invalid SVG config", func(t *testing.T) {
		_, err := DecodeConfig(strings.NewReader("invalid"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestImageDecodeRegistration(t *testing.T) {
	t.Run("Decode via image.Decode with <?xml", func(t *testing.T) {
		img, format, err := image.Decode(strings.NewReader(sampleSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if format != "svg" {
			t.Errorf("expected format 'svg', got %s", format)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})

	t.Run("Decode via image.Decode with <svg", func(t *testing.T) {
		img, format, err := image.Decode(strings.NewReader(minimalSVG))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if format != "svg" {
			t.Errorf("expected format 'svg', got %s", format)
		}
		if img == nil {
			t.Fatal("expected non-nil image")
		}
	})
}

