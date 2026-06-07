package repository

import (
	"database/sql"

	"github.com/elias/api-go-elias/internal/domain"
)

type postgresArchivoRepo struct {
	db *sql.DB
}

func NewPostgresArchivoRepository(db *sql.DB) domain.ArchivoRepository {
	return &postgresArchivoRepo{db: db}
}

func (r *postgresArchivoRepo) Save(a *domain.Archivo) error {
	query := `INSERT INTO archivos (nombre, tipo, ruta, usuario_id) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(query, a.Nombre, a.Tipo, a.Ruta, a.UsuarioID).Scan(&a.ID)
}

func (r *postgresArchivoRepo) GetAll() ([]domain.Archivo, error) {
	rows, err := r.db.Query(`SELECT id, nombre, tipo, ruta, usuario_id FROM archivos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var archivos []domain.Archivo
	for rows.Next() {
		var a domain.Archivo
		if err := rows.Scan(&a.ID, &a.Nombre, &a.Tipo, &a.Ruta, &a.UsuarioID); err != nil {
			return nil, err
		}
		archivos = append(archivos, a)
	}
	return archivos, nil
}
