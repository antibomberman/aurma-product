package elasticModels

import "aurma_product/internal/models"

type Product struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Price       int    `json:"price" db:"price"`
	Slug        string `json:"slug" db:"slug"`
	Count       int    `json:"count" db:"count"`
	IsActive    bool   `json:"is_active" db:"is_active"`
	CompanyName string `json:"company_name" db:"company_name"`
	Barcode     string `json:"barcode" db:"barcode"`
	Mnn         string `json:"mnn" db:"mnn"`
	IssueForm   string `json:"issue_form" db:"issue_form"`
}

func (p *Product) ToProductDetail(images []models.ProductImage) models.ProductDetail {
	return models.ProductDetail{
		Id:          p.Id,
		Title:       p.Title,
		Price:       p.Price,
		Slug:        p.Slug,
		Count:       p.Count,
		IsActive:    p.IsActive,
		CompanyName: p.CompanyName,
		Barcode:     p.Barcode,
		Mnn:         p.Mnn,
		IssueForm:   p.IssueForm,
		Images:      images,
	}
}
