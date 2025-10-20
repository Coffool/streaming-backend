package main

import (
	"auth-service/config"
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/repositories"
	"auth-service/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Cargar variables de entorno (.env opcional en desarrollo)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se encontró archivo .env, usando variables de entorno del sistema")
	} else {
		log.Println("✅ Archivo .env cargado")
	}

	cfg := config.LoadConfig()

	// Conexión a PostgreSQL
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatalf("❌ Error conectando a PostgreSQL: %v", err)
	}
	log.Println("✅ PostgreSQL conectado")

	// Repositorios
	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// Servicios
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, refreshTokenRepo)

	// Router
	r := gin.Default()

	// Rutas públicas
	r.POST("/register", handlers.Register(userService))
	r.POST("/login", handlers.Login(authService))
	r.POST("/refresh", handlers.Refresh(authService))

	// Rutas protegidas
	auth := r.Group("/", middleware.AuthMiddleware())
	auth.POST("/logout", handlers.Logout(authService))
	auth.PUT("/update", handlers.UpdateUser(userService))
	auth.GET("/user/me", handlers.GetCurrentUser(userService))

	fmt.Printf("🚀 Auth-Service corriendo en http://localhost:%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("❌ Error arrancando servidor:", err)
	}
}
