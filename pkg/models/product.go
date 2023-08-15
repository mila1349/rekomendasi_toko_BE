package models

type Product struct {
	Id                 int
	Name               string
	Url                string
	Img                string
	DiscountPercentage int
	OriginalPrice      int
	Price              int
	Rating             int
	ShopInfo           Shop
	Key                Key
}

type ProductThen struct {
	Id                 int
	Name               string
	Url                string
	Img                string
	DiscountPercentage int
	OriginalPrice      int
	Price              int
	Rating             int
	ShopInfo           Shop
}

type Shop struct {
	Id   string
	Name string
	Url  string
	City string
}

type Key struct {
	Marketplace string
	Keyword     string
	Barang      string
}

type Rank struct {
	ShopId          string
	Place           int
	Marketplace     string
	IncompleteItems []Key
}
