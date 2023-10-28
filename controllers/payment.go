package controllers

import (
	
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
	"github.com/razorpay/razorpay-go"
)


// RazorPay initiates a payment process through Razorpay.
// @Summary Initiate RazorPay Payment
// @Description Initiates a payment process through Razorpay for the authenticated user.
// @Tags razor pay
// @ID RazorPay
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {json} RazorPayResponse "RazorPay payment page"
// @Failure 400 {json} ErrorResponse "Error while initiating RazorPay payment"
// @Router /user/razorpay [get]
func RazorPay(c *gin.Context) {
	// Check if the user is authenticated and retrieve their user object
	user, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{
			"Error": "User not authenticated",
		})
		return
	}

	// Fetch the user's ID based on their username
	id := user.(models.User).User_ID

	DB := database.InitDB()

	var userdata models.User
	// fetch the user id
	result := DB.Find(&userdata, "User_ID = ?", id)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	// Fetch the total price from the table carts
	var amount uint
	row := DB.Table("carts").Where("user_id = ?", id).Select("SUM(total_price)").Row()
	err := row.Scan(&amount)

	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Sending the payment details to Razorpay
	client := razorpay.NewClient(os.Getenv("RAZOR_kEY"), os.Getenv("RAZOR_SECRET"))
	data := map[string]interface{}{
		"amount":   amount * 100,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	// Creating the payment details for the client order
	body, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"Error test": err,
		})
		return
	}

	// Rendering the HTML page with user & payment details
	value := body["id"]

	c.HTML(200, "app.html", gin.H{
		"userid":     user.(models.User).User_ID,
		"totalprice": amount,
		"paymentid":  value,
	})
}



// RazorPaySuccess processes a successful RazorPay payment and creates an order.
// @Summary Process RazorPay Success
// @Description Processes a successful RazorPay payment and creates an order.
// @Tags razor pay
// @ID RazorPaySuccess
// @Produce json
// @Param order_id query string true "Order ID"
// @Param payment_id query string true "Payment ID"
// @Param signature query string true "Payment Signature"
// @Param total query string true "Total Amount Paid"
// @Security ApiKeyAuth
// @Success 200 {json} RazorPaySuccesResponse "Payment processed successfully"
// @Failure 400 {json} ErrorResponse "Error while processing payment"
// @Router /user/payment/success [get]
func RazorPaySuccess(c *gin.Context) {
	//getting user details from middleware
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	orderid := c.Query("order_id")
	paymentid := c.Query("payment_id")
	signature := c.Query("signature")
	totalamount := c.Query("total")

	err := database.DB.Create(&models.RazorPay{
		User_Id:          uint(userId),
		RazorPayment_id:  paymentid,
		Signature:        signature,
		RazorPayOrder_id: orderid,
		AmountPaid:       totalamount,
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"Error": err.Error(),
		})
		return
	}

	//searching for database all cart data
	var cartdata []models.Cart
	err = database.DB.Where("user_id=?", userId).Find(&cartdata).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Please check your cart",
		})
		return
	}
	//getting total price of cart
	var totalprice uint
	err = database.DB.Table("carts").Select("SUM(total_price)").Where("user_id=?", userId).Scan(&totalprice).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Failed to find total price",
			"message": "cart is empty",
		})
		return
	}
	//checking stock level
	var product models.Product
	for _, v := range cartdata {
		database.DB.First(&product, v.Product_ID)
		if product.Stock-v.Quantity < 0 {
			c.JSON(400, gin.H{
				"error": "Please check quantity",
			})
			return
		}
	}

	var order models.Order
	if err := c.Bind(&order); err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	database.DB.Create(&models.Payment{
		Payment_Type:   "RAZOR PAY",
		Total_Amount:   totalprice,
		Payment_Status: "Completed",
		User_ID:        userId,
		Date:           time.Now(),
	})
	var payment models.Payment
	database.DB.Last(&payment)
	var address models.Address
	err = database.DB.Where("user_id=? AND address_id=?", userId, order.Address_ID).First(&address).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find address choose different id",
		})
		return
	}
	err = database.DB.Create(&models.Order{
		User_ID:     userId,
		Address_ID:  order.Address_ID,
		Total_Price: totalprice,
		Payment_ID:  payment.Payment_ID,
		Status:      "processing",
	}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var cartbrand struct {
		Brand_Name    string
		Category_Name string
	}

	var order1 models.Order
	database.DB.Last(&order1)
	for _, cartdata := range cartdata {
		database.DB.Table("products").Select("brand_name,category_name").Where("id=?", cartdata.Product_ID).Scan(&cartbrand)
		err = database.DB.Create(&models.OrderItem{
			Order_ID:    order1.Order_ID,
			User_ID:     userId,
			Product_ID:  cartdata.Product_ID,
			Address_ID:  order.Address_ID,
			Brand:       cartbrand.Brand_Name,
			Category:    cartbrand.Category_Name,
			Quantity:    cartdata.Quantity,
			Price:       cartdata.Price,
			Total_Price: cartdata.Total_Price,
			Discount:    cartdata.Category_Offer + cartdata.Coupon_Discount,
			Cart_ID:     cartdata.ID,
			Status:      "processing",
			Created_at:  time.Now(),
		}).Error

		if err != nil {
			break
		}
	}
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//reducing the stock count in the database
	var products models.Product
	for _, v := range cartdata {
		database.DB.First(&products, v.Product_ID)
		database.DB.Model(&models.Product{}).Where("id=?", v.Product_ID).Update("stock", product.Stock-v.Quantity)
	}

	err = database.DB.Delete(&models.Cart{}, "user_id=?", userId).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//giving success message
	c.JSON(200, gin.H{
		"message": "successfully ordered your cart",
	})
}


// Success displays a success page after payment completion.
// @Summary Payment Success
// @Description Displays a success page after payment completion.
// @Tags razor pay
// @ID Success
// @Produce html
// @Param id query int true "Payment ID"
// @Success 200 {json} PaymentSuccessResponse "Payment success page"
// @Failure 400 {json} ErrorResponse "Error while displaying success page"
// @Router /user/success [get]
func Success(c *gin.Context) {

	pid, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
