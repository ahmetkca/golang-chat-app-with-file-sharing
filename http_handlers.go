package main

import (
	_ "embed"
	"net/http"
)

// embeding index.html file into go executable binary
/*--------------------*/

//go:embed index.html
var indexHTML []byte

/*--------------------*/

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html") // setting response type to be html file
	w.Write(indexHTML)                          // writing the embeded index.html file back to the user
}
