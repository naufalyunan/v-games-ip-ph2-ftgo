package models

import "gorm.io/gorm"

type Rental struct {
	gorm.Model
	Status    string `gorm:"type:varchar(250);not null" json:"status"`
	PaymentID uint   `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:PaymentID;references:ID" json:"payment_id"`

	//Assc
	Payment     *Payment      `gorm:"foreignKey:PaymentID;references:ID" json:"payment,omitempty"`
	RentalItems []*RentalItem `gorm:"foreignKey:RentalID" json:"rental_items,omitempty"`
}
