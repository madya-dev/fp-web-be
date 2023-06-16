package dto

type GenerateInput struct {
	EmployeeID   int     `json:"employee_id" validate:"required"`
	Bonus        float64 `json:"bonus" validate:"required"`
	StartPeriode string  `json:"start_periode" validate:"required,min=10,max=10"`
	EndPeriode   string  `json:"end_periode" validate:"required,min=10,max=10"`
}
