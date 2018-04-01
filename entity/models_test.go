package entity_test

import (
	"github.com/crhntr/mongox/entity"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const (
	UserCol = "user"
	TeamCol = "team"
	PostCol = "post"
)

type (
	User struct {
		entity.Entity `bson:",inline"`
		Teams         []objectid.ObjectID `bson:"teams"`
	}

	Team struct {
		entity.Entity `bson:",inline"`
	}

	Post struct {
		entity.Entity `bson:",inline"`
		N             int `bson:"n"`
	}
)

func (this User) Ref() entity.EntityReference {
	return entity.EntityReference{UserCol, this.ID}
}

func (this Team) Ref() entity.EntityReference {
	return entity.EntityReference{TeamCol, this.ID}
}

func (this Post) Ref() entity.EntityReference {
	return entity.EntityReference{PostCol, this.ID}
}
