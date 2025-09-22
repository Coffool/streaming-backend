// repositories/refresh_token_repository.go
package repositories

import (
	"auth-service/models"

	"gorm.io/gorm"
)

type RefreshTokenRepositoryInterface interface {
	FindByToken(token string) (*models.RefreshToken, error)
	Create(refreshToken *models.RefreshToken) error // ← NUEVO MÉTODO
	Update(refreshToken *models.RefreshToken) error
	DeleteByUserID(userID uint) error
}

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepositoryInterface {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) Update(refreshToken *models.RefreshToken) error {
	return r.db.Save(refreshToken).Error
}

func (r *RefreshTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}
