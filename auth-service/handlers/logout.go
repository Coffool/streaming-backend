package handlers

import (
	"auth-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener user_id del contexto
		userIDAny, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			return
		}

		// JWT numérico viene como float64
		userID := uint(userIDAny.(float64))

		// Eliminar refresh token de DB
		db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

		c.JSON(http.StatusOK, gin.H{
			"message": "logged_out",
		})
	}
}
