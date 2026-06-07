package domain

import "time"

type Estudiante struct {
	ID              uint      `json:"id"`
	Nombre          string    `json:"nombre"`
	ApellidoPat     string    `json:"apellido_pat"`
	ApellidoMat     string    `json:"apellido_mat"`
	FechaNacimiento time.Time `json:"fecha_nacimiento"`
	Matricula       string    `json:"matricula"`
	Carrera         string    `json:"carrera"`
	FechaIncripcion time.Time `json:"fecha_incripcion"`
}

type EstudianteRepository interface {
	Create(e *Estudiante) error
	FindByID(id uint) (*Estudiante, error)
	FindAll() ([]*Estudiante, error)
	Update(e *Estudiante) error
	Delete(id uint) error
	FindByFechaIncripcion(desde, hasta time.Time) ([]*Estudiante, error)
}

type EstudianteService interface {
	CrearEstudiante(nombre, apellidoPat, apellidoMat, matricula, carrera string, fechaNac, fechaInc time.Time) (*Estudiante, error)
	ObtenerEstudiante(id uint) (*Estudiante, error)
	ObtenerTodos() ([]*Estudiante, error)
	ActualizarEstudiante(id uint, nombre, apellidoPat, apellidoMat, matricula, carrera string) (*Estudiante, error)
	EliminarEstudiante(id uint) error
	BuscarPorFechaIncripcion(desde, hasta time.Time) ([]*Estudiante, error)
}
