package handlers

import (
	"auth-service/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Refresh(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RefreshInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token requerido"})
			return
		}

		now := time.Now()

		// Buscar refresh token en DB
		var stored models.RefreshToken
		if err := db.Where("token = ?", input.RefreshToken).First(&stored).Error; err != nil {
			// No revelar detalles al usuario
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no se pudo renovar el token"})
			return
		}

		// Verificar expiración
		if now.After(stored.ExpiresAt) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no se pudo renovar el token"})
			return
		}

		// Buscar usuario asociado
		var user models.User
		if err := db.First(&user, stored.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no se pudo renovar el token"})
			return
		}

		// Generar nuevo access token
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "defaultsecret"
		}

		accessTTL, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
		claims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"exp":      now.Add(accessTTL).Unix(),
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		accessTokenString, err := accessToken.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar access token"})
			return
		}

		// Rotar refresh token si está cerca de expirar (<20% de TTL)
		refreshTTL, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
		if stored.ExpiresAt.Sub(now) < refreshTTL/5 {
			stored.Token = uuid.New().String()
			stored.ExpiresAt = now.Add(refreshTTL)
			db.Save(&stored)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "token renovado exitosamente",
			"access_token":  accessTokenString,
			"refresh_token": stored.Token,
			"expires_in":    int(accessTTL.Seconds()),
		})
	}
}
