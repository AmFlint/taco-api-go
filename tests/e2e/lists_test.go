package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/models"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"gopkg.in/mgo.v2/bson"
	"github.com/AmFlint/taco-api-go/tests/utils/generator"
)

const (
	TESTING__LIST_NAME  = "Testing list"
	TESTING__LIST_ORDER = 1
	UPDATED__LIST_NAME = "Updated list"
	// data creation
	// delete endpoint
	genListForDeleteName = "about to be deleted"
	// view endpoint
	genListForViewName = "about to be viewed"
	// update endpoint
	genListForUpdateName = "to be updated"
)

// TODO: Create Helpers for Resource creations -> Tests run in parrallell which means reusing an id from above test may not work

func getListsBaseUrl(boardIdList bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/lists/", boardIdList.Hex())
}

func getlistURL(boardIdList, listId bson.ObjectId) string {
	return fmt.Sprintf("%s%s/", getListsBaseUrl(boardIdList), listId.Hex())
}

func getInvalidlistURL(boardIdList bson.ObjectId) string {
	return fmt.Sprintf("%s%s/", getListsBaseUrl(boardIdList), "2")
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

func getListForUpdate() *models.List {
	list := models.NewList()
	list.Name = genListForUpdateName
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

func getValidListUpdate() []byte {
	list := make(map[string]interface{})
	list["name"] = UPDATED__LIST_NAME
	return helpers.JsonEncode(list)
}

/* ------------------------------------------
   -------------- Test Suite ----------------
   ------------------------------------------ */

var (
	boardIdList bson.ObjectId
)

func init() {
	boardIdList = bson.NewObjectId()
}

// ---- Test Create Endpoint ---- //
func TestCreateListEndpoint(t *testing.T) {
	t.Run("Create List with valid informations", func(t *testing.T) {
		listURL := getListsBaseUrl(boardIdList)
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
		listURL := getListsBaseUrl(boardIdList)
		list := getInvalidListEmptyTitle()
		req, _ := http.NewRequest("POST", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Create List with Invalid Title (bad format)", func(t *testing.T) {
		listURL := getListsBaseUrl(boardIdList)
		list := getInvalidListBadTitle()

		req, _ := http.NewRequest("POST", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Create List with empty request body", func(t *testing.T) {
		listURL := getListsBaseUrl(boardIdList)
		req, _ := http.NewRequest("POST", listURL, nil)
		response := utils.ExecuteRequest(req)
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

// ---- Test View Endpoint ---- //
func TestViewListHandler(t *testing.T) {
	testedListID := generator.GenerateListAndGetID(t, getListForView())

	t.Run("View an existing list with valid ID", func(t *testing.T) {
		listURL := getlistURL(boardIdList, testedListID)

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
			listURL := getlistURL(boardIdList, bson.NewObjectId())

			req, _ := http.NewRequest("GET", listURL, nil)
			response := utils.ExecuteRequest(req)

			utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("View a list with invalid ID", func(t *testing.T) {
		listURL := getInvalidlistURL(boardIdList)

		req, _ := http.NewRequest("GET", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

// ---- Test Delete EndPoint ---- //
func TestDeleteListHandler(t *testing.T) {
	testedListID := generator.GenerateListAndGetID(t, getListForDelete())

	t.Run("Delete an Existing task with valid Id", func(t *testing.T) {
		listURL := getlistURL(boardIdList, testedListID)

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
		listURL := getlistURL(boardIdList, bson.NewObjectId())

		req, _ := http.NewRequest("DELETE", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Delete a task with Invalid Object ID", func(t *testing.T) {
		listURL := getInvalidlistURL(boardIdList)

		req, _ := http.NewRequest("DELETE", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestUpdateListHandler(t *testing.T) {
	testedListID := generator.GenerateListAndGetID(t, getListForUpdate())

	t.Run("Update a list with valid informations", func(t *testing.T) {
		listURL := getlistURL(boardIdList, testedListID)
		list := getValidListUpdate()

		req, _ := http.NewRequest("PATCH", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusOK)

		var updatedList models.List
		if err := json.Unmarshal(response.Body.Bytes(), &updatedList); err != nil {
			t.Error("[Error] in Update List endpoint, could not unmarshal response body")
		}

		utils.AssertStringEqualsTo(t, updatedList.Name, UPDATED__LIST_NAME)
	})

	t.Run("Update a list with invalid informations (Bad format title)", func(t *testing.T) {
		listURL := getlistURL(boardIdList, testedListID)
		list := getInvalidListBadTitle()

		req, _ := http.NewRequest("PATCH", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Update non existing List", func(t *testing.T) {
		listURL := getlistURL(boardIdList, bson.NewObjectId())
		list := getValidListUpdate()

		req, _ := http.NewRequest("PATCH", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusNotFound)
	})

	t.Run("Update list with empty request body", func(t *testing.T) {
		listURL := getlistURL(boardIdList, testedListID)

		req, _ := http.NewRequest("PATCH", listURL, nil)
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Update list with Invalid Object ID", func(t *testing.T) {
		listURL := getInvalidlistURL(boardIdList)
		list := getValidListUpdate()

		req, _ := http.NewRequest("PATCH", listURL, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}
