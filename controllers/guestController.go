package controllers

import (
	"net/http"

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
