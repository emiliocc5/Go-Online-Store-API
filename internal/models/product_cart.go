package dbModel

type ProductCart struct {
	tableName struct{} `pg:"products_carts"`
	Id        int      `pg:"id"`
	ProductId int      `pg:"productid"`
	CartId    int      `pg:"cartid"`
}
