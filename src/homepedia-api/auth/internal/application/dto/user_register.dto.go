package dto

// UserRegisterDTO is a data transfer object for user registration

type UserRegisterDTO struct {
	Username string `json:"username" validate:"required,gte=5,lte=20"`
	Password string `json:"password" validate:"gte=8,lte=20,required"`
	Email    string `json:"email" validate:"email,required"`
}
