package response

type GetCartResponse struct {
	Products []ProductResponse `json:"products"`
}
