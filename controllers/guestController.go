package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// buildCheckedInSlice retorna slice de n bools inicializado em false.
// Usado ao criar guest com acompanhantes para inicializar flags de consumo.
func buildCheckedInSlice(n int) []bool {
	out := make([]bool, n)
	return out
}

type GuestInput struct {
	Name                  string   `json:"nome" binding:"required"`
	CompanionQty          int      `json:"quantidade_acompanhante"`
	CompanionNames        []string `json:"nome_acompanhante"`
	CompanionQRCodes      []string `json:"companion_qr_codes"`
	Email                 string   `json:"email_convidado"`
	Phone                 string   `json:"numero_convidado"`
	CompanionEmails       []string `json:"emails_acompanhantes"`
	CompanionPhones       []string `json:"numeros_acompanhantes"`
	RelacoesAcompanhantes []string `json:"relacoes_acompanhante"`
}

type GuestResponse struct {
	ID                     uint            `json:"ID"`
	Nome                   string          `json:"nome"`
	QrCode                 string          `json:"qr_code"`
	EntradaRegistrada      bool            `json:"entrada_registrada"`
	DataEntrada            *time.Time      `json:"data_entrada"`
	QuantidadeAcompanhante int             `json:"quantidade_acompanhante"`
	NomeAcompanhante       json.RawMessage `json:"nome_acompanhante"`
	RelacoesAcompanhante   json.RawMessage `json:"relacoes_acompanhante"`
	CompanionQRCodes       json.RawMessage `json:"companion_qr_codes"`
	CompanionCheckedIn     json.RawMessage `json:"companion_checked_in"`
}

type ClientGroup struct {
	UserID     uint            `json:"user_id"`
	UserName   string          `json:"user_name"`
	Total      int64           `json:"total"`
	Convidados []GuestResponse `json:"convidados"`
}

func toJSON(v interface{}) datatypes.JSON {
	b, _ := json.Marshal(v)
	return datatypes.JSON(b)
}

func CreateGuest(c *gin.Context) {
	userID, _ := c.Get("userID")
	var input GuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}
	guest := models.Guest{
		Name:               input.Name,
		UserID:             userID.(uint),
		CompanionQty:       input.CompanionQty,
		CompanionNames:     toJSON(input.CompanionNames),
		Email:              input.Email,
		Phone:              input.Phone,
		CompanionRelations: toJSON(input.RelacoesAcompanhantes),
		CompanionEmails:    toJSON(input.CompanionEmails),
		CompanionPhones:    toJSON(input.CompanionPhones),
		CompanionQRCodes:   toJSON(input.CompanionQRCodes),
		CompanionCheckedIn: toJSON(buildCheckedInSlice(len(input.CompanionQRCodes))),
	}
	if err := config.DB.Create(&guest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar convidado"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Convidado adicionado com sucesso!", "data": guest})
}

