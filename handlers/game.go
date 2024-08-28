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

func GetGames(c echo.Context) error {
	var games []*models.Game
	if err := config.DB.Preload("DLCs").Preload("Reviews").Find(&games).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success getting games",
		Data:    games,
	})
}

func GetGameByID(c echo.Context) error {
	gameID := c.Param("id")

	var game models.Game
	if err := config.DB.Preload("DLCs").Preload("Reviews").Where("id = ?", gameID).First(&game).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Game not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success getting game with ID" + gameID,
		Data:    game,
	})
}

func UpdateGameStock(c echo.Context) error {
	gameID := c.Param("id")
	var request models.UpdateGameStockPayload

	//input validation
	if err := c.Bind(&request); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	//find game if exists
	var updatedGame models.Game
	if err := config.DB.Preload("DLCs").Preload("Reviews").Where("id = ?", gameID).First(&updatedGame).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("Game not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	//if exists proceed to update
	if err := config.DB.Model(&updatedGame).Where("id = ?", gameID).Update("stock", request.NewStock).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "Success update game with ID" + fmt.Sprintf("%d", updatedGame.ID),
		Note:    "Stock is changed to " + fmt.Sprintf("%d", updatedGame.Stock),
		Data:    updatedGame,
	})
}
