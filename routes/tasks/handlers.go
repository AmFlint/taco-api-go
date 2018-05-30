package tasks


import (
	"net/http"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/dao"
	"github.com/AmFlint/taco-api-go/config/database"
	"encoding/json"
	"github.com/AmFlint/taco-api-go/models"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	validator2 "gopkg.in/validator.v2"
	log "github.com/sirupsen/logrus"
	"github.com/AmFlint/taco-api-go/constants"
	"github.com/AmFlint/taco-api-go/helpers/logger"
)

//TODO: Create Middleware for initiating TaskDAO
//TODO: Create Middleware for vars["taskId"]
//TODO: Create Error message constant for Not Found
var taskLogger *log.Entry

func init() {
	taskLogger = log.WithField(constants.ResourceKeyLogger, "tasks")
}

func TaskIndexHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerListLogger, r.URL.Path, r.Method)
	taskDao := dao.NewTaskDAO(database.GetDatabaseConnection())
	tasks, err := taskDao.FindAll()
	//tasks := []models.Task {
	//	{TaskId: bson.NewObjectId(), Title: "Test Title", Description: "test description", Status: "done"},
	//	{TaskId: bson.NewObjectId(), Title: "Second task", Description: "Second task desc", Status: "in progress"},
	//}

	if err != nil {
		handlerLogger.Fatal("Could not connect to DB to retrieve Tasks")
	}

	helpers.RespondWithJson(w, 200, tasks)
}

func TaskViewHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerViewLogger, r.URL.Path, r.Method)
	//defer r.Body.Close()
	vars := mux.Vars(r)
	taskIdVar := vars["taskId"]

	if isObjectId := bson.IsObjectIdHex(taskIdVar); !isObjectId {
		handlerLogger.Warn("User provided invalid ObjectID for task Id paremeter")
		helpers.RespondWithError(w, http.StatusBadRequest, "Parameter task id is not a valid ObjectID")
		return
	}

	taskId := bson.ObjectIdHex(taskIdVar)

	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())
	task, err := taskDAO.FindById(taskId)
	if err != nil {
		handlerLogger.Warnf("Task does not exist for provided id: %s", taskId.Hex())
		helpers.RespondWithError(w, http.StatusNotFound, "Task does not exist")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, task)
	return
}

// ---- Endpoint to create a Task ---- //
func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerCreateLogger, r.URL.Path, r.Method)
	var task models.Task
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&task); err != nil {
		handlerLogger.Warnf("User sent data with wrong format, got error: %s", err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	validator := validator2.NewValidator()
	validator.SetTag("onCreate")
	if errs := validator.Validate(task); errs != nil {
		handlerLogger.Warnf("Validation failed on task %s, got error: %s", task, errs.Error())
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
	handlerLogger := logger.GenerateLogger(constants.HandlerDeleteLogger, r.URL.Path, r.Method)
	vars := mux.Vars(r)

	taskIdVar := vars["taskId"]

	if isObjectId := bson.IsObjectIdHex(taskIdVar); !isObjectId {
		handlerLogger.Warn("User provided invalid ObjectID for task Id paremeter")
		helpers.RespondWithError(w, http.StatusBadRequest, "Parameter Task is not a valid object id")
		return
	}

	taskId := bson.ObjectIdHex(taskIdVar)
	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())

	task, err := taskDAO.FindByIdAndDelete(taskId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Task Not Found")
		handlerLogger.Warnf("Task not found with id: %s", taskId.Hex())
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, task)
	return
}

// Http Method PUT on Task Resource: Update a Task
func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerUpdateLogger, r.URL.Path, r.Method)
	var task models.Task

	vars := mux.Vars(r)
	taskIdHex := vars["taskId"]

	// Check TaskId
	if isObjectId := bson.IsObjectIdHex(taskIdHex);  !isObjectId {
		handlerLogger.Warn("User provided invalid ObjectID for task Id paremeter")
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid Task Id")
		return
	}

	taskId := bson.ObjectIdHex(taskIdHex)
	taskDao := dao.NewTaskDAO(database.GetDatabaseConnection())

	// Retrieve task from database
	mainTask, err := taskDao.FindById(taskId)
	if err != nil {
		handlerLogger.Warnf("Task not found for id: %s", taskId.Hex())
		helpers.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Parse request body
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		handlerLogger.Warnf("Bad format for Request body, got error: %s", err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bodyJson := helpers.JsonEncode(body)
	// Check that request body types are correct for Task Model
	if err := json.Unmarshal(bodyJson, &task); err != nil {
		handlerLogger.Fatal(err.Error())
		helpers.RespondWithError(w, http.StatusInternalServerError, "Can not unmarshal body")
		return
	}

	// Hydrate Task from request's attributes
	mainTask.HydrateFromMap(body)

	if err := helpers.Validate(mainTask, "onCreate"); err != nil {
		handlerLogger.Warnf("Validation failed for User Input, got error: %s", err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := taskDao.Update(&mainTask); err != nil {
		handlerLogger.Fatal("Error while trying to access database, unreachable")
		helpers.RespondWithError(w, http.StatusInternalServerError, "Server Error during task Update")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, mainTask)
}