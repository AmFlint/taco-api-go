package lists

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/AmFlint/taco-api-go/models"
	"testing"
	"net/http"
	"bytes"
	"github.com/AmFlint/taco-api-go/helpers"
	"github.com/AmFlint/taco-api-go/tests/utils"
	"encoding/json"
	"github.com/AmFlint/taco-api-go/tests/utils/testconfig"
)

const (
	TESTING__LIST_NAME = "Testing list"
	TESTING__LIST_ORDER = 1
)

func getListsBaseUrl(boardId bson.ObjectId) string {
	return fmt.Sprintf("/boards/%s/lists/", boardId.Hex())
}

func getListUrl(boardId, listId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s/", getListsBaseUrl(boardId), listId.Hex())
}

func getInvalidListUrl(boardId bson.ObjectId) string {
	return fmt.Sprintf("%s/%s/", getListsBaseUrl(boardId), "2")
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
	boardId, testedListId bson.ObjectId
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
		listUrl := getListsBaseUrl(boardId)
		list := getValidList()

		req, _ := http.NewRequest("POST", listUrl, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		var createdList models.List
		if err := json.Unmarshal(response.Body.Bytes(), &createdList); err != nil {
			t.Error("[Error] in Create List endpoint, could not unmarshal response body")
		}

		utils.CheckResponseCode(t, response.Code, http.StatusCreated)

		utils.AssertStringEqualsTo(t, createdList.Name, TESTING__LIST_NAME)
		// TODO: Implement "order" tests when board is implemented
		//utils.AssertIntEqualsTo(t, createdList.Order, 1)

		testedListId = createdList.ListId
	})

	t.Run("Create List with invalid user data - empty title", func(t *testing.T) {
		listUrl := getListsBaseUrl(boardId)
		list := getInvalidListEmptyTitle()
		req, _ := http.NewRequest("POST", listUrl, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)

		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("Create List with Invalid Title (bad format)", func(t *testing.T) {
		listUrl := getListsBaseUrl(boardId)
		list := getInvalidListBadTitle()

		req, _ := http.NewRequest("POST", listUrl, bytes.NewReader(list))
		response := utils.ExecuteRequest(req)
		utils.CheckResponseCode(t, response.Code, http.StatusBadRequest)
	})
}