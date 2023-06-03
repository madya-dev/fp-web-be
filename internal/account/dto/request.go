package dto

type LoginInput struct {
	Username string `json:"username" validate:"required,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}
