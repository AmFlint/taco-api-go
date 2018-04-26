package dao

import (
	"gopkg.in/mgo.v2"
	"github.com/AmFlint/taco-api-go/config/database"
	"github.com/AmFlint/taco-api-go/models"
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

func (l *ListDAO) findById() {

}

func (l *ListDAO) Insert(list *models.List) error {
	err := prepareQuery(l.Database, ListCollection).Insert(&list)
	return err
}
