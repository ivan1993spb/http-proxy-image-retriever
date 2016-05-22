package main

import (
	"html/template"
	"net/http"
)

var ErrorPageTmpl = template.Must(template.New("error_page").Parse(`<!DOCTYPE html>
<html>
    <head>
        <title>Error</title>
    </head>
    <body>
        <h1>Error</h1>
        <p>{{.}}</p>
    </body>
</html>
`))

func HTTPErrorHTML(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	ErrorPageTmpl.Execute(w, error)
}