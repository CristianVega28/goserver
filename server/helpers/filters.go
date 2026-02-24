package server_helpers

import "github.com/CristianVega28/goserver/core/models"

type (
	PaginationResponse struct {
		Page      int `json:"page"`
		TotalPage int `json:"total_page"`
		Data      any `json:"data"`
		Current   int `json:"current"`
		Previous  int `json:"previous"`
	}
)

const per_page = 10

func FilterPagination(page int, model models.Models[map[string]any]) (PaginationResponse, error) {
	var pagination PaginationResponse

	pagination.Page = page

	return pagination, nil
}
