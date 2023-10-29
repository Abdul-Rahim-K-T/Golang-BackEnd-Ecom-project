package models

type Category_Offer struct{
	ID     			uint `json:"id" gorm:"primaryKey;unique"`
	Category_Id	uint`json:"categoryid" gorm:"not null"`
	Offer				bool`json:"offer" gorm:"not null"`
	Percentage	uint`json:"percentage"`
}