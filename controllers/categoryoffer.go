package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)
type categoryoffers struct {
	Category_Id uint `json:"categoryid"`
	Percentage  uint `json:"percentage"`
}

// AddCategoryOffer adds an offer to a category by ID.
// @Summary Add an offer to a category
// @Description Add an offer to a category by its ID
// @Tags CategoryOffers
// @Accept json
// @Produce json
// @Param categoryid body categoryoffers true "Category ID and Percentage for the offer"
// @Success 200 {json} AddCategoryOfferResponse "Successfully added an offer to a category"
// @Failure 400 {json} AddCategoryOfferResponse "Error in adding category offer"
// @Router /admin/categoryoffers [post]
func AddCategoryOffer(c *gin.Context) {
	
	var categoryoffer categoryoffers
	if err := c.Bind(&categoryoffer); err != nil {
		c.JSON(400, gin.H{
			"error": "Binding error",
		})
		return
	}

	var offer models.Category_Offer
	row := database.DB.Where("category_id = ? AND offer = true", categoryoffer.Category_Id).First(&offer).RowsAffected
	if row > 0 {
		c.JSON(400, gin.H{
			"error": "This offer already exist",
		})
		return
	}

	if categoryoffer.Percentage < 1 && categoryoffer.Percentage > 99 {
		c.JSON(400, gin.H{
			"error": "Invalid offer percentage",
		})
		return
	}
	var category models.Category
	err := database.DB.First(&category, categoryoffer.Category_Id).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This category doesn't exist in your database",
		})
		return
	}
	err = database.DB.Create(&models.Category_Offer{
		Category_Id: categoryoffer.Category_Id,
		Offer:       true,
		Percentage:  categoryoffer.Percentage,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully added category offer",
	})
}


// ListCatagoryOffer lists all category offers.
// @Summary List all category offers
// @Description List all category offers including category name, offer status, and percentage
// @Tags CategoryOffers
// @Accept json
// @Produce json
// // @Success 200 {json} ListCategoryOfferResponse "List of category offers"
// @Router /admin/categoryoffers [get]

func ListCatagoryOffer(c *gin.Context) {
	type categoryoffer struct {
		Category_Name string
		Offer         bool
		Percentage    uint
	}
	var offers []categoryoffer
	database.DB.Table("category_offers").Select("categories.category_name,category_offers.offer,category_offers.percentage").
		Joins("INNER JOIN categories ON categories.category_id=category_offers.category_id").
		Scan(&offers)

	c.JSON(200, gin.H{
		"message": offers,
	})
}


// DeleteCategoryOffer cancels an offer by ID.
// @Summary Cancel an offer by ID
// @Description Cancel an offer by its ID
// @Tags CategoryOffers
// @Accept json
// @Produce json
// @Param offer_id path int true "Offer ID to cancel"
/// @Success 200 {json} DeleteCategoryOfferResponse "Successfully cancelled an offer"
// @Failure 400 {json} DeleteCategoryOfferResponse "Error in cancelling offer"
// @Router /admin/categoryoffers/{offer_id} [put]
func DeleteCategoryOffer(c *gin.Context) {
	offid, err := strconv.Atoi(c.Param("offer_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find offerid",
		})
		return
	}

	var offer models.Category_Offer
	row := database.DB.First(&offer, offid).RowsAffected
	if row == 0 {
		c.JSON(400, gin.H{
			"error": "Failed to find offer in database",
		})
		return
	}
	if !offer.Offer {
		c.JSON(400, gin.H{
			"error": "This offer already cancelled",
		})
		return
	}
	err = database.DB.Model(&models.Category_Offer{}).Where("id=?", offid).Update("offer", false).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully cancelled offer",
	})
}