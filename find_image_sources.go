package main

import (
	"io"

	"golang.org/x/net/html"
)

func findImageSources(stopChan <-chan struct{}, r io.Reader) ([]string, error) {
	z := html.NewTokenizer(r)
	sources := make([]string, 0)

	for {
		tokenType := z.Next()
		switch tokenType {
		case html.ErrorToken:
			if err := z.Err(); err != nil && err != io.EOF {
				return nil, err
			}
			return sources, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := z.Token()

			if token.Data == "img" && len(token.Attr) > 0 {
				for _, attr := range token.Attr {
					if attr.Key == "src" && attr.Val != "" {
						sources = append(sources, attr.Val)
						break
					}
				}
			}
		}
	}
	return nil, nil
}
