package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)
type brands struct {
	BrandName string
}

// AddBrand creates a new brand.
// @Summary Create a new brand.
// @Description Create a new brand with the specified name.
// @Tags Brands
// @Produce json
// @Param brand body brands true "Brand name to be created"
// @Success 200 {json} BrandResponse "Successfully created brand"
// @Failure 400 {json} BrandResponse "Error message"
// @Router /admin/brands [post]
func AddBrand(c *gin.Context) {
	var brand brands
	if err := c.Bind(&brand); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var dtbrand models.Brand
	database.DB.Where("brand_name=?", brand.BrandName).First(&dtbrand)

	if dtbrand.Brand_Name == brand.BrandName {
		c.JSON(400, gin.H{
			"error": "This brand already exist",
		})
		return
	}

	database.DB.Create(&models.Brand{Brand_Name: brand.BrandName})
	c.JSON(200, gin.H{
		"message": "successfully created " + brand.BrandName,
	})
}