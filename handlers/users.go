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

// @Summary Update Deposit Amount
// @Description Update a user deposit value based on it's user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.DepoPayload true "Depo Amount"
// @Success 200 {object} models.Response
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /users/depo [PUT]
func TopUp(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	var payload models.DepoPayload

	if err := c.Bind(&payload); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	if payload.Amount < 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Deposit must be a positive number"))
	}

	//validate if user exist
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("User not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	//if exists, update
	if err := config.DB.Model(&user).Where("id = ?", userID).Update("deposit", payload.Amount).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: fmt.Sprintf("Success update deposit of user with ID %d to IDR%.2f", uint(userID), user.Deposit),
		Data:    user,
	})
}
