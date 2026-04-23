package validator

import (
	"fmt"
	"strings"
)

var supportedExtensions = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"}

func Validate(filename string) error {

	lowerFilename := strings.ToLower(filename)
	for _, ext := range supportedExtensions {
		if strings.HasSuffix(lowerFilename, ext) {
			return nil
		}
	}

	err := fmt.Errorf("Unsupported file type: %s", filename)
	return err
}
