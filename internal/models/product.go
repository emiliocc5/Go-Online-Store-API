package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	CategoryId  int     `json:"categoryId"`
	Label       string  `json:"label"`
	Type        int     `json:"type"`
	DownloadUrl string  `json:"downloadUrl"`
	Weight      float64 `json:"weight"`
}
