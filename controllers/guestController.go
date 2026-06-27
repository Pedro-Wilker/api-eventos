package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/models"
	"github.com/gin-gonic/gin"
)

type GuestInput struct {
	Name                  string   `json:"nome" binding:"required"`
	CompanionQty          int      `json:"quantidade_acompanhante"`
	CompanionNames        []string `json:"nome_acompanhante"`
	Email                 string   `json:"email_convidado"`
	Phone                 string   `json:"numero_convidado"`
	CompanionEmails       []string `json:"emails_acompanhantes"`
	CompanionPhones       []string `json:"numeros_acompanhantes"`
	RelacoesAcompanhantes []string `json:"relacoes_acompanhante"`
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
		CompanionNames:     input.CompanionNames,
		Email:              input.Email,
		Phone:              input.Phone,
		CompanionRelations: input.RelacoesAcompanhantes,
		CompanionEmails:    input.CompanionEmails,
		CompanionPhones:    input.CompanionPhones,
	}

	if err := config.DB.Create(&guest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar convidado"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Convidado adicionado com sucesso!", "data": guest})
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
		CompanionNames:     input.CompanionNames,
		Email:              input.Email,
		Phone:              input.Phone,
		CompanionRelations: input.RelacoesAcompanhantes,
		CompanionEmails:    input.CompanionEmails,
		CompanionPhones:    input.CompanionPhones,
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
	guest.CompanionNames = input.CompanionNames
	guest.Email = input.Email
	guest.Phone = input.Phone
	guest.CompanionRelations = input.RelacoesAcompanhantes
	guest.CompanionEmails = input.CompanionEmails
	guest.CompanionPhones = input.CompanionPhones

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

func findGuestByCode(c *gin.Context, codigo string) (*models.Guest, error) {
	role, _ := c.Get("role")
	userID, _ := c.Get("userID")

	var guest models.Guest
	query := config.DB

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Where("id = ? OR qr_code = ?", codigo, codigo).First(&guest).Error
	if err != nil {
		return nil, err
	}
	return &guest, nil
}

func FindGuestByCodeHandler(c *gin.Context) {
	codigo := c.Param("codigo")

	guest, err := findGuestByCode(c, codigo)
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

	guest, err := findGuestByCode(c, input.Codigo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "invalido",
			"error":  "Convidado não encontrado no sistema.",
		})
		return
	}

	if guest.EntradaRegistrada {
		c.JSON(http.StatusConflict, gin.H{
			"status":   "duplicado",
			"data":     guest,
			"mensagem": "Entrada já registrada em " + guest.DataEntrada.Format("02/01/2006 15:04"),
		})
		return
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)
	now := time.Now()

	guest.EntradaRegistrada = true
	guest.DataEntrada = &now
	guest.ValidatedBy = &uid

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
