package main

import (
	"regexp"
)

var ExpDetectDataURL = regexp.MustCompile(`(?i)^\s*data:([a-z]+\/[a-z0-9\-\+]+(;[a-z\-]+\=[a-z0-9\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$`)

func IsDataUrl(s string) bool {
	return ExpDetectDataURL.MatchString(s)
}
