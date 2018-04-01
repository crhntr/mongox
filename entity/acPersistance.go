package entity

import "github.com/globalsign/mgo"

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

func (ref EntityReference) ReadPermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.ReadPermitted(ids...)
}

func (ref EntityReference) UpdatePermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.UpdatePermitted(ids...)
}

func (ref EntityReference) DeletePermitted(db *mgo.Database, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.DeletePermitted(ids...)
}

func (ref EntityReference) PersistClearUDR(db *mgo.Database, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll": Map{ACPath + ".u": ids, ACPath + ".d": ids, ACPath + ".r": ids},
	})
}

func (ref EntityReference) PersistPermitRead(db *mgo.Database, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".u": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".r": Map{"$each": ids}},
	})
}

func (ref EntityReference) PersistPermitUpdate(db *mgo.Database, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".u": Map{"$each": ids}},
	})
}

func (ref EntityReference) PersistPermitDelete(db *mgo.Database, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".u": ids},
		"$addToSet": Map{ACPath + ".d": Map{"$each": ids}},
	})
}

func (ref EntityReference) PersistPublic(db *mgo.Database) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": true},
	})
}

func (ref EntityReference) PersistPrivate(db *mgo.Database) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": false},
	})
}
