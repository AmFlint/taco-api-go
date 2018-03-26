package models

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/AmFlint/taco-api-go/validators"
)

// Task Structure, represents Task Document from Database
type Task struct {
	TaskId      bson.ObjectId `bson:"_id" json:"taskId"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Status      string        `bson:"status" json:"status"`
	Points      float64           `bson:"points" json:"points"`
}

const (
	TASK__VALIDATOR__INVALID_TYPE_POINTS = "Invalid type for field 'points', expected to be a number"
	TASK__VALIDATOR__INVALID_EMPTY_TITLE = "Invalid value for field 'title': Empty field"
	TASK__STATUS__IN_PROGRESS 			 = "in progress"
	TASK__STATUS__DONE 			 		 = "done"
	TASK__TITLE = "title"
	TASK__DESCRIPTION = "description"
	TASK__STATUS = "status"
	TASK__POINTS = "points"
)

// Set Default Status to a Task Entity
func (t *Task) SetDefaultStatus() {
	t.Status = TASK__STATUS__IN_PROGRESS
}

func GetDefaultTaskStatus() string {
	return TASK__STATUS__IN_PROGRESS
}

// Retrieve Valid Values for a Task's Status field
func GetValidTaskStatuses() []string {
	return []string{
		TASK__STATUS__IN_PROGRESS,
		TASK__STATUS__DONE,
	}
}

// Validation Rules/Methods for a Task Title field
func  validateTitle(errors *[]string, title interface{}) {
	// Check if field is empty
	emptyErr := validators.NotEmpty(title, TASK__TITLE)
	validators.CheckValidationErrors(emptyErr, errors)
	// Leave validation if field is empty, continue would cause a panic
	if emptyErr != nil {
		return
	}

	// Checking if given parameter is of type string
	typeErr := validators.IsString(title, TASK__TITLE)
	validators.CheckValidationErrors(typeErr, errors)

	// Title is not a string, return before causing a panic during next validations
	if typeErr != nil {
		return
	}
	// Check if title length is greater than 200 characters
	maxErr := validators.Maxlen(title.(string), 200, TASK__TITLE)
	// Check if title length is lesser than 4 characters
	minErr := validators.Minlen(title.(string), 4, TASK__TITLE)

	// If validation variables contain errors, append them to errors slice reference
	validators.CheckValidationErrors(maxErr, errors)
	validators.CheckValidationErrors(minErr, errors)
}

// Validation Rules/Methods for a Task Description field
func validateDescription(errors *[]string, d interface{}) {
	if validators.NotEmpty(d, TASK__DESCRIPTION) != nil {
		return
	}
	// Check that Description parameter is of type String
	typeErr := validators.IsString(d, TASK__DESCRIPTION)
	validators.CheckValidationErrors(typeErr, errors)
	// If type is not string, return to prevent panic from next validations
	if typeErr != nil {
		return
	}

	// Check that description doesn't exceed 500 characters
	maxErr := validators.Maxlen(d.(string), 500, TASK__DESCRIPTION)
	// If validation variables contain errors, append them to errors slice reference
	validators.CheckValidationErrors(maxErr, errors)
}

// Validation Rules/Methods for a Task Points field
func validatePoints(errors *[]string, p interface{}) {
	// Skip validation if field "points" not provided = optionional
	if validators.NotEmpty(p, TASK__DESCRIPTION) != nil {
		return
	}

	typeErr := validators.IsFloat(p, TASK__POINTS)
	validators.CheckValidationErrors(typeErr, errors)
	// Leave function if points is not a float as next validations would cause a panic
	if typeErr != nil {
		return
	}

	maxErr := validators.MaxF(p.(float64), 100, TASK__POINTS)
	minErr := validators.MinF(p.(float64), 0, TASK__POINTS)

	// If validation variables contain errors, append them to errors slice reference
	validators.CheckValidationErrors(maxErr, errors)
	validators.CheckValidationErrors(minErr, errors)
}

// Validation Rules/Methods for a Task Status field
func validateStatus(errors *[]string, s interface{}) {
	// Check that Description parameter is of type String
	typeErr := validators.IsString(s, TASK__STATUS)
	validators.CheckValidationErrors(typeErr, errors)
	// If type is not string, return to prevent panic from next validations
	if typeErr != nil {
		return
	}

	err := validators.In(s.(string), GetValidTaskStatuses(), TASK__STATUS)
	validators.CheckValidationErrors(err, errors)
}

// Validate Task fields, return array of error messages
func ValidateTask(task map[string]interface{}) []string {
	var e []string

	// Validate each
	validatePoints(&e, task[TASK__POINTS])
	validateDescription(&e, task[TASK__DESCRIPTION])
	validateTitle(&e, task[TASK__TITLE])
	validateStatus(&e, task[TASK__STATUS])

	return e
}
