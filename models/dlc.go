package models

import "gorm.io/gorm"

type DLC struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(250); not null" json:"name"`
	Description string  `gorm:"type:varchar(250); not null" json:"description"`
	Stock       int     `gorm:"type:int; not null" json:"stock"`
	DLCPrice    float64 `gorm:"type:float;not null;default:0;check:dlc_price >= 0" json:"dlc_price"`
	GameID      uint    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:GameID;references:ID" json:"game_id"`

	//Asc
	Game Game `gorm:"foreignKey:GameID;references:ID" json:"-"`
}
