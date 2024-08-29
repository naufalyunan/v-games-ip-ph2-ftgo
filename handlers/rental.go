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

func Pay(c echo.Context) error {
	paymentID := c.Param("id")

	//bind the request payload
	var payload models.UpdatePaymentStatusPayload

	if err := c.Bind(&payload); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid Input"))
	}

	//get payment entity
	var payment models.Payment
	if err := config.DB.Preload("Cart.CartItems").Where("id = ?", paymentID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("No Payment Found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	var rental models.Rental
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		// if valid, then update the payment status
		if err := tx.Model(&payment).Where("id = ?", paymentID).Update("payment_status", payload.Status).Error; err != nil {
			return utils.NewInternalError("Internal Server Error")
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
