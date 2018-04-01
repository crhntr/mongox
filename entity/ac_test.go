package entity_test

import (
	"testing"

	"github.com/crhntr/mongox/entity"
)

func TestACL(t *testing.T) {
	user0 := User{Entity: entity.New()}
	user1 := User{Entity: entity.New()}
	// user2 := User{Entity: entity.New()}

	team0 := Team{Entity: entity.New()}
	team1 := Team{Entity: entity.New()}
	post0 := Post{Entity: entity.New()}

	identityZeroVal := entity.EntityReference{}
	if err := post0.SetCreator(identityZeroVal); err == nil {
		t.Fatal()
	}
	if err := post0.SetCreator(user0.GetEntityReference()); err != nil {
		t.Fatal()
	}
	if err := post0.SetCreator(user0.GetEntityReference()); err == nil {
		t.Fatal()
	}

	if !post0.DeletePermitted(user0.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user0.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user0.GetEntityReference()) {
		t.Fatal()
	}

	if post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	if !post0.ReadPermitted(user1.GetEntityReference(), user0.GetEntityReference()) {
		t.Fatal()
	}

	if !post0.ReadPermitted(user0.GetEntityReference(), user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitDelete(team0.GetEntityReference(), team1.GetEntityReference())

	post0.ClearUDR(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitRead(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitUpdate(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitDelete(user1.GetEntityReference())
	if !post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitUpdate(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.PermitRead(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if !post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}

	post0.ClearUDR(user1.GetEntityReference())
	if post0.DeletePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.UpdatePermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
	if post0.ReadPermitted(user1.GetEntityReference()) {
		t.Fatal()
	}
}
