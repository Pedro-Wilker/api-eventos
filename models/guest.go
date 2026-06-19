package models

import "gorm.io/gorm"

type Guest struct {
	gorm.Model
	Name            string   `gorm:"not null"`
	UserID          uint     `gorm:"not null"`
	CompanionQty    int      `gorm:"default:0"`
	CompanionNames  []string `gorm:"type:jsonb;serializer:json"`
	Email           string
	Phone           string
	CompanionEmails []string `gorm:"type:jsonb;serializer:json"`
	CompanionPhones []string `gorm:"type:jsonb;serializer:json"`
}
