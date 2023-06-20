package dto

type EditInput struct {
	Username string  `json:"username" validate:"required"`
	Role     int     `json:"role" validate:"gte=0"`
	Name     string  `json:"name" validate:"required"`
	Age      int     `json:"age" validate:"required"`
	Salary   float64 `json:"salary" validate:"gte=0"`
	Position string  `json:"position" validate:"required"`
	Status   int     `json:"status" validate:"required"`
}
