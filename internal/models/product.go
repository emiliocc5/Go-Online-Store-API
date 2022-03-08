package models

type Product struct {
	Id          int     `gorm:"primarykey"`
	CategoryId  int     `json:"categoryId"`
	Label       string  `json:"label"`
	Type        int     `json:"type"`
	DownloadUrl string  `json:"downloadUrl"`
	Weight      float64 `json:"weight"`
}
