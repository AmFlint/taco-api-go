package tests

import (
	"testing"
	"github.com/AmFlint/taco-api-go/models"
	"net/http"
	"bytes"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

/* -----------------------------------------------------------------
   ----------------------- Configuration ---------------------------
   ----------------------------------------------------------------- */

// -- Structure for Invalid Tasks due to invalid type for Points entry -- //
type taskInvalidPointsType struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Status      string `bson:"status" json:"status"`
	Points      bool `bson:"points" json:"points"`
}

const (
	INVALID__TASK_TITLE = "Invalid task title"
	INVALID__TASK_DESCRIPTION = "Invalid task description"
	TESTING__TASK_TITLE = "Testing title"
	TESTING__TASK_DESCRIPTION = "Testing description"
	TESTING__TASK_POINTS = 9
)

func getBaseUrl(boardId bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/tasks", boardId.Hex())
}

func getTaskUrl(boardId bson.ObjectId, taskId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId), taskId.Hex())
}

func getInvalidTaskUrl(boardId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId), "0")
}

// -- Get Json encoded (stringified) Struct for invalid task -> Wrong Points entry -- //
func getTaskInvalidPointType() []byte {
	task := taskInvalidPointsType{
		Title: INVALID__TASK_TITLE,
		Description: INVALID__TASK_DESCRIPTION,
		Points: false,
	}
	return helpers.JsonEncode(task)
}

// -- Get Json encoded (stringified) Struct for invalid task -> Empty/Missing entries -- //
func getTaskInvalidMissingInformations() []byte {
	task := models.Task{
		Title: "",
		Description: "",
	}
	return helpers.JsonEncode(task)
}

// -- Get Json encoded (stringified) Struct for Valid Task -- //
func getTaskValid() []byte {
	task := models.Task{
		Title: TESTING__TASK_TITLE,
		Description: TESTING__TASK_DESCRIPTION,
		Points: TESTING__TASK_POINTS,
	}
	return helpers.JsonEncode(task)
}

/* -----------------------------------------------------------------
   ------------------------ TEST SUITE -----------------------------
   ----------------------------------------------------------------- */

/* --------------------------------
   ----- Create Tasks Endpoint ----
   -------------------------------- */

var testedTaskId bson.ObjectId

func TestCreateTaskEndpoint(t *testing.T) {
	// -- Test to create a task with valid body -> Should Create task -> Return 200 w/ Task Object -- //
	t.Run("Create a Task With Valid Informations", func(t *testing.T) {
		body := getTaskValid()
		req, _ := http.NewRequest("POST", "/boards/1/tasks", bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		// Assert that Response code is 200 / OK
		utils.CheckResponseCode(t, response.Code, http.StatusCreated)

		var m map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Error("[ERR] Could not Unmarshal JSON Response, Invalid Format")
		}

		// Assert that Response task's title == created task title
		utils.AssertStringEqualsTo(t, m["title"].(string), TESTING__TASK_TITLE)

		// Assert that Response task's description == created task description
		utils.AssertStringEqualsTo(t, m["description"].(string), TESTING__TASK_DESCRIPTION)

		// Assert that reponse points == created task points
		responsePoints := m["points"].(float64)
		utils.AssertFloatEqualsTo(t, responsePoints, TESTING__TASK_POINTS)

		// Type Assertion response Map, key "taskId" to type ObjectId for later tests
		taskId := m["taskId"].(string)
		testedTaskId = bson.ObjectIdHex(taskId)
	})

	// -- Test to create a task with invalid body (invalid points type) -> Should NOT Create task -> Return 400 w/ Msg/Code object -- //
	t.Run("Create a Task With Invalid Points type", func(t *testing.T) {
		body := getTaskInvalidPointType()
		req, _ := http.NewRequest("POST", "/boards/1/tasks", bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		// Assert that Response Code is 400 / Bad Request
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)

		var m map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Error("[Error] Could not Unmarshal HTTP Response Body")
		}

		//// Assert that
		utils.AssertMapHasKey(t, m, "errors")
	})

	// -- Test to create a task with invalid body (missing/empty entries) -> Return 400 w/ Msg/Code object -- //
	t.Run("Create a Task with Missing Informations", func(t *testing.T) {
		body := getTaskInvalidMissingInformations()
		req, _ := http.NewRequest("POST", "/boards/1/tasks", bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		// Assert that Response code is 400/Bad Request
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)

		var m map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Error("[Error] Could not Unmarshal HTTP Response Body")
		}
		utils.AssertMapHasKey(t, m, "errors")
	})
}

/* --------------------------------
   ----- View Task Endpoint ----
   -------------------------------- */

func TestViewTaskEndpoint(t *testing.T) {
	// Testing view existing Task with Valid ObjectID
	t.Run("View existing task with valid object id", func(t *testing.T) {
		//taskId := utils.TaskCreate(getTaskValid(), t)
		// TODO: change URL, add "/" to every url
		taskUrl := fmt.Sprintf("%s/%s", getBaseUrl(bson.NewObjectId()), testedTaskId.Hex())

		req, _ := http.NewRequest("GET", taskUrl, nil)

		// Execute request
		response := utils.ExecuteRequest(req)

		// Check that response code == 200
		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var m map[string]interface{}
		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Error("[Error] Could not Unmarshal HTTP Response Body")
		}

		utils.AssertMapHasKey(t, m, "title")
		utils.AssertStringEqualsTo(t, m["title"].(string), TESTING__TASK_TITLE)
		utils.AssertStringEqualsTo(t, m["taskId"].(string), testedTaskId.Hex())
	})

	t.Run("View non existing task with valid object id", func(t *testing.T) {
		taskUrl := getTaskUrl(bson.NewObjectId(), bson.NewObjectId())
		req, _ := http.NewRequest("GET", taskUrl, nil)

		response := utils.ExecuteRequest(req)

		// Check server answers with code 404 not found
		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
		// TODO: Verify server response contains error message: create error response Struct
	})

	t.Run("view task with invalid object id", func(t *testing.T) {
		taskUrl := getInvalidTaskUrl(bson.NewObjectId())

		req, _ := http.NewRequest("GET", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestDeleteTaskEndpoint(t *testing.T) {
	// Delete existing Resource with valid Object Id
	t.Run("Delete Existing Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(bson.NewObjectId(), testedTaskId)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var m map[string]interface{}

		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Fatal("Could not unMarshal response JSON")
		}

		utils.AssertMapHasKey(t, m, "taskId")
		utils.AssertStringEqualsTo(t, m["title"].(string), TESTING__TASK_TITLE)
	})

	t.Run("Delete nonexisting Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(bson.NewObjectId(), testedTaskId)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Delete Task with invalid ObjectId", func(t *testing.T) {
		taskUrl := getInvalidTaskUrl(bson.NewObjectId())

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}
