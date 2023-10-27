package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)

// ListOrders retrieves a list of orders for the authenticated user.
// @Summary List Orders
// @Description Retrieve a list of orders for the authenticated user.
// @ID ListOrders
// @Produce json
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorders [get]
func ListOrders(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}

	var orders []orderDetails

	err = database.DB.Table("order_items").
		Select("products.product_name, order_items.quantity, order_items.price, order_items.total_price, order_items.status, order_items.order_id, brands.brand_name as Brand, categories.category_name as Category, order_items.address_id, images.image").
		Joins("INNER JOIN products ON products.id = order_items.product_id").
		Joins("INNER JOIN images ON images.product_id = order_items.product_id").
		Joins("INNER JOIN brands ON brands.Brand_ID = products.Brand_ID").
		Joins("INNER JOIN categories ON categories.category_id = products.category_id").
		Where("order_items.user_id = ?", userId).
		Limit(limit).Offset(offset).
		Scan(&orders).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}

	// Show the data to the user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}


// ListOrdersWithBrand retrieves a list of orders with a specific brand name.
// @Summary List Orders by Brand
// @Description Retrieve a list of orders for the authenticated user with a specific brand name.
// @ID ListOrdersWithBrand
// @Produce json
// @Param brandname query string true "Partial brand name to filter by"
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Security ApiKeyAuth
// @Success 200 {json} SuccesResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderswithbrand [get]
func ListOrdersWithBrand(c *gin.Context) {
	brandPartialName := c.Query("brandname") // Change the query parameter name

	if brandPartialName == "" {
		c.JSON(400, gin.H{
			"error": "Brand name is required",
		})
		return
	}

	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	// Fetching data from the database and inner joining the product table to retrieve product details, including Brand and Category.
	var orders []orderDetails
	err = database.DB.Table("order_items").
		Select("products.product_name, order_items.quantity, order_items.price, order_items.total_price, order_items.status, order_items.order_id, brands.brand_name as Brand, categories.category_name as Category, order_items.address_id, images.image").
		Joins("INNER JOIN products ON products.id = order_items.product_id").
		Joins("INNER JOIN images ON images.product_id = order_items.product_id").
		Joins("INNER JOIN brands ON brands.Brand_ID = products.Brand_ID").
		Joins("INNER JOIN categories ON categories.category_id = products.category_id").
		Where("order_items.user_id = ? AND brands.brand_name LIKE ?", userId, "%"+brandPartialName+"%"). // Use the LIKE operator for partial search
		Limit(limit).Offset(offset).
		Scan(&orders).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	// Show the data to the user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}


// ListOrdersWithCategory retrieves a list of orders with a specific category name.
// @Summary List Orders by Category
// @Description Retrieve a list of orders for the authenticated user with a specific category name.
// @ID ListOrdersWithCategory
// @Produce json
// @Param categoryname query string true "Partial category name to filter by"
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Security ApiKeyAuth
// @Success 200 {json} ListOrdersCategoryResponse "List of orders"
// @Failure 400 {json} ErrorResponse "Error while fetching orders"
// @Router /user/listorderswithcatagory [get]
func ListOrdersWithCatagory(c *gin.Context) {
	categoryPartialName := c.Query("categoryname") // Change the query parameter name

	if categoryPartialName == "" {
		c.JSON(400, gin.H{
			"error": "Category name is required",
		})
		return
	}

	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
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
		Image        string
	}
	// Fetching data from the database and inner joins product table knowing product details
	var orders []orderDetails

	err = database.DB.Table("order_items").
		Select("products.product_name, order_items.quantity, order_items.price, order_items.total_price, order_items.status, order_items.order_id, brands.brand_name as Brand, categories.category_name as Category, order_items.address_id, images.image").
		Joins("INNER JOIN products ON products.id = order_items.product_id").
		Joins("INNER JOIN images ON images.product_id = order_items.product_id").
		Joins("INNER JOIN brands ON brands.Brand_ID = products.Brand_ID").
		Joins("INNER JOIN categories ON categories.category_id = products.category_id").
		Where("order_items.user_id = ? AND categories.category_name LIKE ?", userId, "%"+categoryPartialName+"%"). // Use the LIKE operator for partial search
		Limit(limit).Offset(offset).
		Scan(&orders).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}
	// Show the data to the user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}


