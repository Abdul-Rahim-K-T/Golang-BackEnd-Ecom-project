package controllers

import (
	
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raheem/Ecom/database"
	"github.com/raheem/Ecom/models"
)


// @Summary Add a new address
// @Description Add a new address to the user's profile
// @Tags addresses
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer token"
// @Param address body models.Address true "Address object to add"
// @Success 200 {json} success
// @Failure 400 {json} error
// @Failure 401 {json} error
// @Router /user/addaddress [post]

func AddAddress(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	var address models.Address
	if err := c.Bind(&address); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(address.Zip_code) != 6 {
		c.JSON(400, gin.H{
			"error": "zip code must be 6 number",
		})
		return
	}
	err := database.DB.Create(&models.Address{
		Building_Name: address.Building_Name,
		City:          address.City,
		State:         address.State,
		Landmark:      address.Landmark,
		Zip_code:      address.Zip_code,
		User_ID:       userId,
	}).Error
	if err != nil {
		c.JSON(400, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"message": "success fully created new address",
	})

}

// EditAddress godoc
// @Summary Edit an existing address
// @Description Edit an existing address for the authenticated user
// @Tags Addresses
// @Accept json
// @Produce json
// @Param address_id path int true "Address ID to be edited"
// @Param updatedata body struct {} true "Update data for the address"
// @Success 200 {object} gin.H "Successfully updated address"
// @Failure 400 {object} map[string]interface{} "Error message"
// @Router /user/editaddress/:{address_id} [put]


func EditAddress(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID
	adid, err := strconv.Atoi(c.Param("address_id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	var address models.Address
	err = database.DB.Where("address_id=? AND user_id=?", adid, userId).First(&address).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find the address",
		})
		return
	}
	var updatedata struct {
		NewBuildingname string `json:"newbuildingname"`
		NewCity         string `json:"newcity"`
		NewState        string `json:"newstate"`
		NewLandmark     string `json:"landmark"`
		NewZip          string `json:"newzip"`
	}
	if err := c.Bind(&updatedata); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(updatedata.NewZip) != 6 {
		c.JSON(400, gin.H{
			"error": "zip code must be 6 number",
		})
		return
	}
	err = database.DB.Model(&models.Address{}).Where("address_id=? AND user_id=?", adid, userId).Updates(map[string]interface{}{"building_name": updatedata.NewBuildingname, "city": updatedata.NewCity, "state": updatedata.NewState, "zip_code": updatedata.NewZip, "landmark": updatedata.NewLandmark}).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully updated address",
	})
}


// ListAddresses godoc
// @Summary List user addresses
// @Description List addresses of the authenticated user
// @Tags Addresses
// @Produce json
// @Success 200 {object} gin.H "List of addresses"
// @Failure 400 {object} map[string]interface{} "Error message"
// @Router /user/listaddresses [get]

func ListAddresses(c *gin.Context) {
	user, _ := c.Get("user")
	userId := user.(models.User).User_ID

	type address struct {
		Address_ID    uint
		Building_Name string
		City          string
		State         string
		Landmark      string
		Zip_Code      string
	}

	var addresses []address
	err := database.DB.Table("addresses").Select("address_id", "building_name", "city", "state", "landmark", "zip_code").Where("user_id=?", userId).Find(&addresses).Error
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Failed to find addressess",
		})
		return
	}
	c.JSON(200, gin.H{
		"addresses": addresses,
	})

}