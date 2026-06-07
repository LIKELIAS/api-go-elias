package domain

import "time"

// User es la entidad principal del dominio
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository es el puerto de salida (driven port) hacia la persistencia
type UserRepository interface {
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserService es el puerto de entrada (driving port) hacia la lógica de negocio
type UserService interface {
	Register(name, email, password string) (*User, error)
	Login(email, password string) (string, error)
	GetUser(id uint) (*User, error)
	GetAllUsers() ([]*User, error)
	UpdateUser(id uint, name, email string) (*User, error)
	DeleteUser(id uint) error
}
