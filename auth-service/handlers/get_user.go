package handlers

import (
	"auth-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(userService services.UserServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el ID del usuario desde el middleware de autenticación
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
			return
		}

		// Convertir a uint (el middleware puede enviar diferentes tipos)
		var userID uint
		switch v := userIDInterface.(type) {
		case uint:
			userID = v
		case float64:
			userID = uint(v)
		case int:
			userID = uint(v)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error procesando ID de usuario"})
			return
		}

		user, err := userService.GetUserDetails(userID)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if err.Error() == "usuario no encontrado" {
				statusCode = http.StatusNotFound
			}

			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
