package main

import (
	"log"
	"regexp"
)

var ExpDetectDataURL = regexp.MustCompile(`(?i)^\s*data:([a-z]+\/[a-z0-9\-\+]+(;[a-z\-]+\=[a-z0-9\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$`)

// IsDataUrl returns true if passed string s is data url
func IsDataUrl(s string) bool {
	return ExpDetectDataURL.MatchString(s)
}

func LogErrorChan(stopChan <-chan struct{}, errChan <-chan error, logger *log.Logger) {
	for {
		select {
		case <-stopChan:
			return
		case err := <-errChan:
			logger.Println(err)
		}
	}
}
