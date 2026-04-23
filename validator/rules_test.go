package validator

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid png",
			filename: "image.png",
			wantErr:  false,
		},
		{
			name:     "valid jpg",
			filename: "photo.jpg",
			wantErr:  false,
		},
		{
			name:     "valid jpeg",
			filename: "picture.jpeg",
			wantErr:  false,
		},
		{
			name:     "valid gif",
			filename: "animation.gif",
			wantErr:  false,
		},
		{
			name:     "valid bmp",
			filename: "bitmap.bmp",
			wantErr:  false,
		},
		{
			name:     "valid webp",
			filename: "image.webp",
			wantErr:  false,
		},
		{
			name:     "valid multiple dots supported extension",
			filename: "image.backup.png",
			wantErr:  false,
		},
		{
			name:     "invalid multiple dots unsupported extension",
			filename: "image.backup.txt",
			wantErr:  true,
		},
		{
			name:     "invalid extension txt",
			filename: "document.txt",
			wantErr:  true,
		},
		{
			name:     "invalid extension pdf",
			filename: "document.pdf",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
		{
			name:     "no extension",
			filename: "image",
			wantErr:  true,
		},
		{
			name:     "hidden file no extension",
			filename: ".hidden",
			wantErr:  true,
		},
		{
			name:     "just extension",
			filename: ".png",
			wantErr:  false,
		},
		{
			name:     "supported uppercase extension",
			filename: "image.PNG",
			wantErr:  false,
		},
		{
			name:     "uppercase filename and extension",
			filename: "PICTURE.JPG",
			wantErr:  false,
		},
		{
			name:     "very long filename supported extension",
			filename: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.png",
			wantErr:  false,
		},
		{
			name:     "very long filename unsupported extension",
			filename: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb.txt",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
