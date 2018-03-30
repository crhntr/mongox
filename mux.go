package mongox

import (
	"github.com/crhntr/mongox/entity"
	"github.com/globalsign/mgo"
)

type MakeResourceFunc func() entity.EntityReferencer

type Mux struct {
	session   *mgo.Session
	webappSrc string
	colMap    map[string]ResourceClosures
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
		colMap:    make(map[string]ResourceClosures),
	}
}

func (mux *Mux) Resource(col string, closures ResourceClosures) {
	if _, found := mux.colMap[col]; found {
		panic("resource with name " + col + "already exists")
	}
	mux.colMap[col] = closures
}
