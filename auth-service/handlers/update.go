// archivo: handlers/update.go
package handlers

import (
	"auth-service/models"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UpdateInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Handler para actualizar datos del usuario autenticado
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// obtener el user_id del JWT
		uid, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			return
		}

		var input UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, uid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuario no encontrado"})
			return
		}

		// Validaciones y actualizaciones
		if input.Name != "" {
			user.Name = input.Name
		}

		if input.Username != "" && input.Username != user.Username {
			var count int64
			db.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "username ya está en uso"})
				return
			}
			user.Username = input.Username
		}

		if input.Email != "" && input.Email != user.Email {
			// validar email con regex
			matched, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`, input.Email)
			if !matched {
				c.JSON(http.StatusBadRequest, gin.H{"error": "formato de email inválido"})
				return
			}
			var count int64
			db.Model(&models.User{}).Where("email = ?", input.Email).Count(&count)
			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email ya está en uso"})
				return
			}
			user.Email = input.Email
		}

		if input.Password != "" {
			if len(input.Password) < 6 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "la contraseña debe tener al menos 6 caracteres"})
				return
			}
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
			user.Password = string(hashedPassword)
		}

		// Guardar cambios
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo actualizar el usuario"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "usuario actualizado correctamente",
			"user": gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"username": user.Username,
				"email":    user.Email,
				"role":     user.Role,
			},
		})
	}
}
