package main

import (
	"net/http"

	"github.com/crhntr/mongox"
)

func main() {
	mux := mongox.New(":27017", "src")

	http.ListenAndServe(":8080", mux)
}
