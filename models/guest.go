package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

	QRCode            string     `json:"qr_code" gorm:"uniqueIndex"`
	EntradaRegistrada bool       `json:"entrada_registrada" gorm:"default:false"`
	DataEntrada       *time.Time `json:"data_entrada"`
	ValidatedBy       *uint      `json:"usuario_validador"`
}

func (g *Guest) BeforeCreate(tx *gorm.DB) (err error) {
	if g.QRCode == "" {
		g.QRCode = uuid.New().String()
	}
	return
}
