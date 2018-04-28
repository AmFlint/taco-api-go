package lists

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AmFlint/taco-api-go/config/database"
	"github.com/AmFlint/taco-api-go/constants"
	"github.com/AmFlint/taco-api-go/dao"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/helpers/logger"
	"github.com/AmFlint/taco-api-go/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	validator2 "gopkg.in/validator.v2"
)

//TODO: Create Middleware for initiating TaskDAO
//TODO: Create Middleware for vars["taskId"]
//TODO: Create Error message constant for Not Found

var listLogger *log.Entry

func init() {
	listLogger = log.WithField(constants.HandlerKeyLogger, constants.ResourceListsLogger)
}

// ListCreateHandler -> Handler for List Creation Endpoint ---- //
func ListCreateHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerCreateLogger, r.URL.Path, r.Method)
	list := models.NewList()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&list); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate User Input for List Entity
	validator := validator2.NewValidator()
	validator.SetTag("onCreate")
	if errs := validator.Validate(list); errs != nil {
		handlerLogger.Errorf("Validating List for endpoint List creation with error: %s", errs.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, errs.Error())
		return
	}

	// Manage Database insertion for this new List
	list.ListId = bson.NewObjectId()
	listDao := dao.NewListDao()
	if err := listDao.Insert(&list); err != nil {
		handlerLogger.Error("Could not insert to database")
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusCreated, list)
}

// ListDeleteHandler -> Handler for List Deletion Endpoint
func ListDeleteHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerDeleteLogger, r.URL.Path, r.Method)

	vars := mux.Vars(r)
	listIDVars := vars["listId"]

	if isObjectID := bson.IsObjectIdHex(listIDVars); !isObjectID {
		handlerLogger.Warn("User provided invalid Object ID for parmeters listId")
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	listID := bson.ObjectIdHex(listIDVars)

	listDAO := dao.NewListDao()
	taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())

	list, err := listDAO.FindByIDAndDelete(listID)
	if err != nil {
		handlerLogger.Warnf("List not found with id: %s", listIDVars)
		helpers.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	if err := taskDAO.DeleteFromListID(listID); err != nil {
		// TODO: Check whether to respond now or ignore as list is already deleted, or better: Use a transaction
		handlerLogger.Error("Could not delete tasks from database")
		//helpers.RespondWithError(w, http.StatusInternalServerError, "Could remove tasks attached to given list")
		//return
	}

	helpers.RespondWithJson(w, http.StatusOK, list)
}
