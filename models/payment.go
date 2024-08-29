package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	CartID        uint    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CartID;references:ID" json:"cart_id"`
	PaymentPrice  float64 `gorm:"type:float;not null;default:0;check:payment_price >= 0" json:"payment_price"`
	PaymentStatus string  `gorm:"type:varchar(250)" json:"payment_status"`
	PaymentMethod string  `gorm:"type:varchar(250);not null" json:"payment_method"`
	Provider      string  `gorm:"type:varchar(250)" json:"provider"`
	CouponCode    string  `gorm:"type:varchar(250)" json:"coupon_code,omitempty"`

	//Assc
	Cart Cart `gorm:"foreignKey:CartID;references:ID" json:"cart"`
}
