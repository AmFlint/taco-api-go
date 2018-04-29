package tasks

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
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
	"github.com/AmFlint/taco-api-go/tests/utils/generator"
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
	// Testing data constants
	invalidTaskTitle = "Invalid task title"
	invalidTaskDescription = "Invalid task description"
	testingTaskTitle = "Testing title"
	testingTaskDescription = "Testing description"
	testingTaskPoints = 9
	testingUpdatedTitle = "Updating tested task title"
	testingUpdatedDescription = "Updated tested task description"
	testingUpdatedStatus = true
	testingUpdatedPoints = 20

	// ---- Data Generation properties ---- //
	// Update endpoint
	genTaskForUpdateTitle = "About to be updated title"
	genTaskForUpdateDescription = "About to be updated description"
	genTaskForUpdatePoints = 6
	// Delete endpoint
	genTaskForDeleteTitle = "About to be deleted title"
	genTaskForDeleteDescription = "About to be deleted description"
	genTaskForDeletePoints = 15
	// View endpoint
	genTaskForViewTitle = "About to be viewed title"
	genTaskForViewDescription = "About to be viewed description"
	genTaskForViewPoints = 5
)

// get Base URL for Tasks endpoints
func getBaseUrl(boardId bson.ObjectId, listId bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/lists/%s/tasks", boardId.Hex(), listId.Hex())
}

// Get Task URL for View/delete/update endpoints
func getTaskUrl(boardId, listId, taskId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId, listId), taskId.Hex())
}

// Get and invalid URL (bad format) for task endpoints
func getInvalidTaskUrl(boardId, listId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s", getBaseUrl(boardId, listId), "0")
}

/* ------------------------------------------
   ----------- Task Generation Data ------------
   ------------------------------------------ */

// Get Task Entity for Update endpoint
func getTaskForUpdate() *models.Task {
	return &models.Task{
		Title: genTaskForUpdateTitle,
		Description: genTaskForUpdateDescription,
		Points: genTaskForUpdatePoints,
	}
}

// Get Task Entity for Delete Endpoint
func getTaskForDelete() *models.Task {
	return &models.Task{
		Title: genTaskForDeleteTitle,
		Description: genTaskForDeleteDescription,
		Points: genTaskForDeletePoints,
	}
}

// Get Task Entity for View Endpoint
func getTaskForView() *models.Task {
	return &models.Task{
		Title: genTaskForViewTitle,
		Description: genTaskForViewDescription,
		Points: genTaskForViewPoints,
	}
}


/* ------------------------------------------
   ----------- Task Testing Data ------------
   ------------------------------------------ */

// -- Get Json encoded (stringified) Struct for invalid task -> Wrong Points entry -- //
func getTaskInvalidPointType() []byte {
	task := taskInvalidPointsType{
		Title: invalidTaskTitle,
		Description: invalidTaskDescription,
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
		Title: testingTaskTitle,
		Description: testingTaskDescription,
		Points: testingTaskPoints,
	}
	return helpers.JsonEncode(task)
}

// Get a json encoded Task with no description
func getTaskUpdateValidNoDescription() []byte {
	task := make(map[string]interface{})
	task["title"] = testingUpdatedTitle
	task["status"] = testingUpdatedStatus
	task["points"] = testingUpdatedPoints
	return helpers.JsonEncode(task)
}

func getTaskUpdateValidDescription() []byte {
	task := make(map[string]interface{})
	task["description"] = testingUpdatedDescription
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

func TestMain(m *testing.M) {
	testconfig.Init(m)
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
		utils.AssertStringEqualsTo(t, responseTask.Title, testingTaskTitle)
		// Assert that Response task's description == created task description
		utils.AssertStringEqualsTo(t, responseTask.Description, testingTaskDescription)
		// Assert that reponse points == created task points
		utils.AssertFloatEqualsTo(t, responseTask.Points, testingTaskPoints)

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
	// Generate Task to test
	testedTaskID := generator.GenerateTaskAndGetID(t, getTaskForView())

	// Testing view existing Task with Valid ObjectID
	t.Run("View existing task with valid object id", func(t *testing.T) {
		//taskId := utils.TaskCreate(getTaskValid(), t)
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)

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
		utils.AssertStringEqualsTo(t, task.Title, genTaskForViewTitle)
		utils.AssertStringEqualsTo(t, task.TaskId.Hex(), testedTaskID.Hex())
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
	// Generating Task To test≤
	testedTaskID := generator.GenerateTaskAndGetID(t, getTaskForUpdate())


	t.Run("Update Task with valid informations (title, status, points) - on existing task", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)
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
		utils.AssertBoolEqualsTo(t, responseTask.Status, testingUpdatedStatus)
		utils.AssertStringEqualsTo(t, responseTask.Title, testingUpdatedTitle)
		utils.AssertStringEqualsTo(t, responseTask.Description, genTaskForUpdateDescription)
		utils.AssertFloatEqualsTo(t, responseTask.Points, testingUpdatedPoints)
	})

	t.Run("Update Task - only description with valid request", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)
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
		utils.AssertBoolEqualsTo(t, responseTask.Status, testingUpdatedStatus)
		utils.AssertStringEqualsTo(t, responseTask.Title, testingUpdatedTitle)
		utils.AssertStringEqualsTo(t, responseTask.Description, testingUpdatedDescription)
		utils.AssertFloatEqualsTo(t, responseTask.Points, testingUpdatedPoints)
	})

	t.Run("Update existing task with invalid points (negative)", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)
		body := getTaskUpdateInvalidPoints()

		req, _ := http.NewRequest("PUT", taskUrl, bytes.NewReader(body))

		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Update existing task with empty title", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)
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
	// Generating Task To test
	testedTaskID := generator.GenerateTaskAndGetID(t, getTaskForDelete())

	// Delete existing Resource with valid Object Id
	t.Run("Delete Existing Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, testedTaskID)

		req, _ := http.NewRequest("DELETE", taskUrl, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var m map[string]interface{}

		err := json.Unmarshal(response.Body.Bytes(), &m)
		if err != nil {
			t.Fatal("Could not unMarshal response JSON")
		}

		utils.AssertMapHasKey(t, m, "taskId")
		utils.AssertStringEqualsTo(t, m["title"].(string), genTaskForDeleteTitle)
		utils.AssertStringEqualsTo(t, m["description"].(string), genTaskForDeleteDescription)
		utils.AssertFloatEqualsTo(t, m["points"].(float64), genTaskForDeletePoints)
	})

	t.Run("Delete nonexisting Task with Valid ObjectId", func(t *testing.T) {
		taskUrl := getTaskUrl(boardId, listId, bson.NewObjectId())

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
