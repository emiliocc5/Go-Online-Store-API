package response

type ProductResponse struct {
	Id          int     `json:"id"`
	Label       string  `json:"label"`
	Category    string  `json:"category"`
	Type        int     `json:"type"`
	DownloadUrl string  `json:"downloadUrl,omitempty"`
	Weight      float64 `json:"weight,omitempty"`
}
