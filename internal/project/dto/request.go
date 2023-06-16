package dto

type NewProjectInput struct {
	Name        string  `json:"name" validate:"required"`
	Client      string  `json:"client" validate:"required"`
	Budget      float64 `json:"budget" validate:"required"`
	StartDate   string  `json:"start_date" validate:"required,max=10,min=10"`
	EndDate     string  `json:"end_date" validate:"required,max=10,min=10"`
	EmployeesID []int   `json:"employees_id" validate:"required"`
}
