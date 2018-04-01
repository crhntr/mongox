package entity_test

import (
	"github.com/crhntr/mongox/entity"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const UserCol = "user"

type User struct {
	entity.Entity `bson:",inline"`
	Teams         []objectid.ObjectID `bson:"teams"`
}

func (this User) Ref() entity.EntityReference {
	return entity.EntityReference{UserCol, this.ID}
}

const TeamCol = "team"

type Team struct {
	entity.Entity `bson:",inline"`
}

func (this Team) Ref() entity.EntityReference {
	return entity.EntityReference{TeamCol, this.ID}
}

const PostCol = "post"

type Post struct {
	entity.Entity `bson:",inline"`
	N             int `bson:"n"`
}

func (this Post) Ref() entity.EntityReference {
	return entity.EntityReference{PostCol, this.ID}
}
