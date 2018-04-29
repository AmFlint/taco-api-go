package generator

import (
	"github.com/AmFlint/taco-api-go/models"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/AmFlint/taco-api-go/helpers"
	"bytes"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"testing"
	"encoding/json"
	"log"
)

var (
	boardsID, listID bson.ObjectId
	taskURL string
)

func init() {
	// Init board/list IDS
	// TODO: Refactor boardsID and listID when it gets implemented
	boardsID, listID = bson.NewObjectId(), bson.NewObjectId()
	// Init task URL
	taskURL = fmt.Sprintf("/boards/%s/lists/%s/tasks/", boardsID.Hex(), listID.Hex())
}

// Generate a Task Entity in Database from a given Task Structure
func GenerateTask(t *testing.T, task *models.Task) models.Task {
	// Request to API CREATE task endpoint
	reqTask := helpers.JsonEncode(task)
	req, _ := http.NewRequest("POST", taskURL, bytes.NewReader(reqTask))
	response := utils.ExecuteRequest(req)
	// Manage response
	utils.CheckResponseCode(t, response.Code, http.StatusCreated)
	var resTask models.Task

	if err := json.Unmarshal(response.Body.Bytes(), &resTask); err != nil {
		t.Error("Could not unmarshal Task Response Body from API Create endpoint")
	}

	log.Print("Task Created Properly!")
	return resTask
}

// GenerateTaskAndGetID = Helper to generate a Task entity and get its ObjectID
func GenerateTaskAndGetID(t *testing.T, task *models.Task) bson.ObjectId {
	taskCreated := GenerateTask(t, task)
	return taskCreated.TaskId
}
