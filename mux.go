package mongox

import "github.com/globalsign/mgo"

type Mux struct {
	session   *mgo.Session
	webappSrc string
}

func New(dbAddr, webappSrc string) *Mux {
	sess, err := mgo.Dial(dbAddr)
	if err != nil {
		panic(err)
	}
	return &Mux{session: sess, webappSrc: webappSrc}
}
