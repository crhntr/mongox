package mongox

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	switch head {
	case "api":
		mux.handleAPI(w, r)
	case "":
		path := mux.webappSrc + "/index.html"
		http.ServeFile(w, r, path)
	case "dist-src":
		path := os.Getenv("GOPATH") + "/src/github.com/crhntr/mongox/src/" + r.URL.Path
		w.Header().Set("Content-Type", contentTypes[filepath.Ext(r.URL.Path)])
		http.ServeFile(w, r, path)
	case "src":
		path := mux.webappSrc + r.URL.Path
		w.Header().Set("Content-Type", contentTypes[filepath.Ext(r.URL.Path)])
		http.ServeFile(w, r, path)
	default:
		http.NotFound(w, r)
	}
}

func (mux *Mux) handleAPI(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	var cols []string

	for key, _ := range mux.colMap {
		cols = append(cols, key)
	}

	w.Header().Set("Content-Type", contentTypes["json"])
	enc := json.NewEncoder(w)
	switch head {
	case "db":
		w.WriteHeader(http.StatusOK)
		enc.Encode(struct {
			Collections []string `json:"collections"`
		}{cols})
	}
}

// shiftPath splits off the first component of p, which will be cleaned of
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
