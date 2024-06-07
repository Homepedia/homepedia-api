package dto

type UserLoginDTO struct {
	Password string `json:"password" validate:"gte=8,lte=20,required,containsUppercase,containsSpecialCharacter"`
	Email    string `json:"email" validate:"email,required"`
}
