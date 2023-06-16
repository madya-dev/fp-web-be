package dto

type LoginInput struct {
	Username string `json:"username" validate:"required,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateInput struct {
	Username string `json:"username" validate:"required,max=15"`
	Email    string `json:"email" validate:"email,required"`
	Role     int    `json:"role" validate:"required"`
}

type EditPasswordInput struct {
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"min=8,required"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqfield=Password"`
}
