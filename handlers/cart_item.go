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

func CreateCartItem(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	var cartItem models.CartItem

	if err := c.Bind(&cartItem); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError(err.Error()))
	}

	if cartItem.GameID == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Game ID must not be empty"))
	}

	if cartItem.Quantity == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Quantity must be greater than zero"))
	}

	if cartItem.StartDate == nil || cartItem.EndDate == nil {
		return utils.HandleError(c, utils.NewBadRequestError("Date must be filled"))
	}

	if err := cartItem.ValidateDates(); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("End Date must be greater than Start Date"))
	}

	//check if a cart for a particular user is already exists
	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Println("Masuk bawah")
				//if not found, create cart
				cart.UserID = uint(userID)
				cart.TotalPrice = 0
				if err := tx.Create(&cart).Error; err != nil {
					return utils.NewInternalError("Internal server error")
				}
			} else {
				fmt.Println("Masuk atas")
				return utils.NewInternalError("Internal server error")
			}
		}

		// create the cartItem and assign to cart
		cartItem.CartID = cart.ID
		if err := tx.Create(&cartItem).Error; err != nil {
			return utils.NewInternalError("Error adding item to cart")
		}

		//retrieve the created cartitem
		if err := tx.Preload("Game").Where("id = ?", cartItem.ID).First(&cartItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.NewNotFoundError("Created cart item not found")
			} else {
				return utils.NewInternalError("Internal server error")
			}
		}

		duration := cartItem.CalculateDaysDifference()
		//then update the total price of cart
		updatedTotalPrice := float64(cart.TotalPrice) + float64(cartItem.Quantity)*cartItem.Game.RentPrice*float64(duration)
		if err := tx.Model(&cart).Update("total_price", updatedTotalPrice).Error; err != nil {
			return utils.NewInternalError("Error updating total price of cart")
		}

		// return nil will commit the whole transaction
		return nil

	}); err != nil {
		return utils.HandleError(c, err.(*utils.APIError))
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: fmt.Sprintf("Success adding game with ID %d and qty of %d to cart", cartItem.GameID, cartItem.Quantity),
		Data:    cartItem,
	})

}
