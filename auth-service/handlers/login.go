package handlers

import (
	"auth-service/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // ✅ Importación de jwt/v5
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginInput struct {
	Identifier string `json:"identifier" binding:"required"` // username o email
	Password   string `json:"password" binding:"required"`
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.Where("username = ? OR email = ?", input.Identifier, input.Identifier).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no encontrado"})
			return
		}

		// Comparar contraseña
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "contraseña incorrecta"})
			return
		}

		// JWT_SECRET
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "defaultsecret" // ⚠️ solo pruebas
		}

		// Access Token
		accessTTL, _ := time.ParseDuration(os.Getenv("ACCESS_TOKEN_TTL"))
		// ✅ Se usa jwt.MapClaims del paquete v5
		claims := jwt.MapClaims{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"exp":      time.Now().Add(accessTTL).Unix(),
		}

		// ✅ La función para crear el token es la misma
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// ✅ La función para firmar el token es la misma
		accessTokenString, err := accessToken.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar access token"})
			return
		}

		// ---- Refresh Token (UUID, 7 días) ----
		refreshTTL, _ := time.ParseDuration(os.Getenv("REFRESH_TOKEN_TTL"))
		refreshToken := uuid.New().String()
		expiry := time.Now().Add(refreshTTL)

		// Guardar refresh token en DB (uno por usuario)
		db.Where("user_id = ?", user.ID).Delete(&models.RefreshToken{}) // limpiar anteriores

		newRT := models.RefreshToken{
			UserID:    user.ID,
			Token:     refreshToken,
			ExpiresAt: expiry,
		}
		if err := db.Create(&newRT).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo guardar refresh token"})
			return
		}

		// Respuesta con tokens
		c.JSON(http.StatusOK, gin.H{
			"message":       "login exitoso",
			"access_token":  accessTokenString,
			"refresh_token": refreshToken,
			"expires_in":    int(accessTTL.Seconds()),
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		})
	}
}
