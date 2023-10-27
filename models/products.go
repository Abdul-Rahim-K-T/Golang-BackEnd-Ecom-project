package models

import "gorm.io/gorm"

type Category struct {
	
	Category_ID		uint		`json:"categoryid" gorm:"primaryKey;autoIncrement"`
	Category_Name string `json:"categoryname" gorm:"not null"`
	Unlist        bool   `json:"unlist" gorm:"default:false"`
}

type Product struct {
	gorm.Model
	Product_Name string `json:"productname" gorm:"not null"`
	Description  string `json:"description" gorm:"not null"`
	Stock        uint   `json:"stock" gorm:"not null"`
	Price        uint   `json:"price" gorm:"not null"`
	Category_Id  uint   `json:"categoryid" gorm:"not null"`
	Brand_ID     uint   `json:"brandid" gorm:"not null"`

	// Define the Category relationship with foreign key
	// Category Category `gorm:"foreignKey:Category_Id"`

	// Define the Brand relationship with foreign key
	// Brand Brand `gorm:"foreignKey:Brand_ID"`

	// Define the Images relationship with foreign key
	// Images []Image `gorm:"foreignKey:ProductID"`
}

type Brand struct {
	
	Brand_ID   uint   `json:"brandid"  gorm:"primaryKey;unique"`
	Brand_Name string `json:"brandname" gorm:"not null"`
}
type Image struct {
	Id 				uint		`json"id" gorm:"primaryKey;unique"`
	Image     string `json:"image" gorm:"not null"`
	ProductID uint   `json:"product_id" gorm:"not null"`

	// Define the Product relationship
	// Product Product `gorm:"foreignKey:ProductID"`
}
