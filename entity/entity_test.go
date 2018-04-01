package entity_test

import (
	"flag"
	"os"
	"testing"

	"github.com/crhntr/mongox/entity"
	"github.com/globalsign/mgo"
)

var (
	dbName             = "mongox_entity_test"
	databaseSession, _ = mgo.DialWithInfo(&mgo.DialInfo{
		Database: dbName,
		Addrs:    []string{":27017"},
	})
)

type mp map[string]interface{}

func TestMain(m *testing.M) {
	flag.Parse()
	defer databaseSession.Close()
	defer databaseSession.DB("").DropDatabase()
	os.Exit(m.Run())
}

func TestEntity(t *testing.T) {
	db := databaseSession.DB("")

	identityZeroVal := entity.EntityReference{}
	if identityZeroVal.Validate() == nil {
		t.Fatal()
	}

	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}
	user2 := User{Entity: entity.New()}

	team0 := Team{Entity: entity.New()}
	team1 := Team{Entity: entity.New()}

	for _, idfr := range []entity.EntityReferencer{user0, user1, user2, team0, team1} {
		if err := idfr.Ref().Validate(); err != nil {
			t.Fatal(err)
		}
	}

	post0 := Post{Entity: entity.New()}

	if _, err := entity.InsertList(db, user0, user1, team0, post0); err != nil {
		t.Fatal(err)
	}
	if _, err := entity.InsertList(db, user0); err == nil {
		t.Fatal(err)
	}

	post0N := 101
	if err := db.C(PostCol).UpdateId(post0.ID, mp{"$set": mp{"n": post0N}}); err != nil {
		t.Fatal(err)
	}
	// t.Log(post0)
	entity.RefreshEntity(db, &post0)
	// t.Log(post0)
	if post0.N != post0N {
		t.Fatal()
	}

	if err := entity.PersistClearUDR(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	if err := entity.PersistPermitRead(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if err := entity.PersistPermitUpdate(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if err := entity.PersistPermitDelete(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if !entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	if err := entity.PersistPermitUpdate(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	if err := entity.PersistPermitRead(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if !entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	if err := entity.PersistClearUDR(db, post0.Ref(), user1.Ref()); err != nil {
		t.Fatal()
	}
	if entity.DeletePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.UpdatePermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}
	if entity.ReadPermitted(db, post0.Ref(), user1.Ref()) {
		t.Fatal()
	}

	if err := entity.PersistClearUDR(db, user2.Ref()); err == nil {
		t.Fatal()
	}
	if err := entity.PersistPermitRead(db, user2.Ref()); err == nil {
		t.Fatal()
	}
	if err := entity.PersistPermitUpdate(db, user2.Ref()); err == nil {
		t.Fatal()
	}
	if err := entity.PersistPermitDelete(db, user2.Ref()); err == nil {
		t.Fatal()
	}

	if entity.ReadPermitted(db, user2.Ref()) {
		t.Fatal()
	}
	if entity.UpdatePermitted(db, user2.Ref()) {
		t.Fatal()
	}
	if entity.DeletePermitted(db, user2.Ref()) {
		t.Fatal()
	}

	if entity.ReadPermitted(db, post0.Ref(), user2.Ref()) {
		t.Fatal()
	}
	user2.Teams = append(user2.Teams, team0.ID)

	entity.PersistPermitRead(db, post0.Ref(), team0.Ref())

	refs := append(entity.EntityReferences(TeamCol, user2.Teams...), user2.Ref())

	if !entity.ReadPermitted(db, post0.Ref(), refs...) {
		t.Fatal()
	}

	entity.UpdateEntity(db, user0, entity.Map{
		"$set": mp{"foo": "bar"},
	})
	post0.Ref().Ref()

	post3 := Post{Entity: entity.New()}
	entity.InsertList(db, post3)
	entity.PersistPublic(db, post3.Ref())

	if !entity.ReadPermitted(db, post3.Ref(), user0.Ref()) {
		t.Fatal()
	}
	entity.PersistPrivate(db, post3.Ref())
	if entity.ReadPermitted(db, post3.Ref(), user0.Ref()) {
		t.Fatal()
	}
}
