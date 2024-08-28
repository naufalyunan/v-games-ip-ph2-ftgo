package models

import (
	"v-games-ip-ph2-ftgo/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName     string  `gorm:"type:varchar(250); not null" json:"full_name"`
	Email        string  `gorm:"type:varchar(250);uniqueIndex; not null" json:"email"`
	Password     string  `gorm:"type:varchar(250); not null" json:"password"`
	Deposit      float64 `gorm:"type:float;not null;default:0;check:deposit >= 0" json:"deposit"`
	JWTToken     string  `gorm:"type:varchar(250)" json:"jwt_token"`
	InputRefCode string  `gorm:"type:VARCHAR(250)" json:"input_ref_code"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password == "" {
		return utils.NewBadRequestError("Password must not be empty")
	}
	hashed, err := utils.HashPassword([]byte(u.Password))
	if err != nil {
		return utils.NewInternalError("Internal server error")
	}
	u.Password = string(hashed)
	return nil
}
