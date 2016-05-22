package main

import (
	"html/template"
	"net/http"
)

var ErrorPageTmpl = template.Must(template.New("error_page").Parse(`<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title>Error</title>
    </head>
    <body>
        <h1>Error</h1>
        <p>{{.}}</p>
    </body>
</html>
`))

// HTTPErrorHTML sends error message with specific status code
func HTTPErrorHTML(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	ErrorPageTmpl.Execute(w, error)
}
