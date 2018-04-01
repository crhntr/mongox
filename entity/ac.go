package entity

import "errors"

var ACPath = "_ac"

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
