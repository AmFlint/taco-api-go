package models

import "gopkg.in/mgo.v2/bson"

// Task Structure, represents Task Document from Database
type Task struct {
	TaskId      bson.ObjectId `bson:"task_id" json:"task_id"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Status      string `bson:"status" json:"status"`
}

