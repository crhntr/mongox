package users

import "github.com/crhntr/mongox/entity"

const UsersCollection = "users"

type User struct {
	entity.Entity `bson:",inline"`
}

func (user User) GetEntityReference() entity.EntityReference {
	return entity.EntityReference{UsersCollection, user.ID}
}
