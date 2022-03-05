package dbModel

type Product struct {
	tableName   struct{} `pg:"products"`
	Id          int      `json:"id" pg:"id"`
	CategoryId  int      `json:"categoryId" pg:"categoryid"`
	Label       string   `json:"label" pg:"label"`
	Type        int      `json:"type" pg:"type"`
	DownloadUrl string   `json:"downloadUrl" pg:"downloadurl"`
	Weight      float64  `json:"weight" pg:"weight"`
}
