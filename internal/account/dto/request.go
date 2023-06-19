package dto

type LoginInput struct {
	Username string `json:"username" validate:"required,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateInput struct {
	Username string  `json:"username" validate:"required,max=15"`
	Email    string  `json:"email" validate:"email,required"`
	Role     int     `json:"role" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Age      int     `json:"age" validate:"required"`
	Salary   float64 `json:"salary" validate:"gte=0"`
	Position string  `json:"position" validate:"required"`
	Status   int     `json:"status" validate:"required"`
}

type EditPasswordInput struct {
	Email          string `json:"email" validate:"email,required"`
	Password       string `json:"password" validate:"min=8,required"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqfield=Password"`
}
