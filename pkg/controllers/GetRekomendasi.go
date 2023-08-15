package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/pkg/models"
	"strings"
)

type ArrayBarang struct {
	Barang   []string
	Location string
}

type Products struct {
	Name        string
	Keyword     string
	Marketplace string
	Product     []models.Product
}

type Items struct {
	Name        string
	key         string
	currentPage int
}

func GetRekomendasi(w http.ResponseWriter, r *http.Request) {
	// sort := r.URL.Query().Get("sortBy")
	initPage := 1
	items := []Items{}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	//store barang in array
	var daftar ArrayBarang
	if err = json.Unmarshal(body, &daftar); err != nil {
		log.Fatalln(err)
	}

	if validasi(daftar.Barang) {
		// make loop for adding "%20" to each barang to make  request
		for i := 0; i < len(daftar.Barang); i++ {
			// t := strings.Split(daftar.Barang[i], " ")
			// p := strings.Join(t, "%20")

			//store initial data user items user in "items"
			items = append(items, Items{
				Name:        daftar.Barang[i],
				key:         daftar.Barang[i],
				currentPage: 1,
			})
		}

		// request for 1st iterate
		products := requestProductAll(items, initPage)

		//clustering
		clusters, cheapest := kelompokin(products)

		//get score
		clusters = recommend(clusters, daftar.Location, cheapest)

		//check if cluster empty
		clusters = isEmpty(clusters, products, items, initPage, daftar.Location)

		if len(clusters) == 0 {
			// send 400 created response
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Tidak ditemukan toko")
		} else {
			clusters := QuickSortDesc(clusters)
			clusters = getAllShopeeName(clusters)

			//send response
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.Rekomendasi{
				TotalToko: len(clusters),
				Toko:      clusters,
			})
		}
	} else {
		// send 400 created response
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Jumlah barang tidak boleh lebih dari 5 dan kurang dari sama dengan 1")
	}
}

func validasi(barang []string) bool {
	if len(barang) > 5 || len(barang) <= 1 {
		return false
	} else {
		return true
	}
}

func requestProductAll(items []Items, page int) []models.Product {
	products := []models.Product{}
	for _, item := range items {
		//make request to Shopee and Tokopedia for each barang
		products = append(products, GetProductTokopedia(item.key, item.Name)...)
		products = append(products, GetProductShopee(item.key, item.Name, page)...)
	}
	return products
}

//check if clusters empty
func isEmpty(clusters []models.Cluster, products []models.Product, items []Items, page int, loc string) []models.Cluster {
	if len(clusters) == 0 {
		cheapest := 0
		//need to give limit here...

		//if cluster is empty => call get another page of products
		products := []models.Product{}
		for _, item := range items {
			//make request to Shopee and Tokopedia for each barang
			products = append(products, GetProductShopee(item.key, item.Name, page)...)
		}

		//clustering...
		clusters, cheapest = kelompokin(products)

		//recommend
		clusters = recommend(clusters, loc, cheapest)

		//call isEmpty
		clusters = isEmpty(clusters, products, items, page+1, loc)
	} else {
		complete := false
		for _, cluster := range clusters {
			//if found cluster tht have all products => send request
			if len(cluster.Products) == len(items) {
				complete = true
			}
		}

		if !complete {
			//rank cluster that have least incomplete products
			rank := []models.Rank{}
			for p, cluster := range clusters {
				rank = append(rank, models.Rank{
					ShopId:          cluster.Shop,
					Place:           p,
					IncompleteItems: foundItems(cluster.Products, items),
					Marketplace:     cluster.Maretplace,
				})
			}

			//ranking of cluster that have incomplete products
			ranking := QuickSortRank(rank)

			//visit each shop and find incomplete products
			for _, rank := range ranking {
				moreProducts := []models.Product{}
				//get more products
				if rank.Marketplace == "Shopee" {
					moreProducts = GetShopeeByShop(rank.ShopId)
					fmt.Println("moreee", moreProducts)
				}

				//find incomplete products
				for _, prod := range moreProducts {

					//looping thru incomplete items
					for _, item := range rank.IncompleteItems {
						if strings.Contains(strings.ToLower(prod.Name), strings.ToLower(item.Barang)) {
							//append to cluster
							clusters[rank.Place].Products = append(clusters[rank.Place].Products, prod)

							//found all? send response
							if len(clusters[rank.Place].Products) == len(items) {
								complete = true
								break
							} else {
								//append products if the products is apart of categories
								products = append(products, models.Product{
									Key: models.Key{
										Barang: item.Barang,
									},
								})
							}

						}
					}
					if complete {
						break
					}
				}
				if complete {
					break
				}
				//not found visit another shop below current rank
			}

		}
	}
	return clusters
}

func foundItems(product []models.Product, items []Items) []models.Key {
	uniqueProducts := []models.Key{}

	for _, item := range items {
		found := false
		for _, existingProduct := range product {
			if item.Name == existingProduct.Key.Barang {
				found = true
				break
			}
		}

		if !found {
			uniqueProducts = append(uniqueProducts, models.Key{
				Barang:  item.Name,
				Keyword: item.key,
			})
		}
	}
	return uniqueProducts
}

func getAllShopeeName(clusters []models.Cluster) []models.Cluster {
	data := []models.Cluster{}
	for _, cluster := range clusters {
		if cluster.Maretplace == "Shopee" {
			shop, url := getDetailShop(cluster.Shop)
			cluster.Shop = shop
			cluster.Products[0].ShopInfo.Url = "https://shopee.co.id/" + url
		}

		data = append(data, cluster)
	}
	return data
}
