package main

import (
	"net/http"

	"github.com/crhntr/mongox"
)

func main() {
	mux := mongox.New("playground", ":27017", "src", "posts", "comments")

	http.ListenAndServe(":8080", mux)
}
