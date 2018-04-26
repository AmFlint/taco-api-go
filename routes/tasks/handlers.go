package tasks


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
	validator2 "gopkg.in/validator.v2"
)

//TODO: Create Middleware for initiating TaskDAO
//TODO: Create Middleware for vars["taskId"]
//TODO: Create Error message constant for Not Found
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

	if err := decoder.Decode(&task); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	validator := validator2.NewValidator()
	validator.SetTag("onCreate")
	if errs := validator.Validate(task); errs != nil {
		log.Printf("Validation failed on task %s, got error: %s", task, errs.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, errs.Error())
		return
	}

	task.SetDefaultStatus()
	task.TaskId = bson.NewObjectId()

	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())

	if err := taskDAO.Insert(&task); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	helpers.RespondWithJson(w, http.StatusCreated, task)
}

// Http Method DELETE on Task resource: Delete a Task
func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r);

	taskIdVar := vars["taskId"]

	if isObjectId := bson.IsObjectIdHex(taskIdVar); !isObjectId {
		helpers.RespondWithError(w, http.StatusBadRequest, "Parameter Task is not a valid object id")
		log.Print("[ERR] Parameter :taskId is not a valid ObjectId in DELETE tasks")
		return
	}

	taskId := bson.ObjectIdHex(taskIdVar)
	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())

	task, err := taskDAO.FindByIdAndDelete(taskId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Task Not Found")
		log.Print("[ERR] Task does not exist in DELETE task")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, task)
	return
}

// Http Method PUT on Task Resource: Update a Task
func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	vars := mux.Vars(r)
	taskIdHex := vars["taskId"]

	// Check TaskId
	if isObjectId := bson.IsObjectIdHex(taskIdHex);  !isObjectId {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Task Id")
		return
	}

	taskId := bson.ObjectIdHex(taskIdHex)
	taskDao := dao.NewTaskDAO(database.GetDatabaseConnection())

	// Retrieve task from database
	mainTask, err := taskDao.FindById(taskId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Parse request body
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bodyJson := helpers.JsonEncode(body)
	if err := json.Unmarshal(bodyJson, &task); err != nil {
		log.Print(err.Error())
		helpers.RespondWithError(w, http.StatusInternalServerError, "Can not unmarshal body")
		return
	}

	// Hydrate Task from request's attributes
	mainTask.HydrateFromMap(body)

	validator := validator2.NewValidator()
	validator.SetTag("onCreate")

	if errs := validator.Validate(mainTask); errs != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, errs.Error())
		return
	}

	if err := taskDao.Update(&mainTask); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error during task Update")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, mainTask)
}