package application

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/elias/api-go-elias/internal/domain"
)

type userService struct {
	repo domain.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(name, email, password string) (*domain.User, error) {
	existing, _ := s.repo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("el correo ya está registrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("credenciales inválidas")
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("error al generar el token")
	}
	return signed, nil
}

func (s *userService) GetUser(id uint) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func (s *userService) GetAllUsers() ([]*domain.User, error) {
	return s.repo.FindAll()
}

func (s *userService) UpdateUser(id uint, name, email string) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	user.Name = name
	user.Email = email
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return errors.New("usuario no encontrado")
	}
	return s.repo.Delete(id)
}
