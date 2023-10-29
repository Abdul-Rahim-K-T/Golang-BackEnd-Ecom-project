package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/controllers"
	"github.com/raheem/Ecom/middleware"
)

func AdminRoutes(r *gin.Engine) {
	router := r.Group("/admin")
	{
		//login and logout admin
		router.POST("/login", controllers.AdminLogin)
		router.GET("/logout", middleware.AdminAuth, controllers.AdminLogout)

		//user
		router.GET("/viewusers", middleware.AdminAuth, controllers.ListUsers)
		router.POST("/blockuser/:user_id", middleware.AdminAuth, controllers.BlockUser)
		router.POST("/unblockuser/:user_id", middleware.AdminAuth, controllers.UnblockUser)

		//CATEGORY
		router.POST("/addcategory", middleware.AdminAuth, controllers.AddCategory)                      //
		router.GET("/viewcategories", middleware.AdminAuth, controllers.ListCategories)                 //
		router.POST("/blockcategory/:category_id", middleware.AdminAuth, controllers.BlockCategory)     //
		router.POST("/unblockcategory/:category_id", middleware.AdminAuth, controllers.UnBlockCategory) //

		//PRODUCTS
		router.POST("/addproduct", middleware.AdminAuth, controllers.AddProduct)
		router.PUT("/editproduct", middleware.AdminAuth, controllers.EditProduct)
		router.POST("/deleteproduct/:product_id", middleware.AdminAuth, controllers.DeleteProduct)
		router.POST("/addimage", middleware.AdminAuth, controllers.AddImage)
		router.PATCH("/editimage",middleware.AdminAuth,controllers.EditImage)

		//BRAND
		router.POST("/addbrand", middleware.AdminAuth, controllers.AddBrand)
		//ORDER MANAGEMENT
		router.GET("/veiwallorders", middleware.AdminAuth, controllers.ViewOrders)
		router.POST("/cancelorder/:order_id", middleware.AdminAuth, controllers.CancelOrder)
		router.PATCH("/changestatus/:order_id", middleware.AdminAuth, controllers.ChangeStatus)

		//SALES REPORT

		//COUPEN
		router.POST("/addcoupon", middleware.AdminAuth, controllers.AddCoupon)
		router.GET("/listcoupons", middleware.AdminAuth, controllers.ListCoupons)
		router.GET("/cancelcoupons", middleware.AdminAuth, controllers.CancelCoupons)

		router.GET("/dashboard", middleware.AdminAuth, controllers.AdminDashboard)
		//CATEGORY OFFERS
		router.POST("/addcategoryoffer", middleware.AdminAuth, controllers.AddCategoryOffer)
		router.GET("/listategoryoffer", middleware.AdminAuth, controllers.ListCatagoryOffer)
		router.GET("/cancelcategoryOffer/:offer_id", middleware.AdminAuth, controllers.DeleteCategoryOffer)
	}
}
