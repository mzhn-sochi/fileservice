package utils

import (
	"strings"
)

const (
	png = "PNG"
	jpg = "JPG"
)

func GetContentTypeFromB64(b64 string) (string, bool) {
	if strings.Contains(b64, png) {
		return "image/png", true
	} else if strings.Contains(b64, jpg) {
		return "image/jpg", true
	}

	return "", false
}
