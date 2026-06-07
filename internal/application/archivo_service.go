package application

import "github.com/elias/api-go-elias/internal/domain"

type archivoService struct {
	repo domain.ArchivoRepository
}

func NewArchivoService(repo domain.ArchivoRepository) domain.ArchivoService {
	return &archivoService{repo: repo}
}

func (s *archivoService) SubirArchivo(nombre, tipo, ruta string, usuarioID uint) (*domain.Archivo, error) {
	a := &domain.Archivo{
		Nombre:    nombre,
		Tipo:      tipo,
		Ruta:      ruta,
		UsuarioID: usuarioID,
	}
	if err := s.repo.Save(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *archivoService) GetArchivos() ([]domain.Archivo, error) {
	return s.repo.GetAll()
}
