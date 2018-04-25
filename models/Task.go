package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Task Structure, represents Task Document from Database
type Task struct {
	TaskId      bson.ObjectId `bson:"_id" json:"taskId"`
	Title       string        `bson:"title" json:"title" onCreate:"nonzero,max=200"`
	Description string        `bson:"description" json:"description" onCreate:"max=500"`
	Status      bool          `bson:"status" json:"status"`
	Points      float64       `bson:"points" json:"points" onCreate:"min=0,max=100"`
}

// Set Default Status to a Task Entity
func (t *Task) SetDefaultStatus() {
	t.Status = false
}

// Hydrate a Task structure from a map of string -> interface
func (t *Task) HydrateFromMap(json map[string]interface{}) {
	if title, ok := json["title"]; ok {
		t.Title = title.(string)
	}

	if description, ok := json["description"]; ok {
		t.Description = description.(string)
	}

	if points, ok := json["points"]; ok {
		t.Points = points.(float64)
	}

	if status, ok := json["status"]; ok {
		t.Status = status.(bool)
	}
}