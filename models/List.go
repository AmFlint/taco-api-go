package models

import "gopkg.in/mgo.v2/bson"

type List struct {
	ListId bson.ObjectId `bson:"_id" json:"listId"`
	Name   string        `bson:"name" json:"name" onCreate:"nonzero,max=30,regexp=^[a-zA-Z-_ ]*$"`
	Order  int           `bson:"order" json:"order"`
	Tasks  []Task        `bson:"tasks" json:"tasks"`
}

// Initialize List structure with empty array of task
func NewList() List {
	list := List{}
	list.Tasks = []Task{}
	return list
}

// Hydrate a List structure from a map of string -> interface
func (l *List) HydrateFromMap(json map[string]interface{}) {
	if name, ok := json["name"]; ok {
		l.Name = name.(string)
	}
}
