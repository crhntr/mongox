package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/crhntr/mongox"
	"github.com/crhntr/mongox/entity"
)

type Comment struct {
	entity.Entity `bson:",inline"`

	Text string `json:"text" bson:"text"`
}

type Post struct {
	entity.Entity `bson:",inline"`
}

func main() {
	mux := mongox.New("playground", ":27017", "src")

	mux.Resource("comments", mongox.ResourceClosures{
		Insert: func(r io.Reader) (interface{}, error) {
			var comment Comment
			err := json.NewDecoder(r).Decode(&comment)
			log.Println(comment)
			comment.Entity = entity.New()
			return comment, err
		},
	})

	http.ListenAndServe(":8080", mux)
}
