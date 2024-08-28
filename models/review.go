package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID;references:ID" json:"user_id"`
	GameID  uint    `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:GameID;references:ID" json:"game_id"`
	Rating  float64 `gorm:"type:float;not null;default:0;check:rating >= 0" json:"rating"`
	Message string  `gorm:"type:varchar(250); not null" json:"message"`

	//Assc
	User *User `gorm:"foreignKey:UserID;references:ID" json:"User,omitempty"`
	Game *Game `gorm:"foreignKey:GameID;references:ID" json:"Game,omitempty"`
}
