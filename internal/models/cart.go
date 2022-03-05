package dbModel

type Cart struct {
	Id       int `pg:"id"`
	ClientId int `pg:"clientid"`
}
