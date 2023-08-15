package models

type Cluster struct {
	Shop        string
	Products    []Product
	AltProducts []Product
	TotalPrice  int
	Maretplace  string
	Score       float32
}

type Rekomendasi struct {
	TotalToko int
	Toko      []Cluster
}
