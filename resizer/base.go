package resizer

import "image"

type Resizer interface {
	Resize(src image.Image, newW int) image.Image
}
