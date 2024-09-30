package models

type ProductPharmacy struct {
	ProductId  int `db:"product_id"  json:"product_id"`
	PharmacyId int `db:"pharmacy_id" json:"pharmacy_id"`
	Price      int `db:"price"       json:"price"`
	Count      int `db:"count"       json:"count"`
}
