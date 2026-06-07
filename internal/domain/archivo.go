package domain

type Archivo struct {
	ID        uint   `json:"id"`
	Nombre    string `json:"nombre"`
	Tipo      string `json:"tipo"`
	Ruta      string `json:"ruta"`
	UsuarioID uint   `json:"usuario_id"`
}

type ArchivoRepository interface {
	Save(a *Archivo) error
	GetAll() ([]Archivo, error)
}

type ArchivoService interface {
	SubirArchivo(nombre, tipo, ruta string, usuarioID uint) (*Archivo, error)
	GetArchivos() ([]Archivo, error)
}
