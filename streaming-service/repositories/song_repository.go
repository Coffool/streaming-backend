package repositories

import (
	"streaming-service/models"

	"gorm.io/gorm"
)

type SongRepository interface {
	GetSongURLByID(id uint) (string, error)
}

type songRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) GetSongURLByID(id uint) (string, error) {
	var song models.Song
	// Solo seleccionamos el campo que nos interesa
	err := r.db.Select("audio_url").First(&song, id).Error
	if err != nil {
		return "", err
	}
	return song.AudioURL, nil
}
