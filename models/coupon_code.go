package models

import (
	"time"

	"gorm.io/gorm"
)

type CouponCode struct {
	gorm.Model
	UserID      uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID;references:ID" json:"user_id"`
	Code        string     `gorm:"type:varchar(250);not null" json:"code"`
	Discount    int        `gorm:"type:int; not null" json:"discount"`
	ExpiredDate *time.Time `gorm:"type:time" json:"expired_date"`
	Usable      bool       `gorm:"type:bool;not null" json:"usable"`

	//Assc
	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
