package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     uint        `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID;references:ID" json:"user_id"`
	TotalPrice float64     `gorm:"type:float;not null;default:0;check:total_price >= 0" json:"total_price"`
	CartItems  []*CartItem `gorm:"foreignKey:CartID" json:"cart_items,omitempty"`

	//Assc
	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
