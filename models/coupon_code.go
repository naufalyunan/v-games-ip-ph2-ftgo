package models

import (
	"time"

	"gorm.io/gorm"
)

type CouponCode struct {
	gorm.Model   `swaggerignore:"true"`
	UserID       uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID;references:ID" json:"user_id"`
	Code         string     `gorm:"type:varchar(250);not null" json:"code"`
	Discount     int        `gorm:"type:int; not null" json:"discount"`
	ExpiredDate  *time.Time `gorm:"type:time" json:"expired_date"`
	UsableBySelf bool       `gorm:"type:bool;not null" json:"usable_by_self"`

	//Assc
	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
