package generator

import (
	"github.com/AmFlint/taco-api-go/tests/utils"
	"net/http"
	"github.com/AmFlint/taco-api-go/models"
	"encoding/json"
	"log"
	"testing"
	"gopkg.in/mgo.v2/bson"
	"github.com/AmFlint/taco-api-go/helpers"
	"bytes"
	"fmt"
)

var (
	boardID bson.ObjectId
	listURL string
)

func init() {
	boardsID = bson.NewObjectId()
	listURL = fmt.Sprintf("/boards/%s/lists/", boardsID.Hex())
}

// Generate a Task Entity in Database from a given Task Structure
func GenerateList(t *testing.T, list *models.List) models.List {
	// Request to API CREATE task endpoint
	reqList := helpers.JsonEncode(list)
	req, _ := http.NewRequest("POST", listURL, bytes.NewReader(reqList))
	response := utils.ExecuteRequest(req)
	// Manage response
	utils.CheckResponseCode(t, response.Code, http.StatusCreated)
	var resList models.List

	if err := json.Unmarshal(response.Body.Bytes(), &resList); err != nil {
		t.Error("Could not unmarshal List Response Body from API Create endpoint")
	}

	log.Print("List Created Properly!")
	return resList
}

// GenerateTaskAndGetID = Helper to generate a Task entity and get its ObjectID
func GenerateListAndGetID(t *testing.T, list *models.List) bson.ObjectId {
	listCreated := GenerateList(t, list)
	return listCreated.ListId
}