package dto

import "mime/multipart"

type NewInput struct {
	Type      int                   `form:"type" validate:"required,oneof=1 2 3"`
	StartDate string                `form:"start_date" validate:"required,min=16,max=16"`
	EndDate   string                `form:"end_date" validate:"required,min=16,max=16"`
	File      *multipart.FileHeader `form:"file"`
}

type EditInput struct {
	CisStatus int `json:"cis_status" validate:"required,oneof=1 2"`
}
