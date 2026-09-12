package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Guest struct {
	gorm.Model
	Name               string         `json:"nome"`
	UserID             uint           `json:"user_id"`
	CompanionQty       int            `json:"quantidade_acompanhante"`
	CompanionNames     datatypes.JSON `json:"nome_acompanhante"`
	CompanionEmails    datatypes.JSON `json:"emails_acompanhantes"`
	CompanionPhones    datatypes.JSON `json:"numeros_acompanhantes"`
	CompanionRelations datatypes.JSON `json:"relacoes_acompanhante"`
	Email              string         `json:"email_convidado"`
	Phone              string         `json:"numero_convidado"`

	// CompanionQRCodes armazena os qr_codes sinteticos gerados pelo frontend
	// para cada acompanhante (uso unico, mesmo que o qr do titular).
	// CompanionCheckedIn rastreia, na mesma ordem, quais ja consumiram entrada.
	CompanionQRCodes   datatypes.JSON `json:"companion_qr_codes"`
	CompanionCheckedIn datatypes.JSON `json:"companion_checked_in"`

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
