package mongox

import (
	"log"

	"github.com/globalsign/mgo"
)

type Mux struct {
	session   *mgo.Session
	webappSrc string
	colMap    map[string]struct{}
}

func New(dbName, dbAddr, webappSrc string, collections ...string) *Mux {
	sess, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{dbAddr},
		Database: dbAddr,
	})
	if err != nil {
		panic(err)
	}

	var colMap = make(map[string]struct{})

	for _, col := range collections {
		colMap[col] = struct{}{}
	}
	colMap["users"] = struct{}{}

	collections = collections[:0]
	for key, _ := range colMap {
		collections = append(collections, key)
	}

	log.Printf("collections: %q", collections)

	return &Mux{session: sess, webappSrc: webappSrc, colMap: colMap}
}
