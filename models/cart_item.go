package models

import (
	"time"

	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	CartID    uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CartID;references:ID" json:"cart_id"`
	GameID    uint       `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:GameID;references:ID" json:"game_id"`
	StartDate *time.Time `gorm:"type:time; not null" json:"start_date"`
	EndDate   *time.Time `gorm:"type:time; not null" json:"end_date"`
	Quantity  int        `gorm:"type:int; not null" json:"quantity"`

	//Assc
	Cart Cart `gorm:"foreignKey:CartID;references:ID" json:"-"`
	Game Game `gorm:"foreignKey:GameID;references:ID" json:"-"`
}
