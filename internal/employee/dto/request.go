package dto

type EditInput struct {
	Name     string  `json:"name" validate:"required"`
	Age      int     `json:"age" validate:"required"`
	Salary   float64 `json:"salary" validate:"required"`
	Position string  `json:"position" validate:"required"`
	Status   int     `json:"status" validate:"required"`
}
