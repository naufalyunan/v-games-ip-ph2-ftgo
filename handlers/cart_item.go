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

	//check the quantity and stock
	// retrieve the game
	var game models.Game
	if err := config.DB.Where("id = ?", cartItem.GameID).First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Game not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	// then compare the stock and quantity to buy
	if game.Stock < cartItem.Quantity {
		return utils.HandleError(c, utils.NewBadRequestError("Not enough stock to buy"))
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

		//check if the cart is already requested payment and pending status, if yes, then return error
		var payments []*models.Payment
		uid := uint(userID)

		if err := tx.Joins("JOIN carts ON carts.id = payments.cart_id").
			Where("carts.user_id = ?", uid).
			Where("payment_status = ?", "PENDING").
			Preload("Cart").
			Preload("Cart.CartItems").
			Find(&payments).Error; err != nil {
			return utils.NewInternalError("Internal server error")
		}

		if len(payments) > 0 && payments[0].PaymentStatus == "PENDING" && payments[0].CartID == cart.ID {
			//if truthy then the cart is already in payment, so the item could not be added unless is paid
			return utils.NewBadRequestError("Cart already requested payment, finish the payment first, only then can make a new cart")
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

		//update the cartItem new stock
		newStock := cartItem.Game.Stock - cartItem.Quantity
		if err := tx.Model(&cartItem.Game).Update("stock", newStock).Error; err != nil {
			return utils.NewInternalError("Error updating game new stock after adding to cart")
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
