package handlers

import (
	"auth-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Refresh(authService services.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input services.RefreshRequest

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token requerido"})
			return
		}

		response, err := authService.RefreshToken(input)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
