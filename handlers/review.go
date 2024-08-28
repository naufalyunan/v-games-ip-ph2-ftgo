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

func GetReviews(c echo.Context) error {
	var reviews []*models.Review

	if err := config.DB.Preload("Game").Preload("User").Find(&reviews).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success getting reviews list",
		Data:    reviews,
	})
}

func CreateReview(c echo.Context) error {
	userID := c.Get("user_id").(float64)
	var review models.Review

	if err := c.Bind(&review); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	if review.Rating == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Rating must not be empty"))
	}

	if review.GameID == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Game ID must not be empty"))
	}

	if review.Rating < 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Rating must be a positive number"))
	}

	if review.Message == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Message must not be empty"))
	}

	//find game id if exists
	var game models.Game
	if err := config.DB.Preload("DLCs").Preload("Reviews").Where("id = ?", review.GameID).First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Game not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	review.UserID = uint(userID)

	//if exists then create review

	if err := config.DB.Create(&review).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Success create review for game with ID" + fmt.Sprintf("%d", review.GameID),
		Data:    review,
	})
}
