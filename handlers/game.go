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

// @Summary Get All Games
// @Description Retrieve a list of all games
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /games [get]
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

// @Summary Create New Game
// @Description create a new Game
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body models.Game true "New Game"
// @Success 201 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /games [POST]
func CreateGame(c echo.Context) error {
	var game models.Game
	if err := c.Bind(&game); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid input"))
	}

	if game.Name == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Name must not be empty"))
	}

	if game.Description == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Description must not be empty"))
	}

	if game.RentPrice == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Rent Price must not be empty"))
	}

	if game.RentPrice < 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Rent Price must be a positive value"))
	}

	if game.Studio == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Studio must not be empty"))
	}

	if game.Stock == 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Stock must not be empty"))
	}

	if game.Stock < 0 {
		return utils.HandleError(c, utils.NewBadRequestError("Stock must be a positive number"))
	}

	//create game
	if err := config.DB.Create(&game).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Success add game with ID " + fmt.Sprintf("%d", game.ID),
		Data:    game,
	})
}

// @Summary Get Game By ID
// @Description Get details of a game by it's ID
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Success 200 {object} models.Response
// @Failure 401 {object} utils.APIError "Unauthorized"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /games/{id} [get]
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

// @Summary Update Game By ID
// @Description Update a game by it's ID
// @Tags games
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Game ID"
// @Param product body models.Game true "Update Game"
// @Success 200 {object} models.Response
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /games/{id} [PUT]
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
