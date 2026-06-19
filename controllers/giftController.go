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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
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
	role, _ := c.Get("role")

	var gifts []models.Gift
	query := config.DB

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&gifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar presentes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts})
}

func UpdateGift(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	var gift models.Gift
	query := config.DB

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&gift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Presente não encontrado ou sem permissão"})
		return
	}

	var input GiftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	gift.Name = input.Name
	gift.Description = input.Description
	gift.TotalReservations = input.TotalReservations
	gift.ReservedQty = input.ReservedQty
	gift.PhotoURL = input.PhotoURL

	config.DB.Save(&gift)
	c.JSON(http.StatusOK, gin.H{"message": "Presente atualizado!", "data": gift})
}

func DeleteGift(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := config.DB
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Where("id = ?", id).Delete(&models.Gift{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao deletar presente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Presente deletado com sucesso!"})
}

func PublicListGifts(c *gin.Context) {
	clienteID := c.Param("user_id")
	var gifts []models.Gift

	if err := config.DB.Where("user_id = ?", clienteID).Find(&gifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar a lista de presentes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gifts})
}

func ReserveGift(c *gin.Context) {
	giftID := c.Param("id")
	var gift models.Gift

	if err := config.DB.First(&gift, giftID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Presente não encontrado"})
		return
	}

	if gift.ReservedQty >= gift.TotalReservations {
		c.JSON(http.StatusConflict, gin.H{"error": "Este presente já foi totalmente reservado!"})
		return
	}

	gift.ReservedQty++
	config.DB.Save(&gift)

	c.JSON(http.StatusOK, gin.H{"message": "Presente reservado com sucesso!", "data": gift})
}
