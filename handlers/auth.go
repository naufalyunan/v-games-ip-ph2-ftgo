package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"v-games-ip-ph2-ftgo/config"
	"v-games-ip-ph2-ftgo/models"
	"v-games-ip-ph2-ftgo/utils"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary Register a new user
// @Description Register a new user with email, password, deposit, jwt_token, input_ref_code, full_name, and role
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param shipment body models.User true "New user"
// @Success 201 {object} models.Response
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /users/register [post]
func Register(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError(err.Error()))
	}

	if user.FullName == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Full name must not be empty"))
	}

	if user.Email == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Email must not be empty"))
	}

	if user.Password == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Password must not be empty"))
	}

	if user.Role == "" {
		user.Role = "user"
	}

	//initial deposit value is 0
	user.Deposit = 0

	if user.Role == "admin" {
		user.InputRefCode = ""
		if err := config.DB.Create(&user).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return utils.NewBadRequestError("Email or username already registered")
			}
			return utils.NewInternalError(err.Error())
		}

		return c.JSON(http.StatusCreated, models.Response{
			Message: "Success create admin with ID " + fmt.Sprintf("%d", user.ID),
			Note:    "Created admin account",
			Data: models.ResponseDataRegister{
				ID:       user.ID,
				FullName: user.FullName,
				Email:    user.Email,
			},
		})
	}

	//if not admin
	var coupon models.CouponCode
	if err := config.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return utils.NewBadRequestError("Email or username already registered")
			}
			return utils.NewInternalError(err.Error())
		}

		//create referral code automatically when user first register
		code, err := utils.GenerateCouponCode(user.ID)
		if err != nil {
			return utils.NewInternalError("Internal server error")
		}
		expDate := time.Now().AddDate(0, 1, 0)
		coupon = models.CouponCode{
			UserID:       user.ID,
			Code:         *code,
			Discount:     0,
			ExpiredDate:  &expDate,
			UsableBySelf: false,
		}

		if err := tx.Create(&coupon).Error; err != nil {
			return utils.NewInternalError("Failed creating coupon code for user")
		}

		// return nil will commit the whole transaction
		return nil
	}); err != nil {
		return utils.HandleError(c, err.(*utils.APIError))
	}

	if user.InputRefCode == "" {
		return c.JSON(http.StatusCreated, models.Response{
			Message: "Success create user with ID " + fmt.Sprintf("%d", user.ID),
			Note:    "Created user doesn't input referral",
			Data: models.ResponseDataRegister{
				ID:           user.ID,
				FullName:     user.FullName,
				Email:        user.Email,
				Deposit:      user.Deposit,
				ReferralCode: coupon.Code,
			},
		})
	}

	if err := config.DB.Transaction(func(tx *gorm.DB) error {
		//if referral code input != "", then find the code in db
		//find the coupon code in db
		var couponReferral models.CouponCode
		if err := tx.Where("code = ?", user.InputRefCode).First(&couponReferral).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// if not found, then commit the transaction, but no coupon is created
				return utils.NewNotFoundError("Referral code not found")
			}
			return utils.NewInternalError("Internal Server Error")
		}
		//if the coupon is found then create new coupon for user to use
		codeToUse, err := utils.GenerateCouponCode(user.ID)
		if err != nil {
			return utils.NewInternalError("Internal server error")
		}
		expDate := time.Now().AddDate(0, 1, 0)
		referralToUse := models.CouponCode{
			UserID:       user.ID,
			Code:         *codeToUse,
			Discount:     10,
			ExpiredDate:  &expDate,
			UsableBySelf: true,
		}
		if err := tx.Create(&referralToUse).Error; err != nil {
			return utils.NewInternalError("Failed creating coupon code for user to use")
		}
		return nil
	}); err != nil {
		if apiErr, ok := err.(*utils.APIError); ok {
			if apiErr.Code == 404 {
				return c.JSON(http.StatusCreated, models.Response{
					Message: "Success create user with ID " + fmt.Sprintf("%d", user.ID),
					Note:    "Inputted referral doesn't match any user, no discount coupon created",
					Data: models.ResponseDataRegister{
						ID:           user.ID,
						FullName:     user.FullName,
						Email:        user.Email,
						Deposit:      user.Deposit,
						ReferralCode: coupon.Code,
					},
				})
			}
			return utils.HandleError(c, err.(*utils.APIError))
		}
	}

	return c.JSON(http.StatusCreated, models.Response{
		Message: "Success create user with ID " + fmt.Sprintf("%d", user.ID),
		Note:    "Inputted referral valid, created discount coupon for user",
		Data: models.ResponseDataRegister{
			ID:           user.ID,
			FullName:     user.FullName,
			Email:        user.Email,
			Deposit:      user.Deposit,
			ReferralCode: coupon.Code,
		},
	})
}

// @Summary Login user
// @Description Authenticate a user with email and password
// @Tags Auth
// @Accept application/json
// @Produce json
// @Param shipment body models.LoginPayload true "Login Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} utils.APIError "Bad Request"
// @Failure 403 {object} utils.APIError "Not Found"
// @Failure 500 {object} utils.APIError "Internal Server Error"
// @Router /users/login [post]
func Login(c echo.Context) error {
	var loginPayload models.LoginPayload
	if err := c.Bind(&loginPayload); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError(err.Error()))
	}

	if loginPayload.Email == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Email must not be empty"))
	}

	if loginPayload.Password == "" {
		return utils.HandleError(c, utils.NewBadRequestError("Password must not be empty"))
	}

	//find user by email
	var user models.User
	if err := config.DB.Where("email = ?", loginPayload.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.HandleError(c, utils.NewNotFoundError("User not found"))
		}
		return utils.HandleError(c, utils.NewInternalError("Internal Server Error"))
	}

	//verify the password if the result exist
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
		return utils.HandleError(c, utils.NewBadRequestError("Invalid Email/Password"))
	}
	key := os.Getenv("KEY")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		})
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal Server Error"))
	}

	//update the last_login_date and jwt_token
	if err := config.DB.Model(&user).Updates(map[string]interface{}{"jwt_token": s}).Error; err != nil {
		return utils.HandleError(c, utils.NewInternalError("Internal server error"))
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": s,
	})
}
