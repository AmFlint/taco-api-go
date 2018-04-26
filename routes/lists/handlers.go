package lists


import (
	"net/http"
	"github.com/AmFlint/taco-api-go/models"
	"encoding/json"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/dao"
	validator2 "gopkg.in/validator.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

//TODO: Create Middleware for initiating TaskDAO
//TODO: Create Middleware for vars["taskId"]
//TODO: Create Error message constant for Not Found

func ListCreateHandler(w http.ResponseWriter, r *http.Request) {
	var list models.List

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&list); err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Validate User Input for List Entity
	validator := validator2.NewValidator()
	validator.SetTag("onCreate")
	if errs := validator.Validate(list); errs != nil {
		log.Printf("[ERR] Validating List for endpoint List creation with error: %s", errs.Error())
		helpers.RespondWithError(w, http.StatusBadRequest, errs.Error())
		return
	}

	// Manage Database insertion for this new List
	list.ListId = bson.NewObjectId()
	listDao := dao.NewListDao()
	if err := listDao.Insert(&list); err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJson(w, http.StatusCreated, list)
}
