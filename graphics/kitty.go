package graphics

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

type Kitty struct {
}

func (k *Kitty) Draw(img image.Image) error {
	if os.Getenv("KITTY_WINDOW_ID") == "" && os.Getenv("TERM") != "xterm-kitty" {
		return fmt.Errorf("kitty graphics not supported or not in a kitty terminal")
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	const chunkSize = 4096
	w := bufio.NewWriter(os.Stdout)

	imageID := 1
	offset := 0
	first := true

	for offset < len(encoded) {
		end := min(offset+chunkSize, len(encoded))
		chunk := encoded[offset:end]

		more := 1
		if end == len(encoded) {
			more = 0
		}

		if first {
			_, err := fmt.Fprintf(
				w,
				"\x1b_Ga=T,f=100,i=%d,m=%d,q=2;%s\x1b\\",
				imageID,
				more,
				chunk,
			)
			if err != nil {
				return err
			}
			first = false
		} else {
			_, err := fmt.Fprintf(
				w,
				"\x1b_Gm=%d;%s\x1b\\",
				more,
				chunk,
			)
			if err != nil {
				return err
			}
		}

		offset = end
	}
	return w.Flush()
}
