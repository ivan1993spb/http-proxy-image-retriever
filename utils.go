package main

import "regexp"

var ExpDetectDataURL = regexp.MustCompile(
	`(?i)^\s*data:([a-z]+\/[a-z0-9\-\+]+(;[a-z\-]+\=[a-z0-9\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$`)

// IsDataUrl returns true if passed string s is data url
func IsDataUrl(s string) bool {
	return ExpDetectDataURL.MatchString(s)
}

func IsBrowserImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/jp2", "image/jpx", "image/jpm", "image/webp", "image/vnd.ms-photo",
		"image/jxr", "image/gif", "image/png", "image/tiff", "image/tiff-fx", "image/svg+xml",
		"image/x‑xbitmap", "image/x‑xbm", "image/bmp", "image/x-bmp", "image/x-icon":
		return true
	}

	return false
}
