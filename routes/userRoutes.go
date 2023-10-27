package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/controllers"
	"github.com/raheem/Ecom/middleware"
)

func UserRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*.html")

	router := r.Group("/user")
	{
		//USERLOGIN
		router.POST("/signup", controllers.Signup)
		router.POST("/signup/validate", controllers.ValidateOtp)
		router.POST("/login", controllers.Login)
		router.GET("/logout", middleware.UserAuth, controllers.Logout)

		//PRODUCT
		router.GET("/listproducts", controllers.ListProducts)
		router.GET("/listproductsquery", controllers.ProductDetails)

		//PRODUCTS FILTER

		router.GET("/filterproducts", controllers.FilterProducts)
		router.GET("/sortproducts", controllers.SortProducts)
		// router.GET("/ascendingfilter", controllers.SortWithAscending)
		// router.GET("/descendingfilter", controllers.SortWithDescending)

		//CART
		router.POST("/addtocart", middleware.UserAuth, controllers.AddToCart) //
		router.GET("/viewcart", middleware.UserAuth, controllers.ListCart)
		router.DELETE("/deletefromcart/:cart_id", middleware.UserAuth, controllers.DeleteFromCart)
		router.PUT("/updatecartquantity", middleware.UserAuth, controllers.UpdateCartQuantity)

		//ADDRESS
		router.POST("/addaddress", middleware.UserAuth, controllers.AddAddress)
		router.PUT("/editaddress/:address_id", middleware.UserAuth, controllers.EditAddress)
		router.GET("/listaddresses", middleware.UserAuth, controllers.ListAddresses)

		//USERDETAILS
		router.GET("/userdetail", middleware.UserAuth, controllers.UserDetail)
		router.PUT("/changepassword", middleware.UserAuth, controllers.ChangePassword)
		router.PUT("/updateprofile", middleware.UserAuth, controllers.UpdateProfile)

		//COD
		router.POST("/checkoutcod", middleware.UserAuth, controllers.CheckOutCOD)

		//LIST ORDERS
		router.GET("/listorders", middleware.UserAuth, controllers.ListOrders)
		router.POST("/cancelorder/:order_id", middleware.UserAuth, controllers.CancelOrderWithId)

		//LIST ORDERS WITH FILTERS AND SORT
		router.GET("/listorderswithbrand", middleware.UserAuth, controllers.ListOrdersWithBrand)
		router.GET("/listorderswithcatagory", middleware.UserAuth, controllers.ListOrdersWithCatagory)
		router.GET("/listorderdescasc", middleware.UserAuth, controllers.ListOrdersdescasc)
		

		//RAZOR PAY
		router.GET("/razorpay", middleware.UserAuth, controllers.RazorPay)
		router.POST("/payment/success", middleware.UserAuth, controllers.RazorPaySuccess)
		router.GET("/success", middleware.UserAuth, controllers.Success)

		//APPLY COUPON
		router.POST("/applycoupon", middleware.UserAuth, controllers.ApplyCoupon)

		//WISHLIST
		router.POST("/addtowishlist", middleware.UserAuth, controllers.AddToWishtList)
		router.GET("/listwishlist/:page", middleware.UserAuth, controllers.ListWishlist)
		router.POST("/wishlist/addtocart", middleware.UserAuth, controllers.AddToCartFromWishlist)
		router.POST("/wishlist/delete/:wishlist_id", middleware.UserAuth, controllers.RemoveFromWishlist)

		router.GET("/createinvoice",middleware.UserAuth,controllers.GeneratePdf)
		router.GET("/downloadinvoice",middleware.UserAuth,controllers.DownloadInvoice)
	}
}
