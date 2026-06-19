package models

import "gorm.io/gorm"

type Guest struct {
	gorm.Model
	Name               string   `json:"nome"`
	UserID             uint     `json:"user_id"`
	CompanionQty       int      `json:"quantidade_acompanhante"`
	CompanionNames     []string `json:"nome_acompanhante" gorm:"type:jsonb"`
	CompanionEmails    []string `json:"emails_acompanhantes" gorm:"type:jsonb"`
	CompanionPhones    []string `json:"numeros_acompanhantes" gorm:"type:jsonb"`
	CompanionRelations []string `json:"relacoes_acompanhante" gorm:"type:jsonb"`
	Email              string   `json:"email_convidado"`
	Phone              string   `json:"numero_convidado"`
}
