package dto

type GenerateInput struct {
	Username     string  `json:"username" validate:"required"`
	Bonus        float64 `json:"bonus" validate:"gte=0"`
	StartPeriode string  `json:"start_periode" validate:"required,min=10,max=10"`
	EndPeriode   string  `json:"end_periode" validate:"required,min=10,max=10"`
}
