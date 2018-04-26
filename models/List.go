package models

import "gopkg.in/mgo.v2/bson"

type List struct {
	ListId bson.ObjectId `bson:"_id" json:"listId"`
	Name   string        `bson:"name" json:"name" onCreate:"nonzero,max=30,regexp=^[a-zA-Z-_ ]*$"`
	Order  int         `bson:"order" json:"order"`
}
