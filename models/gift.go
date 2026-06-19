package models

import "gorm.io/gorm"

type Gift struct {
	gorm.Model
	Name              string `gorm:"not null"`
	Description       string
	ReservedQty       int  `gorm:"default:0"`
	TotalReservations int  `gorm:"not null;default:1"`
	UserID            uint `gorm:"not null"`
	PhotoURL          string
}
