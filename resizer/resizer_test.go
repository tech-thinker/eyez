package resizer

import (
	"image"
	"image/color"
	"testing"
)

// createDummyImage generates a simple flat-colored image with the given dimensions.
func createDummyImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color just to make it a valid non-empty image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	return img
}

func TestResizerImpl(t *testing.T) {
	tests := []struct {
		name         string
		resizer      Resizer
		srcWidth     int
		srcHeight    int
		targetWidth  int
		expectWidth  int
		expectHeight int
	}{
		{
			name:         "CatmullRom ScaleDown Square",
			resizer:      &CatmullRom{},
			srcWidth:     100,
			srcHeight:    100,
			targetWidth:  50,
			expectWidth:  50,
			expectHeight: 50,
		},
		{
			name:         "CatmullRom ScaleDown Portrait",
			resizer:      &CatmullRom{},
			srcWidth:     100,
			srcHeight:    200,
			targetWidth:  50,
			expectWidth:  50,
			expectHeight: 100,
		},
		{
			name:         "CatmullRom ScaleUp Landscape",
			resizer:      &CatmullRom{},
			srcWidth:     200,
			srcHeight:    100,
			targetWidth:  400,
			expectWidth:  400,
			expectHeight: 200,
		},
		{
			name:         "Lanczos ScaleDown Square",
			resizer:      &Lanczos{},
			srcWidth:     100,
			srcHeight:    100,
			targetWidth:  50,
			expectWidth:  50,
			expectHeight: 50,
		},
		{
			name:         "Lanczos ScaleDown Portrait",
			resizer:      &Lanczos{},
			srcWidth:     100,
			srcHeight:    200,
			targetWidth:  50,
			expectWidth:  50,
			expectHeight: 100,
		},
		{
			name:         "Lanczos ScaleUp Landscape",
			resizer:      &Lanczos{},
			srcWidth:     200,
			srcHeight:    100,
			targetWidth:  400,
			expectWidth:  400,
			expectHeight: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := createDummyImage(tt.srcWidth, tt.srcHeight)

			dst := tt.resizer.Resize(src, tt.targetWidth)

			if dst == nil {
				t.Fatalf("Resize returned nil")
			}

			bounds := dst.Bounds()
			if bounds.Dx() != tt.expectWidth {
				t.Errorf("Expected width %d, got %d", tt.expectWidth, bounds.Dx())
			}
			if bounds.Dy() != tt.expectHeight {
				t.Errorf("Expected height %d, got %d", tt.expectHeight, bounds.Dy())
			}
		})
	}
}
