package database

import (
	"database/sql"
	"log"
	"streaming-service/config"
	"sync"

	_ "github.com/lib/pq" // driver PostgreSQL
)

var (
	db   *sql.DB
	once sync.Once
)

// initDB inicializa la conexión una sola vez (singleton).
func initDB() {
	cfg := config.GetConfig()

	conn, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("❌ Error al abrir conexión a PostgreSQL: %v", err)
	}

	// Opcional: ajustar el pool de conexiones
	conn.SetMaxOpenConns(10) // máximo conexiones abiertas
	conn.SetMaxIdleConns(5)  // conexiones inactivas
	// conn.SetConnMaxLifetime(time.Hour) // tiempo máximo de vida

	db = conn
	log.Println("✅ Pool de conexiones a PostgreSQL creado")
}

// GetDB devuelve una conexión activa desde el pool.
// Inicializa la conexión la primera vez que se llama.
func GetDB() *sql.DB {
	once.Do(initDB)
	return db
}
