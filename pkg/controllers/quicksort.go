package controllers

import (
	"server/pkg/models"
)

func QuickSortPrice(arr []models.Cluster) []models.Cluster {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[0]
	var less, greater []models.Cluster

	for _, num := range arr[1:] {
		if num.TotalPrice <= pivot.TotalPrice {
			less = append(less, num)
		} else {
			greater = append(greater, num)
		}
	}

	less = QuickSortPrice(less)
	greater = QuickSortPrice(greater)

	return append(append(less, pivot), greater...)
}

func QuickSortQty(arr []models.Cluster) []models.Cluster {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[0]
	var less, greater []models.Cluster

	for _, num := range arr[1:] {
		if len(num.Products) >= len(pivot.Products) {
			less = append(less, num)
		} else {
			greater = append(greater, num)
		}
	}

	less = QuickSortQty(less)
	greater = QuickSortQty(greater)

	return append(append(less, pivot), greater...)
}

func QuickSortRank(arr []models.Rank) []models.Rank {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[0]
	var less, greater []models.Rank

	for _, num := range arr[1:] {
		if len(num.IncompleteItems) <= len(pivot.IncompleteItems) {
			less = append(less, num)
		} else {
			greater = append(greater, num)
		}
	}

	less = QuickSortRank(less)
	greater = QuickSortRank(greater)

	return append(append(less, pivot), greater...)
}

func QuickSortDesc(arr []models.Cluster) []models.Cluster {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[len(arr)-1]
	var left, right []models.Cluster

	for i := 0; i < len(arr)-1; i++ {
		if arr[i].Score > pivot.Score {
			left = append(left, arr[i])
		} else {
			right = append(right, arr[i])
		}
	}

	sortedLeft := QuickSortDesc(left)
	sortedRight := QuickSortDesc(right)

	return append(append(sortedLeft, pivot), sortedRight...)
}
