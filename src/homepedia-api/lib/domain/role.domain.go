package domain

type Role struct {
	ID   int    `gorm:"primary_key;autoIncrement"`
	Name string `gorm:"unique;not null"`
	Credentials []Credentials `gorm:"foreignKey:RoleID"`
}

// NewRole is a constructor for Role with necessary initialization
func NewRole(id int, name string) *Role {
	return &Role{
		ID:   id,
		Name: name,
	}
}
