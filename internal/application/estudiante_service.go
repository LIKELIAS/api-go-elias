package application

import (
	"errors"
	"time"

	"github.com/elias/api-go-elias/internal/domain"
)

type estudianteService struct {
	repo domain.EstudianteRepository
}

func NewEstudianteService(repo domain.EstudianteRepository) domain.EstudianteService {
	return &estudianteService{repo: repo}
}

func (s *estudianteService) CrearEstudiante(nombre, apellidoPat, apellidoMat, matricula, carrera string, fechaNac, fechaInc time.Time) (*domain.Estudiante, error) {
	e := &domain.Estudiante{
		Nombre:          nombre,
		ApellidoPat:     apellidoPat,
		ApellidoMat:     apellidoMat,
		FechaNacimiento: fechaNac,
		Matricula:       matricula,
		Carrera:         carrera,
		FechaIncripcion: fechaInc,
	}
	if err := s.repo.Create(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *estudianteService) ObtenerEstudiante(id uint) (*domain.Estudiante, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("estudiante no encontrado")
	}
	return e, nil
}

func (s *estudianteService) ObtenerTodos() ([]*domain.Estudiante, error) {
	return s.repo.FindAll()
}

func (s *estudianteService) ActualizarEstudiante(id uint, nombre, apellidoPat, apellidoMat, matricula, carrera string) (*domain.Estudiante, error) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("estudiante no encontrado")
	}
	e.Nombre = nombre
	e.ApellidoPat = apellidoPat
	e.ApellidoMat = apellidoMat
	e.Matricula = matricula
	e.Carrera = carrera
	if err := s.repo.Update(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *estudianteService) EliminarEstudiante(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("estudiante no encontrado")
	}
	return s.repo.Delete(id)
}

func (s *estudianteService) BuscarPorFechaIncripcion(desde, hasta time.Time) ([]*domain.Estudiante, error) {
	return s.repo.FindByFechaIncripcion(desde, hasta)
}
