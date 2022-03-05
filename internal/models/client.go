package models

type Client struct {
	tableName struct{} `pg:"clients"`
	Id        int      `json:"id"`
}
