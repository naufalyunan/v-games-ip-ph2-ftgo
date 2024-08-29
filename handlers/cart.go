package handlers

import (
	"net/http"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/models"
	"v-games-ip-ph2-ftgo/utils"

	"github.com/labstack/echo/v4"
)

func GetCarts(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	userRole := c.Get("role")
	var carts []*models.Cart

	if userRole == "user" {
		if err := config.DB.Preload("CartItems").Where("user_id = ?", userID).Find(&carts).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	} else if userRole == "admin" {
		if err := config.DB.Preload("CartItems").Find(&carts).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success getting carts",
		Data:    carts,
	})

}
