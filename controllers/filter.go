package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
	"gorm.io/gorm"
)

// SortWithAscending sorts and retrieves a list of products in ascending order of their prices.
// @Summary Sort products by ascending price
// @Description Retrieves a list of products sorted in ascending order of their prices.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/ascendingfilter [get]
func SortProducts(c *gin.Context) {
	sortOrder := c.DefaultQuery("sortOrder", "asc") // Get the sort order from the query string, default to "asc"

	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Category_Name string
		Stock         uint
		Image         string // Include the "Image" field
	}

	var products []product
	query := database.DB.Table("products").Select("products.id, products.product_name, products.price, brands.brand_name, categories.category_name, products.stock, images.image").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").
		Joins("INNER JOIN categories ON categories.category_id=products.category_id").
		Joins("INNER JOIN images ON images.product_id=products.id")

	// Determine the sorting order based on the "sortOrder" parameter
	if sortOrder == "asc" {
		query = query.Order("products.price asc")
	} else {
		query = query.Order("products.price desc")
	}

	err := query.Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}

	c.JSON(200, gin.H{
		"products": products,
	})
}

// FilterWithBrands filters products by brand name.
// @Summary Filter products by brand name
// @Description Retrieves a list of products associated with the specified brand name.
// @Tags products
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param brand_name query string true "Brand name to filter by"
// @Security ApiKeyAuth
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/brandfilter [get]
func FilterWithBrands(c *gin.Context) {
	brand := c.Query("brand_name")

	var branddb models.Brand
	err := database.DB.Where("brand_name=?", brand).First(&branddb).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "This brand not exist",
		})
		return
	}
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Category_Name string
		Stock         uint
	}
	var products []product
	err = database.DB.Table("products").Select("products.id,products.product_name,products.price,brands.brand_name,categories.category_name,products.stock").
		Joins("INNER JOIN brands ON brands.brand_id=products.brand_id").Joins("INNER JOIN categories ON categories.category_id=products.category_id").
		Where("products.brand_id=?", branddb.Brand_ID).Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}
	c.JSON(200, gin.H{
		"products": products,
	})

}


// @Summary Filter products
// @Description Filter products by brand, category, and sort them
// @ID FilterProducts
// @Produce json
// @Param brand_name query string false "Brand name to filter by"
// @Param category_name query string false "Category name to filter by"
// @Param sort_by query string false "Field to sort by (product_name, price, stock)"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Success 200 {json} SuccessResponse
// @Failure 400 {string} ErrorResponse
// @Router /user/filterproducts [get]

func FilterProducts(c *gin.Context) {
	brandName := c.Query("brand_name")
	categoryName := c.Query("category_name")
	sortBy := c.DefaultQuery("sort_by", "product_name") // Default to sorting by product name
	sortOrder := c.DefaultQuery("sort_order", "asc")    // Default to ascending order

	// Check if brand filter is provided
	if brandName != "" {
		filterAndSort(c, "brand_name", brandName, sortBy, sortOrder)
		return
	}

	// Check if category filter is provided
	if categoryName != "" {
		filterAndSort(c, "category_name", categoryName, sortBy, sortOrder)
		return
	}

	// If no filter is provided, apply only sorting
	filterAndSort(c, "", "", sortBy, sortOrder)
}

func filterAndSort(c *gin.Context, filterField, filterValue, sortBy, sortOrder string) {
	type product struct {
		ID            uint
		Product_Name  string
		Price         uint
		Brand_Name    string
		Category_Name string
		Stock         uint
		Image         string // Add image field
	}
	var products []product

	query := database.DB.Table("products").
		Select("products.id, products.product_name, products.price, brands.brand_name, categories.category_name, products.stock, images.image").
		Joins("INNER JOIN brands ON brands.brand_id = products.brand_id").
		Joins("INNER JOIN categories ON categories.category_id = products.category_id").
		Joins("INNER JOIN images ON images.product_id = products.id") // Join the images table

	// Apply partial string matching filter
	if filterField != "" && filterValue != "" {
		query = query.Where(filterField+" LIKE ?", "%"+filterValue+"%")
	}

	// Apply sorting
	switch sortBy {
	case "product_name":
		query = applySorting(query, "products.product_name", sortOrder)
	case "price":
		query = applySorting(query, "products.price", sortOrder)
	case "stock":
		query = applySorting(query, "products.stock", sortOrder)
	default:
		query = applySorting(query, "products.product_name", sortOrder) // Default sorting by product name
	}

	err := query.Scan(&products).Error

	if err != nil {
		c.JSON(400, gin.H{
			"error": "database error",
		})
		return
	}

	c.JSON(200, gin.H{
		"products": products,
	})
}




func applySorting(query *gorm.DB, field, sortOrder string) *gorm.DB {
	switch sortOrder {
	case "asc":
		return query.Order(field)
	case "desc":
		return query.Order(field + " DESC")
	default:
		return query.Order(field)
	}
}
