package mongox

import (
	"github.com/globalsign/mgo"
)

type Mux struct {
	session   *mgo.Session
	webappSrc string
	colMap    map[string]ResourceHandlers
}

func New(dbName, dbAddr, webappSrc string) *Mux {
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{dbAddr},
		Database: dbName,
	})
	if err != nil {
		panic(err)
	}
	sess.SetSafe(&mgo.Safe{})

	return &Mux{
		session:   sess,
		webappSrc: webappSrc,
		colMap:    make(map[string]ResourceHandlers),
	}
}

func (mux *Mux) Resource(col string, closures ResourceHandlers) {
	if _, found := mux.colMap[col]; found {
		panic("resource with name " + col + "already exists")
	}
	mux.colMap[col] = closures
}
