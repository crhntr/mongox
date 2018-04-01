package mongox

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/globalsign/mgo"
)

type ResourceHandlers struct {
	Insert func(w http.ResponseWriter, r *http.Request) (interface{}, error)
}

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

	if r.URL.Path != "/" {
		col, _ := shiftPath(r.URL.Path)

		switch r.Method {
		case http.MethodPost:
			// mux.insert(w, r, col)
		case http.MethodGet:
			mux.find(w, r, col)
		case http.MethodPatch:
			mux.update(w, r, col)
		case http.MethodDelete:
			mux.remove(w, r, col)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch head {
		case "db":
			w.WriteHeader(http.StatusOK)

			w.Header().Set("Content-Type", contentTypes["json"])

			var cols []string

			for key, _ := range mux.colMap {
				cols = append(cols, key)
			}

			json.NewEncoder(w).Encode(struct {
				Collections []string `json:"collections"`
			}{cols})
		default:
			http.NotFound(w, r)
		}
	}
}

func (mux *Mux) insert(w http.ResponseWriter, r *http.Request, col string) {
	resource, err := mux.colMap[col].Insert(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	sess := mux.session.Clone()
	defer sess.Close()

	if err := sess.DB("").C(col).Insert(resource); err != nil {
		if mgo.IsDup(err) {
			http.Error(w, "duplicate resource", http.StatusBadRequest)
			return
		}
		if err.Error() == "Document failed validation" {
			http.Error(w, "Document failed validation", http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func (mux *Mux) find(w http.ResponseWriter, r *http.Request, col string) {
	http.Error(w, "error find not implemented", http.StatusNotImplemented)
}
func (mux *Mux) update(w http.ResponseWriter, r *http.Request, col string) {
	http.Error(w, "error update not implemented", http.StatusNotImplemented)
}
func (mux *Mux) remove(w http.ResponseWriter, r *http.Request, col string) {
	http.Error(w, "error remove not implemented", http.StatusNotImplemented)
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
