package dao

import (
	"github.com/AmFlint/taco-api-go/config/database"
	"github.com/AmFlint/taco-api-go/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ListDAO struct {
	Database *mgo.Database
}

const (
	ListCollection = "lists"
)

// Create a TaskDAO structure and set DAO's database, return new struct
func NewListDao() ListDAO {
	l := ListDAO{}
	l.SetDb(database.GetDatabaseConnection())

	return l
}

func (l *ListDAO) SetDb(db *mgo.Database) {
	l.Database = db
}

// Delete a list from the database
func (l *ListDAO) Delete(list *models.List) error {
	err := prepareQuery(l.Database, ListCollection).Remove(&list)
	return err
}

// FindByID -> Find a List by its id
func (l *ListDAO) FindByID(listID bson.ObjectId) (models.List, error) {
	var list models.List
	err := prepareQuery(l.Database, ListCollection).FindId(listID).One(&list)
	return list, err
}

// Insert a list to the database
func (l *ListDAO) Insert(list *models.List) error {
	err := prepareQuery(l.Database, ListCollection).Insert(&list)
	return err
}

// FindByIDAndDelete -> Find a List by ID, if error return empty list with error, then delete list and return deleted list + error
func (l *ListDAO) FindByIDAndDelete(listID bson.ObjectId) (models.List, error) {
	var list models.List
	err := prepareQuery(l.Database, ListCollection).FindId(listID).One(&list)

	if err != nil {
		return list, err
	}
	err = l.Delete(&list)
	// Return deleted list and Error
	return list, err
}

// Update - Update a List Entity
func (l *ListDAO) Update(list *models.List) error {
	return prepareQuery(l.Database, ListCollection).UpdateId(list.ListId, list)
}
