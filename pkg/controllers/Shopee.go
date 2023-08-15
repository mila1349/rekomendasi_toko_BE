package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/pkg/models"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ShopeeDataKey struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Items []struct {
			Id        int    `json:"item_id"`
			ShopId    int    `json:"shop_id"`
			Name      string `json:"title"`
			Image     string `json:"img"`
			PriceInfo struct {
				Price string `json:"price"`
			} `json:"price_info"`
			ShopInfo struct {
				City string `json:"shop_location"`
			} `json:"shop_info"`
			DiscountPercentage int `json:"show_discount"`
			ItemRatings        struct {
				Rating_star int `json:"rating_star"`
			} `json:"item_rating"`
		} `json:"items"`
	} `json:"data"`
}

type ShopeeDataShop struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Page        int `json:"page"`
		Page_size   int `json:"page_size"`
		Total_count int `json:"total_count"`
		Items       []struct {
			Id        int    `json:"item_id"`
			ShopId    int    `json:"shop_id"`
			Name      string `json:"title"`
			Image     string `json:"img"`
			PriceInfo struct {
				Price string `json:"price"`
			} `json:"price_info"`
			ShopInfo struct {
				City string `json:"shop_location"`
			} `json:"shop_info"`
			DiscountPercentage int `json:"show_discount"`
			ItemRatings        struct {
				Rating_star int `json:"rating_star"`
			} `json:"item_rating"`
		} `json:"items"`
	} `json:"data"`
}

type ShopeeShopKey struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ShopName string `json:"shop_name"`
		UserName string `json:"username"`
	} `json:"data"`
}

//function to get shopee products
func GetProductShopee(key string, barang string, page int) []models.Product {
	godotenv.Load(".env")

	products := []models.Product{}
	market := models.Key{
		Keyword:     key,
		Barang:      barang,
		Marketplace: "Shopee",
	}

	// url := "https://shopee-e-commerce-data.p.rapidapi.com/api/sc/shopee/search/items?site=id&keyword=" + key + "&page=" + strconv.Itoa(page) + "&pageSize=60"

	url := "https://shopee-e-commerce-data.p.rapidapi.com/shopee/search/items/v2?site=id&keyword=" + key + "&by=relevancy&page=" + strconv.Itoa(page) + "&pageSize=60"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "961dfb0036mshbd971a9502c707bp184958jsn0fd5d2ce178b")
	req.Header.Add("X-RapidAPI-Host", "shopee-e-commerce-data.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	var data ShopeeDataKey
	json.Unmarshal([]byte(body), &data)
	for _, product := range data.Data.Items {
		price, err := strconv.Atoi(strings.Split(product.PriceInfo.Price, ".0")[0])

		if err != nil {
			fmt.Println("Error conversion price")
		}

		products = append(products, models.Product{
			Id:                 product.Id,
			Name:               product.Name,
			Img:                product.Image,
			DiscountPercentage: product.DiscountPercentage,
			Price:              price,
			ShopInfo: models.Shop{
				Id:   strconv.Itoa(product.ShopId),
				City: product.ShopInfo.City,
			},
			Key: market,
		})
	}
	return products
}

//if limitMaxPage==100 then get another page
func GetShopeeByShop(id string) []models.Product {
	products := []models.Product{}
	limit := Limit{
		Page:  0,
		Limit: 100,
	}

	products = append(products, getShopeeShop(id, products, limit)...)

	return products
}

func getShopeeShop(id string, products []models.Product, limit Limit) []models.Product {
	// url := "https://shopee-e-commerce-data.p.rapidapi.com/api/sc/shopee/shop/items?shop_id=" + id +
	// 	"&site=id&page=" + strconv.Itoa(limit.Page) + "&pageSize=" + strconv.Itoa(limit.Limit)

	url := "https://shopee-e-commerce-data.p.rapidapi.com/shopee/shop/items/v2?site=id&shop_id=" + id + "&page=" + strconv.Itoa(limit.Page) + "&pageSize=" + strconv.Itoa(limit.Limit) + "&by=relevancy"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "961dfb0036mshbd971a9502c707bp184958jsn0fd5d2ce178b")
	req.Header.Add("X-RapidAPI-Host", "shopee-e-commerce-data.p.rapidapi.com")
	fmt.Println("iddd", id)

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	var data ShopeeDataShop
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	fmt.Println("datt", data)

	//store limit for further use
	limit.Length = data.Data.Total_count
	limit.Page = limit.Page + 1
	limit.Limit = data.Data.Page_size

	for _, product := range data.Data.Items {
		price, err := strconv.Atoi(strings.Split(product.PriceInfo.Price, ".0")[0])

		if err != nil {
			fmt.Println("Error conversion price")
		}

		products = append(products, models.Product{
			Id:                 product.Id,
			Name:               product.Name,
			Img:                product.Image,
			DiscountPercentage: product.DiscountPercentage,
			Price:              price,
			ShopInfo: models.Shop{
				Id: strconv.Itoa(product.ShopId),
			},
		})
	}

	// get more products if there are more
	if limit.Length == limit.Limit {
		products = append(products, getShopeeShop(id, products, limit)...)
	}

	return products
}

func getDetailShop(id string) (string, string) {
	url := "https://shopee-e-commerce-data.p.rapidapi.com/shopee/shop/shop_info?site=id&shop_id=" + id
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "961dfb0036mshbd971a9502c707bp184958jsn0fd5d2ce178b")
	req.Header.Add("X-RapidAPI-Host", "shopee-e-commerce-data.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	var data ShopeeShopKey
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	return data.Data.ShopName, data.Data.UserName
}
