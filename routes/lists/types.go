package lists

import (
	"github.com/AmFlint/taco-api-go/models"
)

type ListApiResponse struct {
	Lists []models.List `json:"lists"`
}