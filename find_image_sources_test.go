package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const HTML5_SRC_TEST = `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Error</title>
    </head>
    <body>
        <img src="image.png">
        <img src="test.jpg">
        <img src="path/to/test.jpg">
        <img src="path/to/test12.jpg" />
        <img src="/path/to/test.gif">
        <img src="">
    </body>
</html>
`

func TestFindImageSourcesHTML5(t *testing.T) {
	t.Log("testing html5")
	testFindImageSources(t, HTML5_SRC_TEST)
}

const XHTML_SRC_TEST = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=windows-1251">
        <title>An XHTML 1.0 Strict standard template</title>
    </head>

    <body>
        <p>content here</p>
        <img src="image.png">
        <p>content here</p>
        <img src="test.jpg">
        <p>content here</p>
        <img src="path/to/test.jpg">
        <img src="path/to/test12.jpg" />
        <p>content here</p>
        <p>content here</p>
        <p>content here</p>
        <img src="/path/to/test.gif">
        <img src="">
    </body>
</html>
`

func TestFindImageSourcesXHTML(t *testing.T) {
	t.Log("testing xhtml")
	testFindImageSources(t, XHTML_SRC_TEST)
}

const HTML4_SRC_TEST = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
    <head>
        <title></title>
        <meta http-equiv="Content-Type" content="text/html; charset=windows-1251">
    </head>
    <body>
        <p>content here</p>
        <img src="image.png">
        <p>content here</p>
        <img src="test.jpg">
        <p>content here</p>
        <img src="path/to/test.jpg">
        <img src="path/to/test12.jpg" />
        <p>content here</p>
        <p>content here</p>
        <br/>
        <br/>
        <br/>
        <br/>
        <br/>
        <p>content here</p>
        <img src="/path/to/test.gif">
        <img src="">
    </body>
</html>
`

func TestFindImageSourcesHTML4(t *testing.T) {
	t.Log("testing html4")
	testFindImageSources(t, HTML4_SRC_TEST)
}

const HTML_SRC_TEST = `<html>
    <head>
        <title></title>
    </head>
    <body>
        <p>content here</p>
        <img src="image.png">
        <p>content here</p>
        <img src="test.jpg">
        <p>content here</p>
        <img src="path/to/test.jpg">
        <img src="path/to/test12.jpg" />
        <p>content here</p>
        <p>content here</p>
        <br/>
        <br/>
        <br/>
        <br/>
        <br/>
        <p>content here</p>
        <img src="/path/to/test.gif">
        <img src="">
    </body>
</html>
`

func TestFindImageSourcesHTML(t *testing.T) {
	t.Log("testing html")
	testFindImageSources(t, HTML_SRC_TEST)
}

func testFindImageSources(t *testing.T, html string) {
	r := strings.NewReader(html)
	stopChan := make(chan struct{})

	sources, err := findImageSources(stopChan, r)
	close(stopChan)
	assert.Nil(t, err)
	expected := []string{"image.png", "test.jpg", "path/to/test.jpg",
		"path/to/test12.jpg", "/path/to/test.gif"}
	assert.Equal(t, expected, sources, "invalid result")
}
