// acgo
package entity

import (
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Entity struct {
	ID objectid.ObjectID `json:"_id" bson:"_id"`
	AC `json:"_ac" bson:"_ac"`
}

func New() Entity {
	return Entity{
		ID: objectid.New(),
	}
}

type EntityReference struct {
	Col string            `json:"c" bson:"c"`
	ID  objectid.ObjectID `json:"id" bson:"id"`
}

type EntityReferencer interface {
	Ref() EntityReference
}

type Map map[string]interface{}

var SelectEntityDoc = map[string]int{"_id": 1, ACPath: 1}

func (ref EntityReference) Validate() error {
	if ref.Col == "" || ref.ID == objectid.NilObjectID {
		return fmt.Errorf("invalid identity {%q: %q}", ref.Col, ref.ID)
	}
	return nil
}

func (ref EntityReference) Ref() EntityReference {
	return ref
}

func FilterEntityReferenceList(ids []EntityReference, cutset ...EntityReference) []EntityReference {
	filtered := ids[:0]
	for _, id := range cutset {
		for _, idx := range ids {
			if id.ID != idx.ID && id.Col != idx.Col {
				filtered = append(filtered, idx)
			}
		}
	}
	return filtered
}

func EntityReferences(col string, ids ...objectid.ObjectID) []EntityReference {
	refs := make([]EntityReference, len(ids))
	for i, id := range ids {
		refs[i] = EntityReference{col, id}
	}
	return refs
}
