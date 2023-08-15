package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/pkg/models"
	"strconv"
	"strings"
)

type TokpedInputKey struct {
	Query  string `json:"query"`
	Result []struct {
		Title    string `json:"title"`
		Link     string `json:"link"`
		Slug     string `json:"slug"`
		Image    string `json:"image"`
		Price    string `json:"price"`
		Location string `json:"location"`
		Rating   int    `json:"rating"`
		Sold     string `json:"sold"`
	} `json:"results"`
}

type TokpedInputShop struct {
	Data struct {
		ProductPage  string `json:"productPage"`
		LimitMaxPage string `json:"limitMaxPage"`
		Data         []struct {
			Name        string `json:"name"`
			Product_url string `json:"product_url"`
			Product_id  string `json:"product_id"`
			Price       struct {
				Text string `json:"text_idr"`
			} `json:"price"`
			Images struct {
				Original string `json:"original"`
			} `json:"primary_image"`
			Campaign struct {
				Discounted_percentage string `json:"discounted_percentage"`
			} `json:"campaign"`
		} `json:"data"`
	} `json:"data"`
}

type Limit struct {
	Limit  int
	Page   int
	Length int
}

//function to get shopee products
func GetProductTokopedia(key string, barang string) []models.Product {
	products := []models.Product{}
	market := models.Key{
		Keyword:     key,
		Barang:      barang,
		Marketplace: "Tokopedia",
	}
	// url := "https://tokopediaapi.p.rapidapi.com/" + key + "/" + strconv.Itoa(page) + "/60"
	url := "https://tokopediaapi.p.rapidapi.com/?act=search&query=" + key

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "961dfb0036mshbd971a9502c707bp184958jsn0fd5d2ce178b")
	req.Header.Add("X-RapidAPI-Host", "tokopediaapi.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	var data TokpedInputKey
	json.Unmarshal([]byte(body), &data)
	for i, product := range data.Result {
		//convert price (String) to int ex= "Rp.900.000"
		priceFirst := strings.Split(product.Price, "Rp")[1]
		priceSecond := strings.Split(priceFirst, ".")
		priceResult := strings.Join(priceSecond, "")
		price, err := strconv.Atoi(priceResult)

		toko := strings.Split(strings.Split(product.Link, "https://www.tokopedia.com/")[1], "/")[0]

		if err != nil {
			fmt.Println("Error conversion price")
		}
		products = append(products, models.Product{
			Id:     i,
			Name:   product.Title,
			Url:    product.Link,
			Img:    product.Image,
			Price:  price,
			Rating: product.Rating,
			ShopInfo: models.Shop{
				Name: toko,
				Url:  "https://www.tokopedia.com/" + toko,
				City: product.Location,
			},
			Key: market,
		})
	}

	return products
}

//if limitMaxPage==100 then get another page
func GetTokpedByShop(id string) []models.Product {
	products := []models.Product{}
	limit := Limit{
		Page: 0,
	}

	products = append(products, getTokpedShop(id, products, limit)...)

	return products
}

func getTokpedShop(id string, products []models.Product, limit Limit) []models.Product {
	url := "https://tokopedia-api.p.rapidapi.com/shop/products/" + id + "/" + strconv.Itoa(limit.Page+1)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "5d92f21dc1mshea549cfc3d34a23p185438jsn0971f909d143")
	req.Header.Add("X-RapidAPI-Host", "tokopedia-api.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	var data TokpedInputShop
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	//store limit for further use
	lim, err := strconv.Atoi(strings.Split(data.Data.LimitMaxPage, " ")[0])
	limit.Limit = lim
	limit.Length = len(data.Data.Data)
	limit.Page = limit.Page + 1

	for _, product := range data.Data.Data {
		//convert price (String) to int ex= "Rp 159.000"
		priceFirst := strings.Split(product.Price.Text, " ")[1]
		priceSecond := strings.Split(priceFirst, ".")
		priceResult := strings.Join(priceSecond, "")
		price, err := strconv.Atoi(priceResult)

		if err != nil {
			fmt.Println("Error conversion price")
		}
		id, _ := strconv.Atoi(product.Product_id)
		sale, _ := strconv.Atoi(product.Campaign.Discounted_percentage)
		products = append(products, models.Product{
			Id:                 id,
			Name:               product.Name,
			Url:                product.Product_url,
			Img:                product.Images.Original,
			DiscountPercentage: sale,
			Price:              price,
			ShopInfo: models.Shop{
				Id: strconv.Itoa(id),
			},
		})
	}

	// get more products if there are more
	if limit.Length == limit.Limit {
		products = append(products, getTokpedShop(id, products, limit)...)
	}

	return products

}
