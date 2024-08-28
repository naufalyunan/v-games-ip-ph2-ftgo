package models

import "gorm.io/gorm"

type RentalItem struct {
	gorm.Model
	RentalID   uint `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:RentalID;references:ID" json:"rental_id"`
	CartItemID uint `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CartItemID;references:ID" json:"cart_item_id"`

	//Assc
	Rental   Rental   `gorm:"foreignKey:RentalID;references:ID" json:"-"`
	CartItem CartItem `gorm:"foreignKey:CartItemID;references:ID" json:"-"`
}
