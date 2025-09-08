package main

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Solo cargar .env en desarrollo
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ No se encontró archivo .env, usando variables de entorno del sistema")
		} else {
			log.Println("✅ Archivo .env cargado")
		}
	}

	// DSN de conexión
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("❌ No se encontró la variable de entorno DB_DSN")
	}

	// Conectar PostgreSQL y deshabilitar sentencias preparadas
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Esto resuelve el error de "prepared statement already exists"
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatalf("❌ Error conectando a PostgreSQL: %v", err)
	}
	log.Println("✅ PostgreSQL conectado")

	// Router
	r := gin.Default()
	r.POST("/register", handlers.Register(db))
	r.POST("/login", handlers.Login(db))
	r.POST("/refresh", handlers.Refresh(db)) // endpoint de refresh token

	// Rutas protegidas
	auth := r.Group("/", middleware.AuthMiddleware())
	auth.POST("/logout", handlers.Logout(db))
	auth.PUT("/update", handlers.UpdateUser(db))

	// Puerto servidor
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Auth-Service corriendo en http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Error arrancando servidor:", err)
	}
}
