package controllers

import (
	"server/pkg/models"
)

func kelompokin(products []models.Product) ([]models.Cluster, int) {

	// products := []models.Product{
	// 	{Name: "Produk A", Price: 300, ShopInfo: models.Shop{Id: "113", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "kemeja"}},
	// 	{Name: "Produk B", Price: 100, ShopInfo: models.Shop{Id: "452", Name: "kok", City: "Kota Surabaya"}, Key: models.Key{Marketplace: "Tokopedia", Barang: "kemeja"}},
	// 	{Name: "Produk C", Price: 400, ShopInfo: models.Shop{Id: "113", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "blouse"}},
	// 	{Name: "Produk C", Price: 200, ShopInfo: models.Shop{Id: "113", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "blouse"}},
	// 	{Name: "Produk D", Price: 140, ShopInfo: models.Shop{Id: "452", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "celana"}},
	// 	{Name: "Produk E", Price: 200, ShopInfo: models.Shop{Id: "452", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "kemeja"}},
	// 	{Name: "Produk F", Price: 400, ShopInfo: models.Shop{Id: "245", City: "Kota Bandung"}, Key: models.Key{Marketplace: "Shopee", Barang: "kemeja"}},
	// }

	// Group products by market and shop
	groups := make(map[string]map[string][]models.Product)
	for _, product := range products {
		market := product.Key.Marketplace
		shop := product.ShopInfo.Id

		//tokopeida using name as id
		if market == "Tokopedia" {
			shop = product.ShopInfo.Name
		}

		if _, ok := groups[market]; !ok {
			groups[market] = make(map[string][]models.Product)
		}
		if _, ok := groups[market][shop]; !ok {
			groups[market][shop] = []models.Product{}
		}

		//assign product to its own shop and market
		groups[market][shop] = append(groups[market][shop], product)
	}

	//make clusters
	clusters := []models.Cluster{}
	cheapest := 10000000000
	for category, categoryGroups := range groups {
		//iterate every shop to differentiate unique products
		for shop, products := range categoryGroups {
			//only working on shop that have >=2 product
			if len(products) >= 2 {
				uniqueCategories := make(map[string]bool)
				uniqueProducts := []models.Product{}
				altProducts := []models.Product{}
				price := 0

				for _, product := range products {
					if _, exists := uniqueCategories[product.Key.Barang]; !exists {
						price += product.Price
						uniqueCategories[product.Key.Barang] = true
						uniqueProducts = append(uniqueProducts, product)
					} else {
						altProducts = append(altProducts, product)
					}

				}

				if price < cheapest {
					cheapest = price
				}

				//cluster only accept >=2 products
				if len(uniqueProducts) >= 2 {
					clusters = append(clusters, models.Cluster{
						Shop:        shop,
						Products:    uniqueProducts,
						AltProducts: altProducts,
						Maretplace:  category,
						TotalPrice:  price,
					})
				}
			}
		}
	}

	return clusters, cheapest
}
