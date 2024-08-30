package models

import "gorm.io/gorm"

type Game struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string  `gorm:"type:varchar(250); not null" json:"name"`
	Description string  `gorm:"type:varchar(250); not null" json:"description"`
	Genre       string  `gorm:"type:varchar(250); not null" json:"genre"`
	RentPrice   float64 `gorm:"type:float;not null;default:0;check:rent_price >= 0" json:"rent_price"`
	Studio      string  `gorm:"type:varchar(250); not null" json:"studio"`
	Stock       int     `gorm:"type:int; not null" json:"stock"`

	//Assc
	Reviews []*Review `gorm:"foreignKey:GameID" json:"reviews,omitempty"`
	DLCs    []*DLC    `gorm:"foreignKey:GameID" json:"DLCs,omitempty"`
}
