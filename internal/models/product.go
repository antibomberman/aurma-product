package models

import (
	"database/sql"
	pb "github.com/antibomberman/aurma-protos/gen/go/product"
)

type Product struct {
	Id       int          `db:"id"`
	Title    string       `db:"title"`
	Slug     string       `db:"slug"`
	IsActive sql.NullBool `db:"is_active"`
}
type ProductSearchWithData struct {
	Id          int            `json:"id" db:"id"`
	Title       string         `json:"title" db:"title"`
	Slug        string         `json:"slug" db:"slug"`
	IsActive    sql.NullBool   `json:"is_active" db:"is_active"`
	CompanyName sql.NullString `json:"company_name" db:"company_name"`
	Barcode     sql.NullString `json:"barcode" db:"barcode"`
	Mnn         sql.NullString `json:"mnn" db:"mnn"`
	IssueForm   sql.NullString `json:"issue_form" db:"issue_form"`
}

type ProductDetail struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Price       int            `json:"price"`
	Slug        string         `json:"slug"`
	Count       int            `json:"count"`
	IsActive    bool           `json:"is_active"`
	CompanyName string         `json:"company_name"`
	Barcode     string         `json:"barcode" `
	Mnn         string         `json:"mnn" `
	IssueForm   string         `json:"issue_form" `
	Images      []ProductImage `json:"images"`
}

// ProductElastic is the struct for product_list index in Elasticsearch
type ProductElastic struct {
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

func (p *ProductElastic) ToProductDetail(images []ProductImage) ProductDetail {
	return ProductDetail{
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
func (p *ProductSearchWithData) ToProductElastic(price, count int) ProductElastic {
	return ProductElastic{
		Id:          p.Id,
		Title:       p.Title,
		Price:       price,
		Slug:        p.Slug,
		Count:       count,
		IsActive:    p.IsActive.Bool,
		CompanyName: p.CompanyName.String,
		Barcode:     p.Barcode.String,
		Mnn:         p.Mnn.String,
		IssueForm:   p.IssueForm.String,
	}
}
func (p *ProductDetail) ToPbProduct() *pb.Product {
	images := make([]*pb.ProductImage, 0, len(p.Images))
	if len(p.Images) > 0 {
		for _, img := range p.Images {
			images = append(images, img.ToPbImages())
		}
	}

	return &pb.Product{
		Id:          int32(p.Id),
		Title:       p.Title,
		Price:       int32(p.Price),
		Slug:        p.Slug,
		Count:       int32(p.Count),
		IsActive:    p.IsActive,
		CompanyName: p.CompanyName,
		Barcode:     p.Barcode,
		Mnn:         p.Mnn,
		IssueForm:   p.IssueForm,
		Images:      images,
	}

}
