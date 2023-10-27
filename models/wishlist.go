package models

type Wishlist struct {
	Wishlist_ID uint `json:"whishlistid" gorm:"primaryKey;unique"`
	Product_ID  uint `json:"productid" gorm:"not null"`
	User_ID     uint `json:"userid" gorm:"not null"`
}