func ListGuestsByClient(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	if role == "admin" {
		// Struct com tags gorm para o GORM mapear corretamente as colunas
		var results []struct {
			UserID        uint            `gorm:"column:user_id"`
			UserName      string          `gorm:"column:user_name"`
			GuestID       uint            `gorm:"column:guest_id"`
			GuestName     string          `gorm:"column:guest_name"`
			QrCode        string          `gorm:"column:qr_code"`
			Entrada       bool            `gorm:"column:entrada_registrada"`
			DataEnt       *time.Time      `gorm:"column:data_entrada"`
			QtdAcomp      int             `gorm:"column:companion_qty"`
			NomesAcomp    json.RawMessage `gorm:"column:companion_names"`
			RelacoesAcomp json.RawMessage `gorm:"column:companion_relations"`
		}

		err := config.DB.Table("guests").
			Select(`
				users.id        AS user_id,
				users.name      AS user_name,
				guests.id       AS guest_id,
				guests.name     AS guest_name,
				guests.qr_code,
				guests.entrada_registrada,
				guests.data_entrada,
				guests.companion_qty,
				guests.companion_names,
				guests.companion_relations
			`).
			Joins("JOIN users ON users.id = guests.user_id").
			Where("guests.deleted_at IS NULL").
			Scan(&results).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar convidados"})
			return
		}

		groups := make(map[uint]*ClientGroup)
		for _, r := range results {
			if _, ok := groups[r.UserID]; !ok {
				groups[r.UserID] = &ClientGroup{
					UserID:     r.UserID,
					UserName:   r.UserName,
					Convidados: []GuestResponse{},
				}
			}
			groups[r.UserID].Convidados = append(groups[r.UserID].Convidados, GuestResponse{
				ID:                     r.GuestID,
				Nome:                   r.GuestName,
				QrCode:                 r.QrCode,
				EntradaRegistrada:      r.Entrada,
				DataEntrada:            r.DataEnt,
				QuantidadeAcompanhante: r.QtdAcomp,
				NomeAcompanhante:       r.NomesAcomp,
				RelacoesAcompanhante:   r.RelacoesAcomp,
			})
			groups[r.UserID].Total++
		}

		var final []ClientGroup
		totalConvidados := 0
		for _, g := range groups {
			final = append(final, *g)
			totalConvidados += int(g.Total)
		}
		c.JSON(http.StatusOK, gin.H{
			"total_clientes":   len(final),
			"total_convidados": totalConvidados,
			"clientes":         final,
		})
		return
	}

	// View CLIENT
	var user models.User
	config.DB.First(&user, userID)
	var guests []models.Guest
	config.DB.Where("user_id = ?", userID).Find(&guests)

	var guestRes []GuestResponse
	for _, g := range guests {
		guestRes = append(guestRes, GuestResponse{
			ID:                     g.ID,
			Nome:                   g.Name,
			QrCode:                 g.QRCode,
			EntradaRegistrada:      g.EntradaRegistrada,
			DataEntrada:            g.DataEntrada,
			QuantidadeAcompanhante: g.CompanionQty,
			NomeAcompanhante:       json.RawMessage(g.CompanionNames),
			RelacoesAcompanhante:   json.RawMessage(g.CompanionRelations),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"user_id":    user.ID,
		"user_name":  user.Name,
		"total":      len(guestRes),
		"convidados": guestRes,
	})
}

func PublicCreateGuest(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuário inválido"})
		return
	}

	var input GuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	guest := models.Guest{
		Name:               input.Name,
		UserID:             uint(userID),
		CompanionQty:       input.CompanionQty,
		CompanionNames:     toJSON(input.CompanionNames),
		Email:              input.Email,
		Phone:              input.Phone,
		CompanionRelations: toJSON(input.RelacoesAcompanhantes),
		CompanionEmails:    toJSON(input.CompanionEmails),
		CompanionPhones:    toJSON(input.CompanionPhones),
		CompanionQRCodes:   toJSON(input.CompanionQRCodes),
		CompanionCheckedIn: toJSON(buildCheckedInSlice(len(input.CompanionQRCodes))),
	}

	if err := config.DB.Create(&guest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar convidado"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Presença confirmada com sucesso!", "data": guest})
}

func ListGuests(c *gin.Context) {
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var guests []models.Guest
	query := config.DB

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&guests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar convidados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": guests})
}

func UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var guest models.Guest
	query := config.DB

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&guest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Convidado não encontrado ou sem permissão"})
		return
	}

	var input GuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	guest.Name = input.Name
	guest.CompanionQty = input.CompanionQty
	guest.CompanionNames = toJSON(input.CompanionNames)
	guest.Email = input.Email
	guest.Phone = input.Phone
	guest.CompanionRelations = toJSON(input.RelacoesAcompanhantes)
	guest.CompanionEmails = toJSON(input.CompanionEmails)
	guest.CompanionPhones = toJSON(input.CompanionPhones)
	if len(input.CompanionQRCodes) > 0 {
		guest.CompanionQRCodes = toJSON(input.CompanionQRCodes)
		// preserva flags ja consumidas; cria novas entradas como false
		var checked []bool
		_ = json.Unmarshal(guest.CompanionCheckedIn, &checked)
		if len(checked) < len(input.CompanionQRCodes) {
			for i := len(checked); i < len(input.CompanionQRCodes); i++ {
				checked = append(checked, false)
			}
			guest.CompanionCheckedIn = toJSON(checked)
		}
	}

	config.DB.Save(&guest)
	c.JSON(http.StatusOK, gin.H{"message": "Convidado atualizado!", "data": guest})
}

func DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := config.DB
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Where("id = ?", id).Delete(&models.Guest{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar convidado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Convidado deletado com sucesso!"})
}

// findGuestByCode busca convidado por QR code em todo o sistema
// (validacao na porta e operacao ampla do evento).
// Retorna o guest e o indice do acompanhante que fez match, ou -1 se foi
// o titular. Aplica filtro por user_id apenas para nao-admin que estejam
// consultando/manipulando proprios convidados (ListGuests/Update/Delete).
func findGuestByCode(c *gin.Context, codigo string) (*models.Guest, int, error) {
	var guest models.Guest

	// 1) titular (qr_code ou id)
	if err := config.DB.Where("qr_code = ? OR id::text = ?", codigo, codigo).First(&guest).Error; err == nil {
		return &guest, -1, nil
	}

	// 2) acompanhante (codigo sintetico gerado no frontend)
	// LATERAL JOIN com jsonb_array_elements_text expande cada string do
	// array como linha, comparando contra o codigo via igualdade direta.
	// Guarda com jsonb_typeof = 'array' no ON para que a funcao nao seja
	// invocada em colunas NULL/escalar (geraria SQLSTATE 22023 e quebraria
	// o checkin de qualquer outro guest).
	var acompGuest models.Guest
	err := config.DB.Raw(`
		SELECT g.* FROM guests g
		JOIN LATERAL jsonb_array_elements_text(companion_qr_codes) AS elem
		  ON jsonb_typeof(companion_qr_codes) = 'array'
		WHERE elem = ?
		LIMIT 1
	`, codigo).Scan(&acompGuest).Error
	if err == nil && acompGuest.ID != 0 {
		var codes []string
		if acompGuest.CompanionQRCodes != nil {
			if err := json.Unmarshal(acompGuest.CompanionQRCodes, &codes); err == nil {
				for i, code := range codes {
					if code == codigo {
						return &acompGuest, i, nil
					}
				}
			}
		}
		// Encontrou guest mas nao achou indice especifico; retorna sem indice
		// (-1) — CheckinGuest trata isso como titular por seguranca.
		return &acompGuest, -1, nil
	}

	return nil, -1, gorm.ErrRecordNotFound
}

func FindGuestByCodeHandler(c *gin.Context) {
	codigo := c.Param("codigo")

	guest, _, err := findGuestByCode(c, codigo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Convidado não encontrado no sistema."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": guest})
}

type CheckinInput struct {
	Codigo string `json:"codigo" binding:"required"`
}

func CheckinGuest(c *gin.Context) {
	var input CheckinInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Informe o código do QR Code ou ID do convidado."})
		return
	}

	guest, compIdx, err := findGuestByCode(c, input.Codigo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "invalido",
			"error":  "Convidado não encontrado no sistema.",
		})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)
	now := time.Now()

	if compIdx == -1 {
		// Titular
		if guest.EntradaRegistrada {
			ts := ""
			if guest.DataEntrada != nil {
				ts = guest.DataEntrada.Format("02/01/2006 15:04")
			}
			c.JSON(http.StatusConflict, gin.H{
				"status":   "duplicado",
				"data":     guest,
				"mensagem": "Entrada do titular já registrada em " + ts,
			})
			return
		}
		guest.EntradaRegistrada = true
		guest.DataEntrada = &now
		guest.ValidatedBy = &uid
	} else {
		// Acompanhante individual
		var checked []bool
		_ = json.Unmarshal(guest.CompanionCheckedIn, &checked)
		for len(checked) <= compIdx {
			checked = append(checked, false)
		}
		if checked[compIdx] {
			// recupera timestamp aproximado do titular (mesma coluna re-aproveitada)
			ts := ""
			if guest.DataEntrada != nil {
				ts = guest.DataEntrada.Format("02/01/2006 15:04")
			}
			c.JSON(http.StatusConflict, gin.H{
				"status":   "duplicado",
				"data":     guest,
				"mensagem": "Entrada do acompanhante já registrada (titular) em " + ts,
			})
			return
		}
		checked[compIdx] = true
		guest.CompanionCheckedIn = toJSON(checked)
		// nao tocamos em EntradaRegistrada do titular; data_entrada guarda o ultimo checkin
		guest.DataEntrada = &now
		guest.ValidatedBy = &uid
	}

	if err := config.DB.Save(&guest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao registrar entrada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "valido",
		"data":     guest,
		"mensagem": "Entrada autorizada com sucesso!",
	})
}

// buildCheckedInSlice definida no topo do arquivo.
