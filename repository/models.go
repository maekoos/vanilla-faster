// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql"
)

type Category struct {
	Slug         string         `json:"slug"`
	Name         string         `json:"name"`
	CollectionID int32          `json:"collectionId"`
	ImageUrl     sql.NullString `json:"imageUrl"`
}

type Collection struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Product struct {
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	Price           string         `json:"price"`
	SubcategorySlug string         `json:"subcategorySlug"`
	ImageUrl        sql.NullString `json:"imageUrl"`
}

type Subcategory struct {
	Slug            string         `json:"slug"`
	Name            string         `json:"name"`
	SubcollectionID int32          `json:"subcollectionId"`
	ImageUrl        sql.NullString `json:"imageUrl"`
}

type Subcollection struct {
	ID           int32  `json:"id"`
	Name         string `json:"name"`
	CategorySlug string `json:"categorySlug"`
}
