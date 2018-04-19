package dao

import (
	"gopkg.in/mgo.v2"
	"github.com/AmFlint/taco-api-go/models"
	"gopkg.in/mgo.v2/bson"
)

type TaskDAO struct {
	Database *mgo.Database
}

const (
	COLLECTION = "tasks"
)

// Create a TaskDAO structure and set DAO's database, return new struct
func NewTaskDAO(db *mgo.Database) TaskDAO {
	t := TaskDAO{}
	t.SetDb(db)

	return t
}

// Prepare base queries by Initiating connection to the Collection
func prepareQuery(db *mgo.Database) *mgo.Collection {
	return db.C(COLLECTION)
}

func (t *TaskDAO) SetDb(db *mgo.Database) {
	t.Database = db
}

// Find All tasks from Database
func (t *TaskDAO) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	err := prepareQuery(t.Database).Find(bson.M{}).All(&tasks)
	return  tasks, err
}

func (t *TaskDAO) FindById(taskId bson.ObjectId) (models.Task, error) {
	var task models.Task
	err := prepareQuery(t.Database).FindId(taskId).One(&task)

	return task, err
}

func (t *TaskDAO) Delete(task *models.Task) error {
	err := prepareQuery(t.Database).Remove(&task)
	return err
}

func (t *TaskDAO) Update(task *models.Task) error {
	err := prepareQuery(t.Database).UpdateId(task.TaskId, &task)
	return err
}

func (t *TaskDAO) Insert(task *models.Task) error {
	err := prepareQuery(t.Database).Insert(&task)
	return err
}