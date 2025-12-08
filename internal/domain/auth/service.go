package auth

import (
	"errors"

	"github.com/arturo/autohost-cloud-api/internal/platform"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// Service encapsula la lógica de negocio de autenticación
type Service struct {
	repo Repository
}

// NewService crea una nueva instancia del servicio de autenticación
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Register registra un nuevo usuario
func (s *Service) Register(email, name, password string) (userID string, err error) {
	// Verificar si el usuario ya existe
	existing, _ := s.repo.FindUserByEmail(email)
	if existing != nil {
		return "", ErrUserAlreadyExists
	}

	// Hashear contraseña
	hash, err := platform.HashPassword(password)
	if err != nil {
		return "", err
	}

	// Crear usuario
	return s.repo.CreateUser(email, name, hash)
}

// Login autentica a un usuario
func (s *Service) Login(email, password string) (*User, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := platform.CheckPassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
