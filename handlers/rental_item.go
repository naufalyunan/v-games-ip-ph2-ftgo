package handlers

import (
	"fmt"
	"net/http"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/models"
	"v-games-ip-ph2-ftgo/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @Summary Update Rental Item Status
// @Description update status field in the rental items table
// @Tags rental-item
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.RentalItem true "Update Rental Item"
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /rental-item/{id} [PUT]
func UpdateRentalItem(c echo.Context) error {
	rentalItemID := c.Param("id")

	var rentalItem models.RentalItem

	if err := c.Bind(&rentalItem); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	var check models.RentalItem
	if err := config.DB.Preload("Rental").Preload("CartItem").Where("id = ?", rentalItemID).First(&check).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Rental Item not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	//if exist, proceed to update
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Model(&rentalItem).Where("id = ?", rentalItemID).Update("status", rentalItem.Status).Error; err != nil {
			return utils.NewInternalError("Error updating rental item status")
		}

		// return nil will commit the whole transaction
		return nil
	}); err != nil {
		return utils.HandleError(c, err.(*utils.APIError))
	}

	// check all rental items of cart, if all is returned, change the rentals to returned
	// get all rental_items
	var rentalItems []*models.RentalItem

	if err := config.DB.
		Preload("CartItem").
		Where("rental_id = ?", check.RentalID).
		Find(&rentalItems).Error; err != nil {
		return utils.NewInternalError("Error getting rental items")
	}
	// return c.JSON(http.StatusOK, rentalItems)
	status := true
	for _, renIt := range rentalItems {
		if renIt.Status != "RETURNED" {
			status = false
		}
	}
	//update rental as done if all items is returned
	if status {
		fmt.Println("SEMUA SUDAH RETURNED")
		var rental models.Rental
		if err := config.DB.Model(&rental).Where("id = ?", check.RentalID).Update("status", rentalItem.Status).Error; err != nil {
			return utils.NewInternalError("Error updating rental item status")
		}
	}

	//GET RENTAL ITEM FOR DISPLAY
	if err := config.DB.Preload("Rental").Preload("CartItem").Where("id = ?", rentalItemID).First(&rentalItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Rental Item not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: fmt.Sprintf("Success update rental item with ID %d status as %s", check.ID, rentalItem.Status),
		Data:    rentalItem,
	})
}
