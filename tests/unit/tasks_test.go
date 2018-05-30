package unit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/AmFlint/taco-api-go/config/database"
	"github.com/AmFlint/taco-api-go/dao"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/models"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"gopkg.in/mgo.v2/bson"
)

/* -----------------------------------------------------------------
   ----------------------- Configuration ---------------------------
   ----------------------------------------------------------------- */

const (
	// Testing data constants
	invalidTaskTitle          = "Invalid task title"
	invalidTaskDescription    = "Invalid task description"
	testingTaskTitle          = "Testing title"
	testingTaskDescription    = "Testing description"
	testingTaskPoints         = 9
	testingUpdatedTitle       = "Updating tested task title"
	testingUpdatedDescription = "Updated tested task description"
	testingUpdatedStatus      = true
	testingUpdatedPoints      = 20

	// ---- Data Generation properties ---- //
	// Update endpoint
	genTaskForUpdateTitle       = "About to be updated title"
	genTaskForUpdateDescription = "About to be updated description"
	genTaskForUpdatePoints      = 6
	// Delete endpoint
	genTaskForDeleteTitle       = "About to be deleted title"
	genTaskForDeleteDescription = "About to be deleted description"
	genTaskForDeletePoints      = 15
	// View endpoint
	genTaskForViewTitle       = "About to be viewed title"
	genTaskForViewDescription = "About to be viewed description"
	genTaskForViewPoints      = 5
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
		Title:       genTaskForUpdateTitle,
		Description: genTaskForUpdateDescription,
		Points:      genTaskForUpdatePoints,
	}
}

// Get Task Entity for Delete Endpoint
func getTaskForDelete() *models.Task {
	return &models.Task{
		Title:       genTaskForDeleteTitle,
		Description: genTaskForDeleteDescription,
		Points:      genTaskForDeletePoints,
	}
}

// Get Task Entity for View Endpoint
func getTaskForView() *models.Task {
	return &models.Task{
		Title:       genTaskForViewTitle,
		Description: genTaskForViewDescription,
		Points:      genTaskForViewPoints,
	}
}

/* ------------------------------------------
   ----------- Task Testing Data ------------
   ------------------------------------------ */

// -- Get Json encoded (stringified) Struct for invalid task -> Wrong Points entry -- //
func getTaskInvalidPointType() []byte {
	task := make(map[string]interface{})
	task["title"] = invalidTaskTitle
	task["description"] = invalidTaskDescription
	task["points"] = false
	return helpers.JsonEncode(task)
}

// -- Get Json encoded (stringified) Struct for invalid task -> Empty/Missing entries -- //
func getTaskInvalidMissingInformations() []byte {
	task := models.Task{
		Title:       "",
		Description: "",
	}
	return helpers.JsonEncode(task)
}

func getTaskInvalidTooLongTitle() []byte {
	task := models.Task{
		Title:       "testing too long string, I need more than 200 characters in order to test if validation fails for this too long title, because task title should not exceed 200 characters, at least on this application, you know",
		Description: "testing description too long title task",
		Points:      10,
	}
	return helpers.JsonEncode(task)
}

// -- Get Json encoded (stringified) Struct for Valid Task -- //
func getTaskValid() []byte {
	task := models.Task{
		Title:       testingTaskTitle,
		Description: testingTaskDescription,
		Points:      testingTaskPoints,
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
	boardId, listId bson.ObjectId
)

func init() {
	boardId = bson.NewObjectId()
	listId = bson.NewObjectId()
}

/* --------------------------------
   ----- Create Tasks Endpoint ----
   -------------------------------- */

func TestUnitCreateTask(t *testing.T) {
	// -- Test to create a task with valid body -> Should Create task -> Return 200 w/ Task Object -- //
	t.Run("Create a Task With Valid Informations", func(t *testing.T) {
		body := getTaskValid()
		url := getBaseUrl(boardId, listId)

		response, err := utils.ExecuteRequestAndGetResponse("POST", url, bytes.NewReader(body))
		if err != nil {
			t.Error("Error processing HTTP Request and Unmarshaling Response")
		}

		if _, ok := response["taskId"]; !ok {
			t.Error("taskId Not found in JSON Response")
		}

		taskDAO := dao.NewTaskDAO(database.GetDatabaseConnection())
		taskID := bson.ObjectIdHex(response["taskId"].(string))
		testedTask, err := taskDAO.FindById(taskID)
		if err != nil {
			t.Error("Error while fetching Created task inside Database")
		}
		utils.AssertStringEqualsTo(t, testedTask.Title, testingTaskTitle)
		utils.AssertStringEqualsTo(t, testedTask.Description, testingTaskDescription)
		utils.AssertBoolEqualsTo(t, testedTask.Status, false)
		utils.AssertFloatEqualsTo(t, testedTask.Points, testingTaskPoints)
	})

	//// -- Test to create a task with invalid body (invalid points type) -> Should NOT Create task -> Return 400 w/ Msg/Code object -- //
	//t.Run("Create a Task With Invalid Points type", func(t *testing.T) {
	//	body := getTaskInvalidPointType()
	//	req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))
	//
	//	// Execute Request and retrieve response
	//	response := utils.ExecuteRequest(req)
	//
	//	checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
	//})
	//
	//// -- Test to create a task with invalid body (missing/empty entries) -> Return 400 w/ Msg/Code object -- //
	//t.Run("Create a Task with Missing Informations", func(t *testing.T) {
	//	body := getTaskInvalidMissingInformations()
	//	req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))
	//
	//	// Execute Request and retrieve response
	//	response := utils.ExecuteRequest(req)
	//
	//	// Assert that Response code is 400/Bad Request
	//	utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	//
	//	checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
	//})
	//
	//t.Run("Create a task with Too Long title", func(t *testing.T) {
	//	body := getTaskInvalidTooLongTitle()
	//	req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), bytes.NewReader(body))
	//
	//	response := utils.ExecuteRequest(req)
	//
	//	checkResponseCodeAndErrorMessage(t, response.Code, response.Body.Bytes())
	//})
	//
	//t.Run("Create a Task with empty request body", func(t *testing.T) {
	//	req, _ := http.NewRequest("POST", getBaseUrl(boardId, listId), nil)
	//	// Execute Request and retrieve response
	//	response := utils.ExecuteRequest(req)
	//	// Assert that Response code is 400/Bad Request
	//	utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	//})
}
