package graphics

import (
	"image"
	"math"
	"testing"
)

func TestUnicodeClamp(t *testing.T) {
	u := &Unicode{}
	tests := []struct {
		name string
		v    float64
		want int
	}{
		{"negative", -10.5, 0},
		{"zero", 0, 0},
		{"positive", 128.4, 128},
		{"over_max", 260.9, 255},
		{"max", 255, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.clamp(tt.v); got != tt.want {
				t.Errorf("Unicode.clamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnicodeBrightness(t *testing.T) {
	u := &Unicode{}
	tests := []struct {
		name    string
		r, g, b uint32
		want    float64
	}{
		// Note: The method internally shifts right by 8, so it expects 16-bit color values
		{"black", 0, 0, 0, 0},
		{"white", 0xffff, 0xffff, 0xffff, 255}, // 0xffff >> 8 = 255. 0.299*255 + 0.587*255 + 0.114*255 = 255
		{"red", 0xffff, 0, 0, 0.299 * 255},
		{"green", 0, 0xffff, 0, 0.587 * 255},
		{"blue", 0, 0, 0xffff, 0.114 * 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := u.brightness(tt.r, tt.g, tt.b)
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("Unicode.brightness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnicodeDrawValid(t *testing.T) {
	u := &Unicode{}
	// Test valid RGBA image (should not panic)
	rgbaImg := image.NewRGBA(image.Rect(0, 0, 4, 4))

	err := u.Draw(rgbaImg)
	if err != nil {
		t.Errorf("Unicode.Draw() error = %v, want nil", err)
	}
}

func TestUnicodeDrawPanic(t *testing.T) {
	u := &Unicode{}
	// Test invalid image type (should panic due to lack of type assertion 'ok' check in implementation)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Unicode.Draw() was expected to panic on non-*image.RGBA")
		}
	}()

	grayImg := image.NewGray(image.Rect(0, 0, 4, 4))
	_ = u.Draw(grayImg)
}
