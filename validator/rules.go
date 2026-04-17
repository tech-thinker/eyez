package validator

import (
	"fmt"
	"strings"
)

var supportedExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"}

func Validate(filename string) error {

	for _, ext := range supportedExtensions {
		if strings.HasSuffix(filename, ext) {
			return nil
		}
	}

	err := fmt.Errorf("Unsupported file type: %s", filename)
	return err
}
