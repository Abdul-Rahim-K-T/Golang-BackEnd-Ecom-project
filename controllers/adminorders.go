package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)


// ViewOrders lists orders with pagination.
// @Summary List orders with pagination.
// @Description List orders with pagination, specifying the page and limit.
// @Tags ORDER MANAGEMENT
// @Produce json
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Success 200 {json} []orderDetails "List of orders"
// @Failure 500 {json} ErrorResponse "Internal server error"
// @Router /admin/orders [get]

func ViewOrders(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
			c.JSON(500, gin.H{
					"error": "Query page error",
			})
			return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
			c.JSON(500, gin.H{
					"error": "query limit error",
			})
			return
	}
	offset := (page - 1) * limit

	type orderDetails struct {
			Product_Name string
			Quantity     string
			Price        uint
			Total_Price  uint
			Status       string
			Order_ID     uint
			Brand        string
			Category     string
			Address_ID   uint
	}
	// Fetching data from the database and joining the product table to know product details
	var orders []orderDetails
	err = database.DB.Table("order_items").
			Select("products.product_name, order_items.quantity, order_items.price, order_items.total_price, order_items.status, order_items.order_id, brands.brand_name as Brand, categories.category_name as Category, order_items.address_id").
			Joins("INNER JOIN products ON products.id = order_items.product_id").
			Joins("INNER JOIN brands ON brands.Brand_ID = products.Brand_ID").
			Joins("INNER JOIN categories ON categories.category_id = products.category_id").
			Limit(limit).Offset(offset).
			Scan(&orders).Error
	if err != nil {
			c.JSON(400, gin.H{
					"error": "Scanning error",
			})
			return
	}

	c.JSON(200, gin.H{
			"orders": orders,
	})
}

// @Summary Cancel an order.
// @Description Cancel an order by its order ID.
// @Tags ORDER MANAGEMENT
// @Produce json
// @Param order_id path int true "Order ID to be canceled"
// @Success 200 {json} SuccessResponse "Successfully canceled order"
// @Failure 400 {json} ErrorResponse "Error message"
// @Router /admin/orders/cancel/{order_id} [post]

func CancelOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Integer convertion error",
		})
		return
	}
	var order models.Order
	err = database.DB.Where("order_id=?", id).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find this order",
		})
		return
	}

	if order.Status == "cancelled" {
		c.JSON(400, gin.H{
			"error": "This order already cancelled",
		})
		return
	}

	err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Updation error",
		})
		return
	}
	var orderItems []models.OrderItem
	if err = database.DB.Where("order_id=?", id).Find(&orderItems).Error; err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find order items",
		})
		return
	}
	var product models.Product
	for _, v := range orderItems {
		database.DB.First(&product, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock+v.Quantity)
	}
	err = database.DB.Model(&models.Payment{}).Where("payment_id=?", order.Payment_ID).Update("payment_status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully cancelled order",
	})

}

type data struct {
	Status string
}

// @Summary Change the status of an order.
// @Description Change the status of an order by its order ID, providing the new status.
// @Tags ORDER MANAGEMENT
// @Produce json
// @Param order_id path int true "Order ID to be updated"
// @Param status body data true "New status for the order"
// @Success 200 {json} SuccessResponse "Successfully updated status"
// @Failure 400 {json} ErrorResponse "Error message"
// @Router /admin/orders/change-status/{order_id} [patch]

func ChangeStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Integer convertion error",
		})
		return
	}
	var order models.Order
	err = database.DB.Where("order_id=?", id).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find this order",
		})
		return
	}
	var input data
	if err = c.Bind(&input); err != nil {
		c.JSON(400, gin.H{
			"error": "Binding error",
		})
		return
	}
	if input.Status == "shipped" || input.Status == "pending" || input.Status == "cancelled" {
		if input.Status == "cancelled" {
			if order.Status == "cancelled" {
				c.JSON(400, gin.H{
					"error": "This order already cancelled",
				})
				return
			}

			err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			var orderItems []models.OrderItem
			if err = database.DB.Where("order_id=?", id).Find(&orderItems).Error; err != nil {
				c.JSON(400, gin.H{
					"error": "Failed to find order items",
				})
				return
			}
			var product models.Product
			for _, v := range orderItems {
				database.DB.First(&product, v.Product_ID)
				database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock+v.Quantity)
			}
			err = database.DB.Model(&models.Payment{}).Where("payment_id=?", order.Payment_ID).Update("payment_status", "cancelled").Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "updation error",
				})
				return
			}

			c.JSON(200, gin.H{
				"message": "successfully cancelled order",
			})
		} else {
			err = database.DB.Model(&models.Order{}).Where("order_id=?", id).Update("status", input.Status).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", id).Update("status", input.Status).Error
			if err != nil {
				c.JSON(400, gin.H{
					"error": "Updation error",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "successfully updated status",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"error": "this status not applicable",
		})
	}
}
