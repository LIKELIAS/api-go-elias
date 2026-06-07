package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// NewPostgresDB crea y retorna la conexión a PostgreSQL.
// Soporta dos modos:
//  1. Variable DATABASE_URL (connection string completa) — usado en Lambda/Neon
//  2. Variables separadas DB_HOST, DB_PORT, etc. — usado en local con Docker
func NewPostgresDB() (*sql.DB, error) {
	var dsn string

	if url := os.Getenv("DATABASE_URL"); url != "" {
		// Modo producción: Neon PostgreSQL connection string
		dsn = url
	} else {
		// Modo local: variables separadas
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir la conexión: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con PostgreSQL: %w", err)
	}

	return database, nil
}
