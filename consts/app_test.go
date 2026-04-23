package consts

import "testing"

// This test is just to avoid unintentional changes to constant values.
func TestConstantsValues(t *testing.T) {
	// Defaults
	if DEFAULT_WIDTH != 80 {
		t.Errorf("Expected DEFAULT_WIDTH to be 80, got %v", DEFAULT_WIDTH)
	}
	if DEFAULT_GRAPHICS != "unicode" {
		t.Errorf("Expected DEFAULT_GRAPHICS to be 'unicode', got %v", DEFAULT_GRAPHICS)
	}
	if DEFAULT_ALGORITHM != "catmull-rom" {
		t.Errorf("Expected DEFAULT_ALGORITHM to be 'catmull-rom', got %v", DEFAULT_ALGORITHM)
	}

	// Output modes
	if GRAPHICS_UNICODE != "unicode" {
		t.Errorf("Expected GRAPHICS_UNICODE to be 'unicode', got %v", GRAPHICS_UNICODE)
	}
	if GRAPHICS_ASCII != "ascii" {
		t.Errorf("Expected GRAPHICS_ASCII to be 'ascii', got %v", GRAPHICS_ASCII)
	}
	if GRAPHICS_KITTY != "kitty" {
		t.Errorf("Expected GRAPHICS_KITTY to be 'kitty', got %v", GRAPHICS_KITTY)
	}

	// Algorithms
	if ALGO_CATMULL_ROM != "catmull-rom" {
		t.Errorf("Expected ALGO_CATMULL_ROM to be 'catmull-rom', got %v", ALGO_CATMULL_ROM)
	}
	if ALGO_LANCZOS != "lanczos" {
		t.Errorf("Expected ALGO_LANCZOS to be 'lanczos', got %v", ALGO_LANCZOS)
	}
}
