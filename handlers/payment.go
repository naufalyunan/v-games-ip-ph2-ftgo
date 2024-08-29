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

func CreatePayment(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	var payment models.Payment
	var response models.Response

	//initial note (no discount used)
	response.Note = "Using no discount code"

	if err := c.Bind(&payment); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	//find if cart_id exists
	var cart models.Cart
	if err := config.DB.Where("id = ?", payment.CartID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Cart not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	if cart.UserID != uint(userID) {
		return utils.HandleError(c, utils.NewUnauthorizedError("Unauthorized action"))
	}

	//check if the user input coupon code
	if payment.CouponCode != "" {
		//if yes, check if the coupon code exists
		var coupon models.CouponCode

		if err := config.DB.Where("code = ?", payment.CouponCode).First(&coupon).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				//revert to original price if coupon code not found
				payment.PaymentPrice = cart.TotalPrice
				response.Note = "Invalid coupon code, revert to original price"
			} else {
				return utils.HandleError(c, utils.NewInternalError("Internal server error"))
			}
		} else {
			//if exists, adjust the payment price
			disc := cart.TotalPrice * float64(coupon.Discount) / 100
			fmt.Println("discount = ", float64(coupon.Discount))
			payment.PaymentPrice = cart.TotalPrice - disc
			response.Note = fmt.Sprintf("Using coupon with code %s and get %d%% discount", coupon.Code, coupon.Discount)
		}

	} else {
		payment.PaymentPrice = cart.TotalPrice
	}

	//if exists then validate other input
	if payment.PaymentPrice == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Cannot create payment for empty cart"))
	}

	payment.PaymentStatus = "pending"

	// create payment entity
	if err := config.DB.Create(&payment).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	//after creating payment, delete cart and discount code
	if err := config.DB.Delete(&cart, cart.ID).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error deleting cart after creating payment"))
	}

	if payment.PaymentPrice != cart.TotalPrice {
		if err := config.DB.Where("code = ?", payment.CouponCode).Delete(&models.CouponCode{}).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Error deleting cart after creating payment"))
		}
	}

	response.Message = "Success creating payment unit"
	response.Data = payment

	return c.JSON(http.StatusCreated, response)

}
