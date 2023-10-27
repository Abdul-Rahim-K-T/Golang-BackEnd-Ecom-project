package controllers

import (
	
	
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)

type datas struct {
	Productid uint `json:"productid"`
}

// AddToWishtList allows the authenticated user to add a product to their wishlist.
// @Summary Add to Wishlist
// @Description Allows the authenticated user to add a product to their wishlist.
// @ID AddToWishlist
// @Accept json
// @Produce json
// @Param user body datas true "Product ID"
// @Success 200 {json} AddToWishlistResponse "Product successfully added to wishlist."
// @Failure 400 {json} ErrorResponse "Invalid input or product error"
// @Security Bearer
// @Router /user/addtowishlist [post]
func AddToWishtList(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	var data datas
	if err := c.Bind(&data); err != nil {
		c.JSON(500, gin.H{
			"error": "Binding error",
		})
		return
	}
	var wishlist models.Wishlist
	row := database.DB.Where("product_id=?", data.Productid).First(&wishlist).RowsAffected
	if row > 0 {
		c.JSON(400, gin.H{
			"error": "This product already in wishlist",
		})
		return
	}
	var product models.Product
	err := database.DB.First(&product, data.Productid).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This product doesnot exist in database",
		})
		return
	}
	err = database.DB.Create(&models.Wishlist{
		Product_ID: data.Productid,
		User_ID:    userId,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully " + product.Product_Name + " added to wishlist",
	})
}

// ListWishlist retrieves the authenticated user's wishlist items.
// @Summary List Wishlist
// @Description Retrieves the authenticated user's wishlist items.
// @ID ListWishlist
// @Produce json
// @Param page query int false "Page number"
// @Success 200 {json} ListWishListResponse "Wishlist items"
// @Failure 400 {json} ErrorResponse "Database error or no products found in wishlist"
// @Security Bearer
// @Router /user/listwishlist/{page} [get]
func ListWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, _ := strconv.Atoi(c.Query("page"))
	limit := 3
	offset := (page - 1) * limit

	type wishlist struct {
		Product_Name string
		Price        uint
	}
	var wishlists []wishlist
	err := database.DB.Table("wishlists").Select("products.product_name,products.price").
		Joins("INNER JOIN products ON products.id=wishlists.product_id").
		Where("user_id=?", userId).
		Limit(limit).Offset(offset).
		Scan(&wishlists).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"wishlist": wishlists,
	})
}

type carts1 struct {
	Wishlist_ID uint `json:"wishlistid"`
	Quantity    uint `json:"quantity"`
}

/// AddToCartFromWishlist adds a product from the user's wishlist to the cart.
// @Summary Add to Cart from Wishlist
// @Description Adds a product from the user's wishlist to the cart.
// @ID AddToCartFromWishlist
// @Accept json
// @Produce json
// @Success 200 {json} SuccessResponse "Products successfully added to the cart"
// @Failure 400 {json} ErrorResponse "Database or product error"
// @Security Bearer
// @Router /user/wishlist/addtocart [post]
func AddToCartFromWishlist(c *gin.Context) {
	// Extract user ID from the context
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	// Fetch all products in the user's wishlist
	var wishlistedProducts []models.Wishlist
	err := database.DB.Where("user_id = ?", userId).Find(&wishlistedProducts).Error
	if err != nil {
			c.JSON(500, gin.H{
					"error": "Database error",
			})
			return
	}

	if len(wishlistedProducts) == 0 {
			c.JSON(400, gin.H{
					"error": "No products found in wishlist",
			})
			return
	}

	// Check if any products are eligible for category offers
	var id []uint
	database.DB.Table("category_offers").Select("id").Where("offer = true").Scan(&id)
	eligibleCategories := make(map[uint]bool)
	for _, v := range id {
			eligibleCategories[v] = true
	}

	// Iterate through wishlisted products and add them to the cart
	for _, wishlistedProduct := range wishlistedProducts {
			// Retrieve the product details from the product table
			var product models.Product
			err = database.DB.First(&product, wishlistedProduct.Product_ID).Error
			if err != nil {
					c.JSON(400, gin.H{
							"error": "Failed to find product in the database",
					})
					return
			}

			// Calculate cart values based on category offers
			var cart models.Cart
			err = database.DB.Where("product_id = ? AND user_id = ?", wishlistedProduct.Product_ID, userId).First(&cart).Error
			if err != nil {
					cart = models.Cart{
							Product_ID:     wishlistedProduct.Product_ID,
							Quantity:       1, // You can set the initial quantity as needed.
							Price:          product.Price,
							Category_Id:    product.Category_Id,
							Total_Price:    product.Price,
							User_ID:        userId,
							Category_Offer: 0, // Initialize to 0.
					}
			} else {
					// Update the existing cart values
					var result models.Category_Offer
					database.DB.Where("category_id = ?", product.Category_Id).First(&result)
					cart.Quantity += 1 // You can update the quantity as needed.
					cart.Total_Price = (product.Price * cart.Quantity)
					if eligibleCategories[product.Category_Id] {
							cart.Total_Price -= (product.Price * cart.Quantity * result.Percentage / 100)
					}
					cart.Category_Offer = (product.Price * cart.Quantity * result.Percentage / 100)
			}

			// Create or update the cart record
			if err := database.DB.Save(&cart).Error; err != nil {
					c.JSON(400, gin.H{
							"error": err.Error(),
					})
					return
			}
	}

	c.JSON(200, gin.H{
			"message": "Successfully added wishlist products to the cart",
	})
}




// RemoveFromWishlist removes a product from the user's wishlist.
// @Summary Remove from Wishlist
// @Description Removes a product from the user's wishlist.
// @ID RemoveFromWishlist
// @Param wishlist_id path int true "Wishlist item ID"
// @Success 200 {json} SuccesResponse "Product successfully removed from wishlist."
// @Failure 400 {json} ErrorResponse "Integer conversion error or failed to find in wishlist"
// @Security Bearer
// @Router /user/wishlist/delete/{wishlist_id} [delete]
func RemoveFromWishlist(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	wishlistId, err := strconv.Atoi(c.Param("wishlist_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"errro": "integer converting error",
		})
		return
	}
	err = database.DB.Where("wishtlist_id=? AND user_id=?", wishlistId, userId).Delete(&models.Wishlist{}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "failed to find in wishlist",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted from wishlist",
	})
}
