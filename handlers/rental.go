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

// @Summary Update Payment Status
// @Description update field in the payments table, also create a new rental and rental items automatically
// @Tags pay
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.UpdatePaymentStatusPayload true "Update Payment"
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /pay/{id} [PUT]
func Pay(c echo.Context) error {
	paymentID := c.Param("id")

	//bind the request payload
	var payload models.UpdatePaymentStatusPayload

	if err := c.Bind(&payload); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid Input"))
	}

	//get payment entity
	var payment models.Payment
	if err := config.DB.Preload("Cart.User").Preload("Cart.CartItems").Where("id = ?", paymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("No Payment Found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	if payment.PaymentStatus == "PAID" {
		return utils.HandleError(c, utils.NewBadRequestError("Cart already Paid"))
	}

	//check if deposit is enough, if not returned bad request
	if payment.PaymentPrice > payment.Cart.User.Deposit {
		return utils.HandleError(c, utils.NewBadRequestError("Unsufficient Balance to Pay"))
	}

	var rental models.Rental
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		// if valid, then update the payment status
		if err := tx.Model(&payment).Where("id = ?", paymentID).Update("payment_status", payload.Status).Error; err != nil {
			return utils.NewInternalError("Internal Server Error")
		}

		currentDepo := payment.Cart.User.Deposit - payment.PaymentPrice
		//update the user 
		if err := tx.Model(&models.User{}).Where("id = ?", payment.Cart.UserID).Update("deposit", currentDepo).Error; err != nil {
			return utils.NewInternalError("Fail to update user's deposit value")
		}

		//create rental after updating payment
		rental = models.Rental{
			Status:    "ON RENTAL",
			PaymentID: payment.ID,
		}

		if err := tx.Create(&rental).Error; err != nil {
			return utils.NewInternalError("Error creating rental entity")
		}

		//save each rental item to rental entity
		for _, cartItem := range payment.Cart.CartItems {
			rentalItem := models.RentalItem{
				RentalID:   rental.ID,
				CartItemID: cartItem.ID,
				Status:     "ON RENTAL",
			}
			if err := tx.Create(&rentalItem).Error; err != nil {
				return utils.NewInternalError("Error creating rental item entity")
			}
		}
		// return nil will commit the whole transaction
		return nil
	}); err != nil {
		return utils.HandleError(c, err.(*utils.APIError))
	}

	//get the rental for display
	if err := config.DB.Preload("Payment.Cart").Preload("RentalItems").Where("id = ?", rental.ID).First(&rental).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("No Rental Found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	//after completing, delete cart and discount code
	if err := config.DB.Delete(&models.Cart{}, payment.Cart.ID).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Error deleting cart after creating payment"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: fmt.Sprintf("Success update payment status with ID %d to %s and create rental with ID %d", payment.ID, payment.PaymentStatus, rental.ID),
		Data:    rental,
	})
}

// @Summary Get All Rentals
// @Description Retrieve a list of all rentals
// @Tags rentals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /rentals [get]
func GetRentals(c echo.Context) error {
	userID, ok := c.Get("user_id").(float64)
	if !ok {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid user ID"))
	}

	userRole, ok := c.Get("role").(string)
	if !ok {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid user role"))
	}

	var rentals []models.Rental

	if userRole == "user" {
		// Convert float64 userID to uint as needed
		uid := uint(userID)

		if err := config.DB.Joins("JOIN payments ON payments.id = rentals.payment_id").
			Joins("JOIN carts ON carts.id = payments.cart_id").
			Where("carts.user_id = ?", uid).
			Preload("Payment").
			Preload("Payment.Cart").
			Preload("RentalItems").
			Find(&rentals).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	} else if userRole == "admin" {
		if err := config.DB.Joins("JOIN payments ON payments.id = rentals.payment_id").
			Joins("JOIN carts ON carts.id = payments.cart_id").
			Preload("Payment").
			Preload("Payment.Cart").
			Preload("RentalItems").
			Find(&rentals).Error; err != nil {
			return utils.HandleError(c, utils.NewInternalError("Internal server error"))
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success getting rentals",
		Data:    rentals,
	})
}
