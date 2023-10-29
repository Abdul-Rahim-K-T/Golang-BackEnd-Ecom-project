package controllers

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)
// @Summary Add a new category
// @Description Add a new category to the database
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object that needs to be added"
// @Success 200 {json} SuccessResponse "Successfully created a category"
// @Failure 400 {json} ErrorResponse "Error in adding category"
// @Failure 401 {json} ErrorResponse "Error in request"
// @Router admin/categories [post]

func AddCategory(c *gin.Context) {
	var category models.Category

	if err := c.Bind(&category); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		fmt.Println("binding error")
		return
	}
	var dbcat models.Category
	database.DB.Where("Category_Name=?", category.Category_Name).First(&dbcat)

	//checking this catagory already exist in database
	if dbcat.Category_Name == category.Category_Name {
		c.JSON(400, gin.H{
			"error": "This catagory already exist",
		})
		fmt.Println("catagory already exist")
		return
	}
	//adding data into database
	result := database.DB.Create(&models.Category{
		Category_Name: category.Category_Name,
	})
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error test": result.Error.Error(),
		})
		fmt.Println("database error")
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully created " + category.Category_Name + " category",
	})
}

// @Summary List all categories
// @Description List all categories from the database
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {json} Succes "Successfully retrieved categories"
// @Failure 400 {json} ErrorResponse "Error in retrieving categories"
// @Router admin/categories [get]

func ListCategories(c *gin.Context) {
	var categories []models.Category

	result := database.DB.Raw("SELECT * FROM categories").Scan(&categories)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"categories": categories,
	})
}

// @Summary Block a category by ID
// @Description Block a category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID to block"
// @Success 200 {json} SuccessResponse "Successfully blocked a category"
// @Failure 400 {json} ErrorResponse "Error in blocking category"
// @Failure 401 {json} ErrorResponse "Error in request"
// @Router admin/categories/{category_id} [put]


func BlockCategory(c *gin.Context) {
	id := c.Param("category_id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(401, gin.H{
			"err": err.Error(),
		})
		return
	}

	var category models.Category
	result := database.DB.First(&category, intId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if category.Unlist {
		c.JSON(401, gin.H{
			"error": "this category already blocked",
		})
		return
	}
	result = database.DB.Model(&models.Category{}).Where("category_id=?", intId).Update("unlist", true)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully blocked " + category.Category_Name + " category",
	})
}

// @Summary Unblock a category by ID
// @Description Unblock a category by its ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param category_id path int true "Category ID to unblock"
// @Success 200 {json} SuccessResponse "Successfully unblocked a category"
// @Failure 400 {json} ErrorResponse "Error in unblocking category"
// @Failure 401 {json} ErrorResponse "Error in request"
// @Router admin/categories/{category_id}/unblock [put]



func UnBlockCategory(c *gin.Context) {
	id := c.Param("category_id")
	intId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(401, gin.H{
			"err": err.Error(),
		})
		return
	}

	var category models.Category
	result := database.DB.First(&category, intId)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if !category.Unlist {
		c.JSON(401, gin.H{
			"error": "this category already unblocked",
		})
		return
	}
	result = database.DB.Model(&models.Category{}).Where("category_id=?", intId).Update("unlist", false)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully unblocked " + category.Category_Name + " category",
	})
}
