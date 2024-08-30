package handlers

import (
	"net/http"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/models"
	"v-games-ip-ph2-ftgo/utils"

	"github.com/labstack/echo/v4"
)

// @Summary Get All Coupons
// @Description Retrieve a list of all Coupons
// @Tags coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /coupons [get]
func GetCoupons(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	userRole := c.Get("role")
	var coupons []*models.CouponCode

	if userRole == "user" {
		if err := config.DB.Where("usable_by_self = ?", true).Where("user_id = ?", userID).Find(&coupons).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	} else if userRole == "admin" {
		if err := config.DB.Where("usable_by_self = ?", true).Find(&coupons).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Here are discount coupons to use",
		Data:    coupons,
	})
}

// @Summary Get Referreal Code
// @Description Retrieve a list of Referral Code
// @Tags coupons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /coupons/referral [get]
func GetReferralCode(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	userRole := c.Get("role")
	var coupons []*models.CouponCode

	if userRole == "user" {
		if err := config.DB.Where("usable_by_self = ?", false).Where("user_id = ?", userID).Find(&coupons).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	} else if userRole == "admin" {
		if err := config.DB.Where("usable_by_self = ?", false).Find(&coupons).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Referral code is as below",
		Data:    coupons,
	})
}
