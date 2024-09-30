package models

import (
	"database/sql"
	pb "github.com/antibomberman/aurma-protos/gen/go/product"
)

type ProductImage struct {
	Links struct {
		Medium    string `json:"medium"`
		Thumbnail string `json:"thumbnail"`
		Raw       string `json:"raw"`
		OgImage   string `json:"og_image"`
	} `json:"links"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (p *ProductImage) ToPbImages() *pb.ProductImage {
	return &pb.ProductImage{
		Links: &pb.ProductImageLinks{
			Medium:    p.Links.Medium,
			Thumbnail: p.Links.Thumbnail,
			Raw:       p.Links.Raw,
			OgImage:   p.Links.OgImage,
		},
		Title:       p.Title,
		Description: p.Description,
	}
}

type ProductGalleryImage struct {
	Id          int            `db:"id"`
	OwnerId     int            `db:"ownerId"`
	Name        sql.NullString `db:"name"`
	Description sql.NullString `db:"description"`
	Rank        int            `db:"rank"`
	Webp        sql.NullBool   `db:"webp"`
}
