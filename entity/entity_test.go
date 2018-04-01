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
		if err := idfr.GetEntityReference().Validate(); err != nil {
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

	if err := post0.GetEntityReference().PersistClearUDR(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}

	if err := post0.GetEntityReference().PersistPermitRead(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if err := post0.GetEntityReference().PersistPermitUpdate(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if err := post0.GetEntityReference().PersistPermitDelete(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if !post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}

	if err := post0.GetEntityReference().PersistPermitUpdate(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}

	if err := post0.GetEntityReference().PersistPermitRead(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}

	if err := post0.GetEntityReference().PersistClearUDR(db, user1.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if post0.GetEntityReference().DeletePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().UpdatePermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.GetEntityReference().ReadPermitted(db, user1.GetEntityReference()) {
		t.Fatal()
	}

	if err := user2.GetEntityReference().PersistClearUDR(db); err == nil {
		t.Fatal()
	}
	if err := user2.GetEntityReference().PersistPermitRead(db); err == nil {
		t.Fatal()
	}
	if err := user2.GetEntityReference().PersistPermitUpdate(db); err == nil {
		t.Fatal()
	}
	if err := user2.GetEntityReference().PersistPermitDelete(db); err == nil {
		t.Fatal()
	}

	if user2.GetEntityReference().ReadPermitted(db) {
		t.Fatal()
	}
	if user2.GetEntityReference().UpdatePermitted(db) {
		t.Fatal()
	}
	if user2.GetEntityReference().DeletePermitted(db) {
		t.Fatal()
	}

	if post0.GetEntityReference().ReadPermitted(db, user2.GetEntityReference()) {
		t.Fatal()
	}
	user2.Teams = append(user2.Teams, team0.ID)

	post0.GetEntityReference().PersistPermitRead(db, team0.GetEntityReference())

	refs := append(entity.EntityReferences(TeamCol, user2.Teams...), user2.GetEntityReference())

	if !post0.GetEntityReference().ReadPermitted(db, refs...) {
		t.Fatal()
	}

	entity.UpdateEntity(db, user0, entity.Map{
		"$set": mp{"foo": "bar"},
	})
	post0.GetEntityReference().GetEntityReference()

	post3 := Post{Entity: entity.New()}
	entity.InsertList(db, post3)
	post3.GetEntityReference().PersistPublic(db)

	if !post3.GetEntityReference().ReadPermitted(db, user0.GetEntityReference()) {
		t.Fatal()
	}
	post3.GetEntityReference().PersistPrivate(db)
	if post3.GetEntityReference().ReadPermitted(db, user0.GetEntityReference()) {
		t.Fatal()
	}
}
