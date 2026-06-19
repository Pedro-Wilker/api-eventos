package controllers

import (
	"net/http"

	"github.com/Pedro-Wilker/api-eventos/config"
	"github.com/Pedro-Wilker/api-eventos/models"
	"github.com/gin-gonic/gin"
)

type GiftInput struct {
	Name              string `json:"nome_presente" binding:"required"`
	Description       string `json:"descricao_presente"`
	ReservedQty       int    `json:"quantidade_reservas"`
	TotalReservations int    `json:"quantas_reservas_por_presente"`
	PhotoURL          string `json:"foto_presente"`
}

func CreateGift(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input GiftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	totalRes := input.TotalReservations
	if totalRes <= 0 {
		totalRes = 1
	}

	gift := models.Gift{
		Name:              input.Name,
		Description:       input.Description,
		ReservedQty:       input.ReservedQty,
		TotalReservations: totalRes,
		UserID:            userID.(uint),
		PhotoURL:          input.PhotoURL,
	}

	if err := config.DB.Create(&gift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao salvar presente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Presente adicionado com sucesso!", "data": gift})
}

func ListGifts(c *gin.Context) {
	userID, _ := c.Get("userID")

	var gifts []models.Gift

	if err := config.DB.Where("user_id = ?", userID).Find(&gifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar presentes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts})
}
