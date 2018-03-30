package mongox

import (
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/globalsign/mgo"
)

type Mux struct {
	session   *mgo.Session
	webappSrc string
}

var contentTypes = map[string]string{
	"html": "text/html; charset=utf-8",
	"js":   "application/javascript",
	"json": "application/json",
	"jpeg": "image/jpeg",
	"png":  "image/png",
	"txt":  "text/plain",
	"css":  "text/css",
}

func New(dbAddr, webappSrc string) *Mux {
	sess, err := mgo.Dial(dbAddr)
	if err != nil {
		panic(err)
	}
	return &Mux{session: sess, webappSrc: webappSrc}
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "":
		path := mux.webappSrc + "/index.html"
		log.Print("index: " + path)
		http.ServeFile(w, r, path)
	case "src":
		path := mux.webappSrc + r.URL.Path
		ext := filepath.Ext(r.URL.Path)
		w.Header().Set("Content-Type", contentTypes[ext])
		http.ServeFile(w, r, path)
	default:
		http.NotFound(w, r)
	}
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
