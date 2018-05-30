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

	// Make sure that request body is not empty
	if r.Body == nil {
		handlerLogger.Warn("Empty Request Body")
		helpers.RespondWithError(w, http.StatusBadRequest, "Empty Request Body")
		return
	}

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

// ListViewHandler -> Handler to View List Endpoint
func ListViewHandler(w http.ResponseWriter, r *http.Request) {
	handlerLogger := logger.GenerateLogger(constants.HandlerViewLogger, r.URL.Path, r.Method)
	vars := mux.Vars(r)
	listIdVars := vars["listId"]

	if isObjectId := bson.IsObjectIdHex(listIdVars); !isObjectId {
		handlerLogger.Warn("Invalid Object ID for list")
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	listID := bson.ObjectIdHex(listIdVars)
	listDAO := dao.NewListDao()

	list, err := listDAO.FindByID(listID)
	// List not found
	if err != nil {
		handlerLogger.Warnf("List not found with id: %s", listIdVars)
		helpers.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	helpers.RespondWithJson(w, http.StatusOK, list)
	return
}

// ListUpdateHandler -> Handler to Update a List Endpoint
func ListUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var mainList models.List
	handlerLogger := logger.GenerateLogger(constants.HandlerUpdateLogger, r.URL.Path, r.Method)
	vars := mux.Vars(r)

	listIdVars := vars["listId"]

	if isObjectId := bson.IsObjectIdHex(listIdVars); !isObjectId {
		handlerLogger.Warn("Invalid Object ID for list")
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid ObjectID")
		return
	}

	listID := bson.ObjectIdHex(listIdVars)
	listDAO := dao.NewListDao()

	// TODO: Create / Use method UpdateByID -> check error type for response
	list, err := listDAO.FindByID(listID)
	if err != nil {
		handlerLogger.Warnf("List not found with id: %s", listIdVars)
		helpers.RespondWithError(w, http.StatusNotFound, "List not found")
		return
	}

	// Make sure that request body is not empty
	if r.Body == nil {
		handlerLogger.Warn("Received Empty request body")
		helpers.RespondWithError(w, http.StatusBadRequest, "Empty body")
		return
	}

	// Manage HTTP Request Body
	// Parse request body
	var body map[string]interface{}
	if err := helpers.DecodeBody(r.Body, &body); err != nil {
		handlerLogger.Warnf("Bad format for Request body, got error: %s", err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bodyJson := helpers.JsonEncode(body)
	// Validate Request body types against List data structure
	if err := json.Unmarshal(bodyJson, &mainList); err != nil {
		handlerLogger.Fatal(err.Error())
		helpers.RespondWithError(w, http.StatusInternalServerError, "Can not unmarshal body")
		return
	}

	// Validate List from Request Body
	if err := helpers.Validate(mainList, "onCreate"); err != nil {
		handlerLogger.Warnf("Could not validate List model, received error: %s", err.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, "Could not validate given data, errors: " + err.Error())
		return
	}

	list.HydrateFromMap(body)
	if err := listDAO.Update(&list); err != nil {
		handlerLogger.Warnf("Could not update list with id: %s, got error: %s", listIdVars, err.Error())
		helpers.RespondWithError(w, http.StatusInternalServerError, "Could not reach Database")
		return
	}
	handlerLogger.Infof("I fucked it all: %s", list)

	helpers.RespondWithJson(w, http.StatusOK, mainList)
}
