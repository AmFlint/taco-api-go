package dao

import (
	"github.com/AmFlint/taco-api-go/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TaskDAO struct {
	Database *mgo.Database
}

const (
	TaskCollection = "tasks"
)

// Create a TaskDAO structure and set DAO's database, return new struct
func NewTaskDAO(db *mgo.Database) TaskDAO {
	t := TaskDAO{}
	t.SetDb(db)

	return t
}

// Prepare base queries by Initiating connection to the Collection
func prepareQuery(db *mgo.Database, collection string) *mgo.Collection {
	return db.C(collection)
}

func (t *TaskDAO) SetDb(db *mgo.Database) {
	t.Database = db
}

// Find All tasks from Database
func (t *TaskDAO) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := prepareQuery(t.Database, TaskCollection).Find(bson.M{}).All(&tasks)
	return tasks, err
}

func (t *TaskDAO) FindById(taskId bson.ObjectId) (models.Task, error) {
	var task models.Task
	err := prepareQuery(t.Database, TaskCollection).FindId(taskId).One(&task)

	return task, err
}

func (t *TaskDAO) Delete(task *models.Task) error {
	err := prepareQuery(t.Database, TaskCollection).Remove(&task)
	return err
}

func (t *TaskDAO) Update(task *models.Task) error {
	err := prepareQuery(t.Database, TaskCollection).UpdateId(task.TaskId, &task)
	return err
}

func (t *TaskDAO) Insert(task *models.Task) error {
	err := prepareQuery(t.Database, TaskCollection).Insert(&task)
	return err
}

// Find a Task by ID, if error return empty task with error, then delete task and return deleted task + error
func (t *TaskDAO) FindByIdAndDelete(taskId bson.ObjectId) (models.Task, error) {
	var task models.Task
	err := prepareQuery(t.Database, TaskCollection).FindId(taskId).One(&task)

	if err != nil {
		return task, err
	}
	err = t.Delete(&task)
	// Return deleted task and Error
	return task, err
}

// DeleteFromListID deletes all tasks which are attached to given listId
func (t *TaskDAO) DeleteFromListID(listID bson.ObjectId) error {
	return prepareQuery(t.Database, TaskCollection).Remove(bson.M{"listId": listID})
}
