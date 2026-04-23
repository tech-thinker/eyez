package graphics

import (
	"image"
	"math"
	"testing"
)

func TestLuminance(t *testing.T) {
	tests := []struct {
		name    string
		r, g, b uint8
		want    float64
	}{
		{"black", 0, 0, 0, 0},
		{"white", 255, 255, 255, 255},
		{"red", 255, 0, 0, 0.2126 * 255},
		{"green", 0, 255, 0, 0.7152 * 255},
		{"blue", 0, 0, 255, 0.0722 * 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := luminance(tt.r, tt.g, tt.b)
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("luminance(%d, %d, %d) = %f; want %f", tt.r, tt.g, tt.b, got, tt.want)
			}
		})
	}
}

func TestGammaCorrect(t *testing.T) {
	tests := []struct {
		name string
		v    uint8
		want uint8
	}{
		{"min", 0, 0},
		{"max", 255, 255},
		// For 128: 128/255 = 0.50196
		// math.Pow(0.50196, 1/2.2) = 0.7297
		// 0.7297 * 255 = 186.08 -> 186
		{"mid", 128, 186},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gammaCorrect(tt.v)
			// allow a small margin if floating point rounding causes off-by-one
			// but math.Pow algorithm yields exactly 186
			if got != tt.want {
				t.Errorf("gammaCorrect(%d) = %d; want %d", tt.v, got, tt.want)
			}
		})
	}
}

func TestASCIIDraw(t *testing.T) {
	a := &ASCII{}

	t.Run("Unsupported Image Type", func(t *testing.T) {
		grayImg := image.NewGray(image.Rect(0, 0, 10, 10))
		err := a.Draw(grayImg)
		if err == nil {
			t.Fatal("expected error for unsupported image type, got nil")
		}
		expectedErr := "expected *image.NRGBA"
		if err.Error() != expectedErr {
			t.Errorf("expected error '%s', got '%v'", expectedErr, err)
		}
	})

	t.Run("Supported Image Type", func(t *testing.T) {
		// Draw() accepts *image.RGBA despite the error message saying *image.NRGBA
		rgbaImg := image.NewRGBA(image.Rect(0, 0, 10, 10))
		err := a.Draw(rgbaImg)
		if err != nil {
			t.Fatalf("expected nil error for valid image, got %v", err)
		}
	})
}