// ListOrdersdescasc retrieves a list of orders with sorting in ascending or descending order.
// @Summary List Orders with Sorting
// @Description Retrieve a list of orders for the authenticated user with optional sorting.
// @ID ListOrdersdescasc
// @Produce json
// @Param page query int true "Page number for pagination"
// @Param limit query int true "Number of items per page"
// @Param sort query string false "Sort order (asc or desc)"
// @Security ApiKeyAuth
// @Success 200 {json} SuccesResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/listorderdescasc [get]
func ListOrdersdescasc(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(500, gin.H{
			"error": "query value error",
		})
		return
	}
	offset := (page - 1) * limit

	sortOrder := c.DefaultQuery("sort", "asc") // Default to ascending order if no sort parameter provided

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
		Image        string
	}

	var orders []orderDetails
	db := database.DB.Table("order_items").
		Select("products.product_name, order_items.quantity, order_items.price, order_items.total_price, order_items.status, order_items.order_id, brands.brand_name as Brand, categories.category_name as Category, order_items.address_id, images.image").
		Joins("INNER JOIN products ON products.id = order_items.product_id").
		Joins("INNER JOIN images ON images.product_id = order_items.product_id").
		Joins("INNER JOIN brands ON brands.Brand_ID = products.Brand_ID").
		Joins("INNER JOIN categories ON categories.category_id = products.category_id").
		Where("user_id = ?", userId).
		Limit(limit).Offset(offset)

	if sortOrder == "asc" {
		db = db.Order("order_items.price asc")
	} else if sortOrder == "desc" {
		db = db.Order("order_items.price desc")
	} else {
		c.JSON(400, gin.H{
			"error": "Invalid sort order parameter",
		})
		return
	}

	err = db.Scan(&orders).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "Scanning error",
		})
		return
	}

	// Show the data to the user
	c.JSON(200, gin.H{
		"orders": orders,
	})
}


// CancelOrderWithId cancels an order with a specific order ID.
// @Summary Cancel Order
// @Description Cancel an order with a specific order ID for the authenticated user.
// @ID CancelOrderWithId
// @Produce json
// @Param order_id path int true "Order ID to cancel"
// @Security ApiKeyAuth
// @Success 200 {json} SuccesResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/cancelorder/{order_id} [put]
func CancelOrderWithId(c *gin.Context) {
	//getting user details from middlewares
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	orderId, err := strconv.Atoi(c.Param("order_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//fetching the data from database
	var order models.Order
	err = database.DB.Where("user_id=? AND order_id=?", userId, orderId).First(&order).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Your order didn't find",
		})
		return
	}

	var payments models.Payment
	database.DB.First(&payments, order.Payment_ID)
	//checking if order is already cancelled or not
	if order.Status == "cancelled" {
		c.JSON(400, gin.H{
			"error": "This order alredy cancelled",
		})
		return
	}
	//updating the status in to cancelled
	err = database.DB.Model(&models.Order{}).Where("user_id=? AND order_id=?", userId, orderId).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	err = database.DB.Model(&models.OrderItem{}).Where("order_id=?", orderId).Update("status", "cancelled").Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "updation error",
		})
		return
	}

	var orderItems []models.OrderItem
	database.DB.Where("order_id=?", orderId).Find(&orderItems)

	totalprices := 0
	var product models.Product
	for _, v := range orderItems {
		totalprices += int(v.Total_Price)
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
	if payments.Payment_Type == "RAZOR PAY" {
		err = database.DB.Model(&models.User{}).Where("user_id=?", userId).Update("wallet", user.(models.User).Wallet+uint(totalprices)).Error
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Database error",
			})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "successfully cancelled order",
	})
}
