package routes

import (
	"net/http"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/dao"
	"github.com/AmFlint/taco-api-go/config/database"
	"log"
	"encoding/json"
	"github.com/AmFlint/taco-api-go/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
)

func TaskIndexHandler(w http.ResponseWriter, r *http.Request) {
	taskDao := dao.NewTaskDAO(database.GetDatabaseConnection())
	tasks, err := taskDao.FindAll()
	//tasks := []models.Task {
	//	{TaskId: bson.NewObjectId(), Title: "Test Title", Description: "test description", Status: "done"},
	//	{TaskId: bson.NewObjectId(), Title: "Second task", Description: "Second task desc", Status: "in progress"},
	//}

	if err != nil {
		log.Fatal("Could not connect to DB to retrieve Tasks")
	}

	helpers.RespondWithJson(w, 200, tasks)
}

func TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	//defer r.Body.Close()
	vars := mux.Vars(r)
	taskIdVar := vars["taskId"]

	if isObjectId := bson.IsObjectIdHex(taskIdVar); !isObjectId {
		helpers.RespondWithError(w, http.StatusBadRequest, "Parameter task id is not a valid ObjectID")
		return
	}

	taskId := bson.ObjectIdHex(taskIdVar)

	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())
	task, err := taskDAO.FindById(taskId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Task does not exist")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, task)
	return
}

// ---- Endpoint to create a Task ---- //
func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	decoder := json.NewDecoder(r.Body)
	taskInput := make(map[string]interface{})

	if err := decoder.Decode(&taskInput); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ERROR__INVALID_PLAYLOAD)
		return
	}

	taskInput["status"] = models.GetDefaultTaskStatus()
	// Validate User Input for Task Creation
	validationErrors := models.ValidateTask(taskInput)

	// If errors were fount during validation, send error messages and bad request status code
	if len(validationErrors) != 0 {
		helpers.RespondWithErrors(w, http.StatusBadRequest, validationErrors)
		return
	}

	validInput, _ := json.Marshal(taskInput)

	if err := json.Unmarshal(validInput, &task); err != nil {
		log.Print(task)
		helpers.RespondWithError(w, http.StatusBadRequest, helpers.ERROR__INVALID_PLAYLOAD)
		return
	}

	task.SetDefaultStatus()
	task.TaskId = bson.NewObjectId()

	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())

	if err := taskDAO.Insert(&task); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	}
	helpers.RespondWithJson(w, http.StatusCreated, task)
}

// Http Method DELETE on Task resource: Delete a Task
func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {

}

// Http Method PUT on Task Resource: Update a Task
func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {

}