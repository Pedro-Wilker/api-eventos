package controllers

import (
	"net/http"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/models"
	"github.com/gin-gonic/gin"
)

type GuestInput struct {
	Name            string   `json:"nome" binding:"required"`
	CompanionQty    int      `json:"quantidade_acompanhante"`
	CompanionNames  []string `json:"nome_acompanhante"`
	Email           string   `json:"email_convidado"`
	Phone           string   `json:"numero_convidado"`
	CompanionEmails []string `json:"emails_acompanhantes"`
	CompanionPhones []string `json:"numeros_acompanhantes"`
}

func CreateGuest(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	var input GuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	guest := models.Guest{
		Name:            input.Name,
		UserID:          userID.(uint),
		CompanionQty:    input.CompanionQty,
		CompanionNames:  input.CompanionNames,
		Email:           input.Email,
		Phone:           input.Phone,
		CompanionEmails: input.CompanionEmails,
		CompanionPhones: input.CompanionPhones,
	}

	if err := config.DB.Create(&guest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar convidado"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Convidado adicionado com sucesso!", "data": guest})
}

func ListGuests(c *gin.Context) {
	userID, _ := c.Get("userID")

	var guests []models.Guest

	if err := config.DB.Where("user_id = ?", userID).Find(&guests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar convidados"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": guests})
}
