package domain

import (
	"time"

	"github.com/google/uuid"
)

type Credentials struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	RoleID    uint      `gorm:"default:2;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// NewCredentials is a constructor for Credentials with necessary initialization
func NewCredentials(username, password, email string, role uint) *Credentials {
	return &Credentials{
		Username: username,
		Password: password,
		Email:    email,
		RoleID:   role,
	}
}
