package controllers

import (
	"server/pkg/models"
	"strings"
)

func recommend(clusters []models.Cluster, userLoc string, cheapest int) []models.Cluster {
	final := []models.Cluster{}

	//iterate thru clusters
	for _, cluster := range clusters {
		score := calQty(len(cluster.Products)) + calPrice(cheapest, cluster.TotalPrice) + calLoc(userLoc, cluster.Products[0].ShopInfo.City)
		cluster.Score = score
		final = append(final, cluster)
	}

	return final
}

func calQty(qty int) float32 {
	qtyScore := 3

	return float32(qty * qtyScore)
}

func calPrice(cheapest int, price int) float32 {
	var m float32
	priceScore := 2
	m = float32(priceScore) * float32(cheapest) / float32(price)

	return m
}

func calLoc(userLoc string, loc string) float32 {
	if strings.Contains(strings.ToLower(loc), strings.ToLower(userLoc)) {
		return 1
	} else {
		return 0
	}
}
