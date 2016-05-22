package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func FindImageSources(r io.Reader) (sources []string, err error) {
	z := html.NewTokenizer(r)
	sources = make([]string, 0)

	for {
		switch tokenType := z.Next(); tokenType {
		case html.ErrorToken:
			if e := z.Err(); e != nil && e != io.EOF {
				err = e
			}
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			if token := z.Token(); strings.ToLower(token.Data) == "img" && len(token.Attr) > 0 {
				for _, attr := range token.Attr {
					if strings.ToLower(attr.Key) == "src" && attr.Val != "" {
						sources = append(sources, attr.Val)
						break
					}
				}
			}
		}
	}

	return
}
