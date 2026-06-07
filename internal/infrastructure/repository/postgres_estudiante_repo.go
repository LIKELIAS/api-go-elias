package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/elias/api-go-elias/internal/domain"
)

type postgresEstudianteRepo struct {
	db *sql.DB
}

func NewPostgresEstudianteRepository(db *sql.DB) domain.EstudianteRepository {
	return &postgresEstudianteRepo{db: db}
}

func (r *postgresEstudianteRepo) Create(e *domain.Estudiante) error {
	query := `
		INSERT INTO estudiantes (nombre, apellido_pat, apellido_mat, fecha_nacimiento, matricula, carrera, fecha_incripcion)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	return r.db.QueryRow(query, e.Nombre, e.ApellidoPat, e.ApellidoMat, e.FechaNacimiento, e.Matricula, e.Carrera, e.FechaIncripcion).Scan(&e.ID)
}

func (r *postgresEstudianteRepo) FindByID(id uint) (*domain.Estudiante, error) {
	e := &domain.Estudiante{}
	query := `SELECT id, nombre, apellido_pat, apellido_mat, fecha_nacimiento, matricula, carrera, fecha_incripcion FROM estudiantes WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&e.ID, &e.Nombre, &e.ApellidoPat, &e.ApellidoMat, &e.FechaNacimiento, &e.Matricula, &e.Carrera, &e.FechaIncripcion)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("estudiante con id %d no encontrado", id)
		}
		return nil, err
	}
	return e, nil
}

func (r *postgresEstudianteRepo) FindAll() ([]*domain.Estudiante, error) {
	query := `SELECT id, nombre, apellido_pat, apellido_mat, fecha_nacimiento, matricula, carrera, fecha_incripcion FROM estudiantes ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*domain.Estudiante
	for rows.Next() {
		e := &domain.Estudiante{}
		if err := rows.Scan(&e.ID, &e.Nombre, &e.ApellidoPat, &e.ApellidoMat, &e.FechaNacimiento, &e.Matricula, &e.Carrera, &e.FechaIncripcion); err != nil {
			return nil, err
		}
		lista = append(lista, e)
	}
	if lista == nil {
		lista = []*domain.Estudiante{}
	}
	return lista, nil
}

func (r *postgresEstudianteRepo) Update(e *domain.Estudiante) error {
	query := `UPDATE estudiantes SET nombre=$1, apellido_pat=$2, apellido_mat=$3, matricula=$4, carrera=$5 WHERE id=$6`
	_, err := r.db.Exec(query, e.Nombre, e.ApellidoPat, e.ApellidoMat, e.Matricula, e.Carrera, e.ID)
	return err
}

func (r *postgresEstudianteRepo) Delete(id uint) error {
	result, err := r.db.Exec(`DELETE FROM estudiantes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("estudiante con id %d no encontrado", id)
	}
	return nil
}

func (r *postgresEstudianteRepo) FindByFechaIncripcion(desde, hasta time.Time) ([]*domain.Estudiante, error) {
	query := `SELECT id, nombre, apellido_pat, apellido_mat, fecha_nacimiento, matricula, carrera, fecha_incripcion 
	          FROM estudiantes WHERE fecha_incripcion BETWEEN $1 AND $2 ORDER BY fecha_incripcion`
	rows, err := r.db.Query(query, desde, hasta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []*domain.Estudiante
	for rows.Next() {
		e := &domain.Estudiante{}
		if err := rows.Scan(&e.ID, &e.Nombre, &e.ApellidoPat, &e.ApellidoMat, &e.FechaNacimiento, &e.Matricula, &e.Carrera, &e.FechaIncripcion); err != nil {
			return nil, err
		}
		lista = append(lista, e)
	}
	if lista == nil {
		lista = []*domain.Estudiante{}
	}
	return lista, nil
}
