package handlers

import (
	"auth-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(authService services.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input services.LoginRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := authService.Login(input)
		if err != nil {
			statusCode := http.StatusInternalServerError

			// Mapear errores específicos a códigos de estado apropiados
			switch err.Error() {
			case "usuario no encontrado", "contraseña incorrecta":
				statusCode = http.StatusUnauthorized
			case "no se pudo generar access token", "no se pudo guardar refresh token", "error limpiando tokens anteriores":
				statusCode = http.StatusInternalServerError
			}

			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
