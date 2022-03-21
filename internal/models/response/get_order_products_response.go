package response

type GetOrderProductsResponse struct {
	Products []ProductResponse `json:"products"`
}
