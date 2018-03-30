// acgo
package entity

import (
	"errors"
	"fmt"

	"github.com/crhntr/mongo-go-driver/bson/objectid"
	"github.com/globalsign/mgo"
)

type Map map[string]interface{}

type Entity struct {
	ID objectid.ObjectID `json:"_id" bson:"_id"`
	AC `json:"_ac" bson:"_ac"`
}

func New() Entity {
	return Entity{
		ID: objectid.New(),
	}
}

var ACPath = "_ac"

var SelectEntityDoc = map[string]int{"_id": 1, ACPath: 1}

type EntityReference struct {
	Col string            `json:"c" bson:"c"`
	ID  objectid.ObjectID `json:"id" bson:"id"`
}

func (ref EntityReference) Validate() error {
	if ref.Col == "" || ref.ID == objectid.NilObjectID {
		return fmt.Errorf("invalid identity {%q: %q}", ref.Col, ref.ID)
	}
	return nil
}

func (ref EntityReference) GetEntityReference() EntityReference {
	return ref
}

type EntityReferencer interface {
	GetEntityReference() EntityReference
}

// AC should be Embeded in structs to be stored in MongoDB
// It should be anotated with the `bson:"ac"` or whatever ACPath is set to.
// When a new object is created, the creator's identity should be passed to SetCreator
// bson tag "inline" should not be set
type AC struct {
	Readers  []EntityReference `json:"r,omitempty" bson:"r,omitempty"`
	Updaters []EntityReference `json:"u,omitempty" bson:"u,omitempty"`
	Deleters []EntityReference `json:"d,omitempty" bson:"d,omitempty"`
	Creator  *EntityReference  `json:"cr,omitempty" bson:"cr,omitempty"`
	Public   bool              `json:"p" bson:"p"`
}

func (ac *AC) SetCreator(id EntityReference) error {
	if ac.Creator != nil {
		return errors.New("creator already set")
	}
	if err := id.Validate(); err != nil {
		return err
	}
	ac.Creator = &id
	return nil
}

func (ac AC) ReadPermitted(ids ...EntityReference) bool {
	if ac.Public {
		return true
	}
	for _, id := range ids {
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Updaters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Readers {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (ac AC) UpdatePermitted(ids ...EntityReference) bool {
	for _, id := range ids {
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
		for _, idInSet := range ac.Updaters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (ac AC) DeletePermitted(ids ...EntityReference) bool {
	for _, id := range ids {
		if ac.Creator != nil && ac.Creator.Col == id.Col && ac.Creator.ID == id.ID {
			return true
		}

		for _, idInSet := range ac.Deleters {
			if idInSet.Col == id.Col && idInSet.ID == id.ID {
				return true
			}
		}
	}
	return false
}

func (id EntityReference) ReadPermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(id.Col).FindId(id.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.ReadPermitted(ids...)
}

func (id EntityReference) UpdatePermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(id.Col).FindId(id.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.UpdatePermitted(ids...)
}

func (id EntityReference) DeletePermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(id.Col).FindId(id.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.DeletePermitted(ids...)
}

func (ac *AC) ClearUDR(ids ...EntityReference) {
	for _, id := range ids {
		ac.Updaters = FilterEntityReferenceList(ac.Updaters, id)
		ac.Deleters = FilterEntityReferenceList(ac.Deleters, id)
		ac.Readers = FilterEntityReferenceList(ac.Readers, id)
	}
}

func (ac *AC) PermitRead(ids ...EntityReference) {
	for _, id := range ids {
		ac.Updaters = FilterEntityReferenceList(ac.Updaters, id)
		ac.Deleters = FilterEntityReferenceList(ac.Deleters, id)
		ac.Readers = FilterEntityReferenceList(ac.Readers, id)

		ac.Readers = append(ac.Readers, id)
	}
}

func (ac *AC) PermitUpdate(ids ...EntityReference) {
	for _, id := range ids {
		ac.Updaters = FilterEntityReferenceList(ac.Updaters, id)
		ac.Deleters = FilterEntityReferenceList(ac.Deleters, id)
		ac.Readers = FilterEntityReferenceList(ac.Readers, id)

		ac.Updaters = append(ac.Updaters, id)
	}
}

func (ac *AC) PermitDelete(ids ...EntityReference) {
	for _, id := range ids {
		ac.Updaters = FilterEntityReferenceList(ac.Updaters, id)
		ac.Deleters = FilterEntityReferenceList(ac.Deleters, id)
		ac.Readers = FilterEntityReferenceList(ac.Readers, id)

		ac.Deleters = append(ac.Deleters, id)
	}
}

func (entity EntityReference) PersistClearUDR(db *mgo.Database, ids ...EntityReference) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$pullAll": Map{ACPath + ".u": ids, ACPath + ".d": ids, ACPath + ".r": ids},
	})
}

func (entity EntityReference) PersistPermitRead(db *mgo.Database, ids ...EntityReference) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$pullAll":  Map{ACPath + ".u": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".r": Map{"$each": ids}},
	})
}

func (entity EntityReference) PersistPermitUpdate(db *mgo.Database, ids ...EntityReference) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".u": Map{"$each": ids}},
	})
}

func (entity EntityReference) PersistPermitDelete(db *mgo.Database, ids ...EntityReference) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".u": ids},
		"$addToSet": Map{ACPath + ".d": Map{"$each": ids}},
	})
}

func (entity EntityReference) PersistPublic(db *mgo.Database) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$set": Map{ACPath + ".p": true},
	})
}

func (entity EntityReference) PersistPrivate(db *mgo.Database) error {
	return db.C(entity.Col).UpdateId(entity.ID, Map{
		"$set": Map{ACPath + ".p": false},
	})
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

func InsertList(db *mgo.Database, entityList ...EntityReferencer) (int, error) {
	for i, entity := range entityList {
		ref := entity.GetEntityReference()
		if err := db.C(ref.Col).Insert(entity); err != nil {
			return len(entityList) - i, err
		}
	}
	return len(entityList), nil
}

func RefreshEntity(db *mgo.Database, entity EntityReferencer) error {
	ref := entity.GetEntityReference()
	return db.C(ref.Col).FindId(ref.ID).One(entity)
}

func UpdateEntity(db *mgo.Database, entity EntityReferencer, updateDoc Map) error {
	ref := entity.GetEntityReference()
	return db.C(ref.Col).UpdateId(ref.ID, updateDoc)
}
