package services

import (
	"fmt"
	"os"
	"path/filepath"
	"streaming-service/repositories"

	"gorm.io/gorm"
)

type SongService interface {
	GetSongURL(id uint) (string, error)
}

type songService struct {
	repo repositories.SongRepository
}

func NewSongService(db *gorm.DB) SongService {
	repo := repositories.NewSongRepository(db)
	return &songService{repo: repo}
}

func (s *songService) GetSongURL(id uint) (string, error) {
	// Ruta relativa guardada en la base de datos
	relativePath, err := s.repo.GetSongURLByID(id)
	if err != nil {
		return "", err
	}

	// Leer la base path desde variable de entorno
	basePath := os.Getenv("CONTENT_BASE_PATH")
	if basePath == "" {
		return "", fmt.Errorf("CONTENT_BASE_PATH no está definido en el entorno")
	}

	// Construir la ruta absoluta segura
	fullPath := filepath.Join(basePath, relativePath)

	return fullPath, nil
}
