package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)

type carts struct {
	ProductId uint `json:"productid"`
	Quantity  uint `json:"quantity"`
}


// AddToCart adds a product to the user's cart.
// @Summary Add a product to the user's cart.
// @Description Add a specified product with a given quantity to the user's shopping cart.
// @Tags Cart
// @Produce json
// @Param cart body carts true "Product ID and Quantity to add"
// @Success 200 {json} AddToCartResponse "Successfully added to cart"
// @Failure 400 {json} AddToCartResponse "Error message"
// @Router /user/addtocart [post]

func AddToCart(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	var cart carts
	if err := c.Bind(&cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to bind",
		})
		c.Abort()
		return
	}
	var product models.Product
	err := database.DB.First(&product, cart.ProductId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This product not found",
		})
		return
	}
	var id []uint
	database.DB.Table("category_offers").Select("id").Where("offer=true").Scan(&id)
	count := 0
	for _, v := range id {
		if v == product.Category_Id {
			count = 1
			break
		}
	}
	if count == 1 {
		var dbcart models.Cart
		err = database.DB.Where("product_id=? AND user_id=?", cart.ProductId, userId).First(&dbcart).Error

		var result models.Category_Offer
		database.DB.Where("category_id=?", product.Category_Id).First(&result)

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:     cart.ProductId,
				Quantity:       cart.Quantity,
				Price:          product.Price,
				Category_Id:    product.Category_Id,
				Total_Price:    (product.Price * cart.Quantity) - ((product.Price * cart.Quantity) * result.Percentage / 100),
				User_ID:        userId,
				Category_Offer: (product.Price * cart.Quantity) * result.Percentage / 100,
			}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success fully added to cart",
			})
			return
		}
		var discount uint
		database.DB.Table("carts").Select("category_offer").Where("product_id=? AND user_id=?", cart.ProductId, userId).Scan(&discount)
		total := ((product.Price) - ((product.Price) * result.Percentage / 100)) * cart.Quantity
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", cart.ProductId, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total, "category_offer": discount + ((product.Price) * result.Percentage / 100)}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "successfully updated cart",
		})
	} else {

		var dbcart models.Cart
		err = database.DB.Where("product_id=? AND user_id=?", cart.ProductId, userId).First(&dbcart).Error

		if err != nil {
			err = database.DB.Create(&models.Cart{
				Product_ID:  cart.ProductId,
				Quantity:    cart.Quantity,
				Price:       product.Price,
				Category_Id: product.Category_Id,
				Total_Price: product.Price * cart.Quantity,
				User_ID:     userId,
			}).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": err.Error(),
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "success fully added to cart",
			})
			return
		}
		total := product.Price * cart.Quantity
		err = database.DB.Model(&models.Cart{}).Where("product_id=? AND user_id=?", cart.ProductId, userId).Updates(map[string]interface{}{"quantity": dbcart.Quantity + cart.Quantity, "total_price": dbcart.Total_Price + total}).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "successfully updated cart",
		})
	}
}


// ListCart retrieves the user's shopping cart items.
// @Summary Retrieve the user's shopping cart items.
// @Description Get the list of products in the user's shopping cart.
// @Tags Cart
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {json} ListCartResponse "List of cart items and total price"
// @Failure 400 {json} AddToCartResponse "Error message"
// @Security ApiKeyAuth
// @Router /user/viewcart [get]
func ListCart(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query error",
		})
		return
	}
	offset := (page - 1) * limit

	type cart struct {
		ID              uint `json:"productid"`
		Category_Id     uint
		Product_name    string
		Quantity        uint
		Price           string
		Total_price     uint
		Coupon_Discount uint
		Category_Offer  uint
		Image           string
	}
	var totalprice int
	var carts []cart
	err = database.DB.Table("carts").
		Select("products.id,products.category_id,products.product_name,carts.quantity,carts.price,carts.total_price,carts.coupon_discount,carts.category_offer,images.image").
		Joins("INNER JOIN products ON products.id=carts.product_id").Joins("INNER JOIN images ON images.product_id=carts.product_id").Where("carts.user_id=?", userId).
		Limit(limit).Offset(offset).
		Scan(&carts).Error

	if err != nil {
		c.JSON(400, gin.H{
			"err": err.Error(),
		})
		return
	}
	database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userId).Scan(&totalprice)
	c.JSON(200, gin.H{
		"total price": totalprice,
		"cartitems":   carts,
	})
}


// DeleteFromCart deletes a product from the user's cart.
// @Summary Delete a product from the user's cart.
// @Description Delete a product from the user's cart by specifying the cart item ID.
// @Tags Cart
// @Produce json
// @Param cart_id path int true "Cart Item ID"
// @Success 200 {json} DeleteCartResponse "Successfully deleted from cart"
// @Failure 400 {json} DeleteCartResponse "Error message"
// @Router /user/deletefromcart/:cart_id [delete]
func DeleteFromCart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cart_id"))
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Println(userId, id)
	result := database.DB.Where("id=? AND user_id=?", id, userId).Delete(&models.Cart{})
	fmt.Println(result.RowsAffected)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully deleted from cart",
	})
}

type cart struct {
	CartId      uint `json:"cartid"`
	NewQuantity uint `json:"newquantity"`
}


// UpdateCartQuantity updates the quantity of a product in the user's cart.
// @Summary Update the quantity of a product in the user's cart.
// @Description Update the quantity of a product in the user's cart by specifying the cart item ID and the new quantity.
// @Tags Cart
// @Produce json
// @Param cart body cart true "Cart Item ID and New Quantity"
// @Success 200 {json} UpdateCartQuantityResponse "Successfully updated quantity"
// @Failure 400 {json} UpdateCartQuantityResponse "Error message"
// @Router /user/updatecartquantity [post]

func UpdateCartQuantity(c *gin.Context) {
	user, _ := c.Get("user")
	id := user.(models.User).User_ID
	var updateData cart
	if err := c.Bind(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if updateData.NewQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "quantity must postive value",
		})
		return
	}
	var dtcart models.Cart
	err := database.DB.First(&dtcart, updateData.CartId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	err = database.DB.Model(&models.Cart{}).Where("user_id=? AND id=?", id, updateData.CartId).Updates(map[string]interface{}{"quantity": updateData.NewQuantity, "total_price": updateData.NewQuantity * dtcart.Price}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully updated quantity",
	})
}
