package response

import (
	"github.com/emiliocc5/online-store-api/internal/models"
)

type GetCartResponse struct {
	Products []models.Product `json:"products"`
}
