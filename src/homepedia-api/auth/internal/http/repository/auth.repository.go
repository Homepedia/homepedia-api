package repository

import (
	"errors"
	"homepedia-api/lib/config"
	"homepedia-api/lib/domain"

	"gorm.io/gorm"
)

const (
	ERROR_DUPLICATE_USERNAME = "username already exists"
	ERROR_DUPLICATE_EMAIL    = "email already exists"
	USER_CREATED             = "user created successfully"
)

var ErrUserNotFound = errors.New("user not found")

type AuthRepositoryError struct {
	Message string
	Success bool
}

type AuthRepository interface {
	// Register est une fonction pour enregistrer un nouvel utilisateur
	Register(credentials *domain.Credentials) AuthRepositoryError
	CheckUsernameOrEmail(username string, email string) AuthRepositoryError
	FindUserByEmail(email string) (*domain.Credentials, error)
	// TODO: Login est une fonction pour connecter un utilisateur
}

// authRepository est la structure qui implémente l'interface AuthRepository
type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		db: config.Connections.Auth,
	}
}

func (r *authRepository) Register(credentials *domain.Credentials) AuthRepositoryError {
	checkUniqueProperties := r.CheckUsernameOrEmail(credentials.Username, credentials.Email)
	if !checkUniqueProperties.Success {
		return AuthRepositoryError{
			Message: checkUniqueProperties.Message,
			Success: false,
		}
	}
	if err := r.db.Create(credentials).Error; err != nil {
		return AuthRepositoryError{
			Message: err.Error(),
			Success: false,
		}
	}
	return AuthRepositoryError{
		Message: USER_CREATED,
		Success: true,
	}
}

func (r *authRepository) CheckUsernameOrEmail(username string, email string) AuthRepositoryError {
	var count int64

	r.db.Model(&domain.Credentials{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return AuthRepositoryError{
			Message: ERROR_DUPLICATE_USERNAME,
			Success: false,
		}
	}

	count = 0

	r.db.Model(&domain.Credentials{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return AuthRepositoryError{
			Message: ERROR_DUPLICATE_EMAIL,
			Success: false,
		}
	}

	return AuthRepositoryError{
		Success: true,
	}
}

func (r *authRepository) FindUserByEmail(email string) (*domain.Credentials, error) {
	var credentials domain.Credentials
	if err := r.db.Where("email = ?", email).First(&credentials).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &credentials, nil
}
