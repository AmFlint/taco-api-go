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
	"log"
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
	TESTING__UPDATED_TITLE = "Updating tested task title"
	TESTING__UPDATED_DESCRIPTION = "Updated tested task description"
	TESTING__UPDATED_STATUS = true
	TESTING__UPDATED_POINTS = 20
)

func getBaseUrl(boardId bson.ObjectId, listId bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/lists/%s/tasks", boardId.Hex(), listId.Hex())
}

func getTaskUrl(boardId, listId, taskId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId, listId), taskId.Hex())
}

func getInvalidTaskUrl(boardId, listId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId, listId), "0")
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

func getTaskInvalidTooLongTitle() []byte {
	task := models.Task{
		Title: "testing too long string, I need more than 200 characters in order to test if validation fails for this too long title, because task title should not exceed 200 characters, at least on this application, you know",
		Description: "testing description too long title task",
		Points: 10,
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

//
func getTaskUpdateValidNoDescription() []byte {
	task := make(map[string]interface{})
	task["title"] = TESTING__UPDATED_TITLE
	task["status"] = TESTING__UPDATED_STATUS
	task["points"] = TESTING__UPDATED_POINTS
	return helpers.JsonEncode(task)
}

func getTaskUpdateValidDescription() []byte {
	task := make(map[string]interface{})
	task["description"] = TESTING__UPDATED_DESCRIPTION
	return helpers.JsonEncode(task)
}

func getTaskUpdateInvalidPoints() []byte {
	task := make(map[string]interface{})
	task["points"] = -1
	return helpers.JsonEncode(task)
}

func getTaskUpdateInvalidTitle() []byte {
	task := make(map[string]interface{})
	task["title"] = ""
	return helpers.JsonEncode(task)
}

/* ----------------------- Local Test Helpers ------------------------ */

func checkResponseCodeAndErrorMessage(t *testing.T, code int, body []byte) {
	utils.CheckResponseCode(t, code, http.StatusBadRequest)
	var res utils.ErrorResponse

	if err := json.Unmarshal(body, &res); err != nil {
		t.Errorf(utils.ERROR__UNMARSHAL_RESPONSE, err.Error())
	}
	utils.AssertNotEmpty(t, res.Message)
}

/* -----------------------------------------------------------------
   ------------------------ TEST SUITE -----------------------------
   ----------------------------------------------------------------- */

var (
	testedTaskId, boardId, listId bson.ObjectId
)

func init() {
	boardId = bson.NewObjectId()
	listId = bson.NewObjectId()
}

/* --------------------------------
   ----- Create Tasks Endpoint ----
   -------------------------------- */

func TestCreateTaskEndpoint(t *testing.T) {
	// -- Test to create a task with valid body -> Should Create task -> Return 200 w/ Task Object -- //
	t.Run("Create a Task With Valid Informations", func(t *testing.T) {
		body := getTaskValid()
		req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		// Assert that Response code is 200 / OK
		utils.CheckResponseCode(t, response.Code, http.StatusCreated)

		var responseTask models.Task

		err := json.Unmarshal(response.Body.Bytes(), &responseTask)
		if err != nil {
			t.Error("[ERR] Could not Unmarshal JSON Response, Invalid Format")
		}

		// Assert that Response task's title == created task title
		utils.AssertStringEqualsTo(t, responseTask.Title, TESTING__TASK_TITLE)
		// Assert that Response task's description == created task description
		utils.AssertStringEqualsTo(t, responseTask.Description, TESTING__TASK_DESCRIPTION)
		// Assert that reponse points == created task points
		utils.AssertFloatEqualsTo(t, responseTask.Points, TESTING__TASK_POINTS)

		utils.AssertBoolEqualsTo(t, responseTask.Status, false)

		// Save created taskId for later tests
		testedTaskId = responseTask.TaskId
	})

	// -- Test to create a task with invalid body (invalid points type) -> Should NOT Create task -> Return 400 w/ Msg/Code object -- //
	t.Run("Create a Task With Invalid Points type", func(t *testing.T) {
		body := getTaskInvalidPointType()
		req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
	})

	// -- Test to create a task with invalid body (missing/empty entries) -> Return 400 w/ Msg/Code object -- //
	t.Run("Create a Task with Missing Informations", func(t *testing.T) {
		body := getTaskInvalidMissingInformations()
		req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))

		// Execute Request and retrieve response
		response := utils.ExecuteRequest(req)

		// Assert that Response code is 400/Bad Request
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)

		checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
	})

	t.Run("Create a task with Too Long title", func(t *testing.T) {
		body := getTaskInvalidTooLongTitle()
		req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))

		response := utils.ExecuteRequest(req)

		checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
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
		taskUrl := fmt.Sprintf("%s/%s", getBaseUrl(boardId, listId), testedTaskId.Hex())

		req, _ := http.NewRequest("GET", taskUrl, nil)

		// Execute request
		response := utils.ExecuteRequest(req)

		// Check that response code == 200
		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var task models.Task
		err := json.Unmarshal(response.Body.Bytes(), &task)
		if err != nil {
			t.Errorf(utils.ERROR__UNMARSHAL_RESPONSE, err.Error())
		}

		// Assert retrieved Task has same title as the one created above, and check they have the same TaskId
		utils.AssertStringEqualsTo(t, task.Title, TESTING__TASK_TITLE)
		utils.AssertStringEqualsTo(t, task.TaskId.Hex(), testedTaskId.Hex())
	})

	t.Run("View non existing task with valid object id", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, bson.NewObjectId())
		req, _ := http.NewRequest("GET", taskUrl, nil)

		response := utils.ExecuteRequest(req)

		// Check server answers with code 404 not found
		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
		// TODO: Verify server response contains error message: create error response Struct
	})

	t.Run("view task with invalid object id", func(t *testing.T) {
		taskUrl := getInvalidTaskUrl(boardId, listId)

		req, _ := http.NewRequest("GET", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestUpdateTaskEndpoint(t *testing.T) {
	log.Print("Logging Task Update")
	t.Run("Update Task with valid informations (title, status, points) - on existing task", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)
		body := getTaskUpdateValidNoDescription()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))

		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var responseTask models.Task

		err := json.Unmarshal(response.Body.Bytes(), &responseTask)
		if err != nil {
			t.Error("[ERR] Could not Unmarshal JSON Response, Invalid Format")
		}

		// Check that fields are updated (and description is the same as before) in Task Response
		utils.AssertBoolEqualsTo(t, responseTask.Status, TESTING__UPDATED_STATUS)
		utils.AssertStringEqualsTo(t, responseTask.Title, TESTING__UPDATED_TITLE)
		utils.AssertStringEqualsTo(t, responseTask.Description, TESTING__TASK_DESCRIPTION)
		utils.AssertFloatEqualsTo(t, responseTask.Points, TESTING__UPDATED_POINTS)
	})

	t.Run("Update Task - only description with valid request", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)
		body := getTaskUpdateValidDescription()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))

		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var responseTask models.Task

		err := json.Unmarshal(response.Body.Bytes(), &responseTask)
		if err != nil {
			t.Error("[ERR] Could not Unmarshal JSON Response, Invalid Format")
		}

		// Check that fields are updated (and description is the same as before) in Task Response
		utils.AssertBoolEqualsTo(t, responseTask.Status, TESTING__UPDATED_STATUS)
		utils.AssertStringEqualsTo(t, responseTask.Title, TESTING__UPDATED_TITLE)
		utils.AssertStringEqualsTo(t, responseTask.Description, TESTING__UPDATED_DESCRIPTION)
		utils.AssertFloatEqualsTo(t, responseTask.Points, TESTING__UPDATED_POINTS)
	})

	t.Run("Update existing task with invalid points (negative)", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)
		body := getTaskUpdateInvalidPoints()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))

		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Update existing task with empty title", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)
		body := getTaskUpdateInvalidTitle()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))
		response := utils.ExecuteRequest(req)
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Update non-existing task", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, bson.NewObjectId())
		body := getTaskUpdateValidNoDescription()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Update non-valid task ID", func(t *testing.T) {
		taskUrl := getInvalidTaskUrl(boardId, listId)
		body := getTaskUpdateValidNoDescription()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestDeleteTaskEndpoint(t *testing.T) {
	// Delete existing Resource with valid Object Id
	t.Run("Delete Existing Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var m map[string]interface{}

		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Fatal("Could not unMarshal response JSON")
		}

		utils.AssertMapHasKey(t, m, "taskId")
		utils.AssertStringEqualsTo(t, m["title"].(string), TESTING__UPDATED_TITLE)
	})

	t.Run("Delete nonexisting Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskId)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Delete Task with invalid ObjectId", func(t *testing.T) {
		taskUrl := getInvalidTaskUrl(boardId, listId)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}
