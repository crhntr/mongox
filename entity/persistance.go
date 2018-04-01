package entity

import "github.com/globalsign/mgo"

func InsertList(db *mgo.Database, entityList ...EntityReferencer) (int, error) {
	for i, entity := range entityList {
		ref := entity.Ref()
		if err := db.C(ref.Col).Insert(entity); err != nil {
			return len(entityList) - i, err
		}
	}
	return len(entityList), nil
}

func RefreshEntity(db *mgo.Database, entity EntityReferencer) error {
	ref := entity.Ref()
	return db.C(ref.Col).FindId(ref.ID).One(entity)
}

func UpdateEntity(db *mgo.Database, entity EntityReferencer, updateDoc Map) error {
	ref := entity.Ref()
	return db.C(ref.Col).UpdateId(ref.ID, updateDoc)
}

func ReadPermitted(db *mgo.Database, ref EntityReference, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.ReadPermitted(ids...)
}

func UpdatePermitted(db *mgo.Database, ref EntityReference, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.UpdatePermitted(ids...)
}

func DeletePermitted(db *mgo.Database, ref EntityReference, ids ...EntityReference) bool {
	var ent Entity
	if err := db.C(ref.Col).FindId(ref.ID).Select(SelectEntityDoc).One(&ent); err != nil {
		return false
	}
	return ent.AC.DeletePermitted(ids...)
}

func PersistClearUDR(db *mgo.Database, ref EntityReference, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll": Map{ACPath + ".u": ids, ACPath + ".d": ids, ACPath + ".r": ids},
	})
}

func PersistPermitRead(db *mgo.Database, ref EntityReference, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".u": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".r": Map{"$each": ids}},
	})
}

func PersistPermitUpdate(db *mgo.Database, ref EntityReference, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".d": ids},
		"$addToSet": Map{ACPath + ".u": Map{"$each": ids}},
	})
}

func PersistPermitDelete(db *mgo.Database, ref EntityReference, ids ...EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$pullAll":  Map{ACPath + ".r": ids, ACPath + ".u": ids},
		"$addToSet": Map{ACPath + ".d": Map{"$each": ids}},
	})
}

func PersistPublic(db *mgo.Database, ref EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": true},
	})
}

func PersistPrivate(db *mgo.Database, ref EntityReference) error {
	return db.C(ref.Col).UpdateId(ref.ID, Map{
		"$set": Map{ACPath + ".p": false},
	})
}
