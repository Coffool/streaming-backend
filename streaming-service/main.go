package main

import (
	"streaming-service/config"
	"streaming-service/database"
	"streaming-service/handlers"
	"streaming-service/middleware"
	"streaming-service/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	cfg := config.GetConfig()

	// Obtener conexión base (sql.DB del pool)
	sqlDB := database.GetDB()

	// Inicializar GORM una sola vez con el esquema "music_streaming"
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "music_streaming.", // 👈 prefijo de esquema
			SingularTable: false,              // mantiene pluralización (Song -> songs)
		},
	})
	if err != nil {
		panic("❌ Error al inicializar GORM")
	}

	// Crear el servicio de canciones UNA sola vez
	songService := services.NewSongService(gormDB)

	r := gin.Default()

	// Ruta de streaming con middleware JWT y servicio inyectado
	protected := r.Group("/", middleware.AuthMiddleware(cfg.JWTSecret))
	protected.GET("/stream", handlers.StreamSongHandler(songService))

	r.Run(":" + cfg.Port)
}
