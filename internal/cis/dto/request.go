package dto

import "mime/multipart"

type NewInput struct {
	Type      int                   `form:"type" validate:"required,oneof=1 2 3"`
	StartDate string                `form:"start_date" validate:"required,min=10,max=10"`
	EndDate   string                `form:"end_date" validate:"required,min=10,max=10"`
	File      *multipart.FileHeader `form:"file"`
}

type EditInput struct {
	CisStatus int `json:"cis_status" validate:"required,oneof=1 2"`
}
