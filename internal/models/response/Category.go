package response

type Category struct {
	Id   int `gorm:"primarykey"`
	Name string
}
