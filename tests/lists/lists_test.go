package lists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/models"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
	"gopkg.in/mgo.v2/bson"
	"github.com/AmFlint/taco-api-go/tests/utils/generator"
)

const (
	TESTING__LIST_NAME  = "Testing list"
	TESTING__LIST_ORDER = 1

	// data creation
	// delete endpoint
	genListForDeleteName = "about to be deleted"
	// view endpoint
	genListForViewName = "about to be viewed"
)

// TODO: Create Helpers for Resource creations -> Tests run in parrallell which means reusing an id from above test may not work

func getListsBaseUrl(boardId bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/lists/", boardId.Hex())
}

func getlistURL(boardId, listId bson.ObjectId) string {
	return fmt.Sprintf("%s%s/", getListsBaseUrl(boardId), listId.Hex())
}

func getInvalidlistURL(boardId bson.ObjectId) string {
	return fmt.Sprintf("%s%s/", getListsBaseUrl(boardId), "2")
}

func getListForDelete() *models.List {
	list := models.NewList()
	list.Name = genListForDeleteName
	return &list
}

func getListForView() *models.List {
	list := models.NewList()
	list.Name = genListForViewName
	return &list
}

// Configuration for basic Lists

func getValidList() []byte {
	list := make(map[string]interface{})
	list["name"] = TESTING__LIST_NAME
	return helpers.JsonEncode(list)
}

func getInvalidListEmptyTitle() []byte {
	list := make(map[string]interface{})
	list["name"] = ""
	return helpers.JsonEncode(list)
}

func getInvalidListBadTitle() []byte {
	list := make(map[string]interface{})
	list["title"] = "testing1209"
	return helpers.JsonEncode(list)
}

/* ------------------------------------------
   -------------- Test Suite ----------------
   ------------------------------------------ */

var (
	boardId bson.ObjectId
)

func init() {
	boardId = bson.NewObjectId()
}

func TestMain(m *testing.M) {
	testconfig.Init(m)
}

// ---- Test Create Endpoint ---- //
func TestCreateListEndpoint(t *testing.T) {
	t.Run("Create List with valid informations", func(t *testing.T) {
		listURL := getListsBaseUrl(boardId)
		list := getValidList()

		req, _ := http.NewRequest("POST", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		var createdList models.List
		if err := json.Unmarshal(response.Body.Bytes(), &createdList); err != nil {
			t.Error("[Error] in Create List endpoint, could not unmarshal response body")
		}

		utils.CheckResponseCode(t, response.Code, http.StatusCreated)

		utils.AssertStringEqualsTo(t, createdList.Name, TESTING__LIST_NAME)
		// TODO: Implement "order" tests when board is implemented
		//utils.AssertIntEqualsTo(t, createdList.Order, 1)
	})

	t.Run("Create List with invalid user data - empty title", func(t *testing.T) {
		listURL := getListsBaseUrl(boardId)
		list := getInvalidListEmptyTitle()
		req, _ := http.NewRequest("POST", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Create List with Invalid Title (bad format)", func(t *testing.T) {
		listURL := getListsBaseUrl(boardId)
		list := getInvalidListBadTitle()

		req, _ := http.NewRequest("POST", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

// ---- Test View Endpoint ---- //
func TestViewListHandler(t *testing.T) {
	testedListID := generator.GenerateListAndGetID(t, getListForView())

	t.Run("View an existing list with valid ID", func(t *testing.T) {
		listURL := getlistURL(boardId, testedListID)

		req, _ := http.NewRequest("GET", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var list models.List
		if err := json.Unmarshal(response.Body.Bytes(), &list); err != nil {
			t.Error("Could not unmarshal response body")
		}

		// Assertions
		utils.AssertStringEqualsTo(t, list.Name, genListForViewName)
	})

	t.Run("View a non existing list with valid ID", func(t *testing.T) {
			listURL := getlistURL(boardId, bson.NewObjectId())

			req, _ := http.NewRequest("GET", listURL, nil)
			response := utils.ExecuteRequest(req)

			utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("View a list with invalid ID", func(t *testing.T) {
		listURL := getInvalidlistURL(boardId)

		req, _ := http.NewRequest("GET", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

// ---- Test Delete EndPoint ---- //
func TestDeleteListHandler(t *testing.T) {
	testedListID := generator.GenerateListAndGetID(t, getListForDelete())

	t.Run("Delete an Existing task with valid Id", func(t *testing.T) {
		listURL := getlistURL(boardId, testedListID)

		req, _ := http.NewRequest("DELETE", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		// Decode the Response
		var list models.List
		err := json.Unmarshal(response.Body.Bytes(), &list)
		if err != nil {
			t.Error("Could not unmarshal JSON")
		}

		utils.AssertStringEqualsTo(t, list.Name, genListForDeleteName)
	})

	t.Run("Delete a non-existing task with valid ObjectId", func(t *testing.T) {
		listURL := getlistURL(boardId, bson.NewObjectId())

		req, _ := http.NewRequest("DELETE", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Delete a task with Invalid Object ID", func(t *testing.T) {
		listURL := getInvalidlistURL(boardId)

		req, _ := http.NewRequest("DELETE", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}
