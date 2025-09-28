// handlers/register.go
package handlers

import (
	"auth-service/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Username  string `json:"username" binding:"required,min=3,max=20"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Birthdate string `json:"birthdate" binding:"required"` // formato: YYYY-MM-DD
	IsArtist  bool   `json:"is_artist"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Parse birthdate
		birthdate, err := time.Parse("2006-01-02", input.Birthdate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "fecha inválida, formato esperado YYYY-MM-DD"})
			return
		}

		// Validar si username o email ya existen
		var existing models.User
		if err := db.Where("username = ? OR email = ?", input.Username, input.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "usuario o email ya registrados"})
			return
		}

		// Encriptar contraseña
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo encriptar la contraseña"})
			return
		}

		role := "user"
		if input.IsArtist {
			role = "artist"
		}

		user := models.User{
			Name:         "user",
			Role:         role,
			Username:     input.Username,
			Email:        input.Email,
			Password:     string(hashedPassword),
			Birthdate:    birthdate,
			Registerdate: time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "usuario creado exitosamente",
			"user": gin.H{
				"id":        user.ID,
				"username":  user.Username,
				"email":     user.Email,
				"birthdate": user.Birthdate.Format("2006-01-02"),
				"role":      user.Role,
			},
		})
	}
}
