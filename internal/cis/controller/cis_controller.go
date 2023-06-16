package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"hrd-be/internal/cis/dto"
	globalResponse "hrd-be/internal/global/response"
	"hrd-be/model"
	"hrd-be/pkg/database"
	inputValidator "hrd-be/pkg/validator"
	"math"
	"path/filepath"
	"strconv"
	"time"
)

func NewCisHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		var newInput dto.NewInput
		if err := c.Bind(&newInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(newInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		filename := uuid.New().String()
		ext := filepath.Ext(newInput.File.Filename)
		allowedTypes := []string{".jpeg", ".jpg", ".png", ".pdf", ".webp"}
		status := false
		for _, allowedType := range allowedTypes {
			if ext == allowedType {
				status = true
				break
			}
		}
		if !status {
			response.DefaultBadRequest()
			response.Data = map[string]string{"errors": "file input type not allowed"}
			c.AbortWithStatusJSON(response.Code, response)
		}

		db := database.Connection()
		startDate, _ := time.Parse("2006-01-02", newInput.StartDate)
		endDate, _ := time.Parse("2006-01-02", newInput.EndDate)

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := c.SaveUploadedFile(newInput.File, filepath.Join("files", filename+ext)); err != nil {
				return errors.New("failed upload file")
			}

			cisDetail := model.CisDetail{
				StartDate: startDate,
				EndDate:   endDate,
				File:      filename + ext,
			}
			if err := tx.Create(&cisDetail).Error; err != nil {
				return err
			}

			cis := model.Cis{
				CisTypeID:   newInput.Type,
				CisDetailID: cisDetail.ID,
				EmployeeID:  newInput.EmployeeID,
			}
			if err := tx.Create(&cis).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			response.DefaultInternalError()
			response.Data = map[string]string{"errors": err.Error()}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultCreated()
		response.Message = "cis created successfully"
		c.JSON(response.Code, response)
	}
}

func GetAllCisHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response

		currentPage := c.Query("page")
		currentPageInt, _ := strconv.Atoi(currentPage)
		if currentPageInt < 1 {
			currentPageInt = 1
		}
		perPage := 20
		fistData := (currentPageInt * perPage) - perPage

		type Cis struct {
			ID     int    `json:"id"`
			Type   string `json:"type"`
			Status string `json:"status"`
			Name   string `json:"name"`
		}

		var cis []model.Cis
		var totalData int64
		var totalPage int
		db := database.Connection()
		db.Model(&model.Cis{}).Count(&totalData)
		totalPage = int(math.Ceil(float64(totalData) / float64(perPage)))
		db.Preload("CisType").
			Preload("CisStatus").
			Preload("CisDetail").
			Preload("Employee").
			Limit(perPage).Offset(fistData).Find(&cis)

		var cleanCis []Cis
		for _, each := range cis {
			var clean Cis
			clean.ID = each.ID
			clean.Type = each.CisType.Name
			clean.Status = each.CisStatus.Name
			clean.Name = each.Employee.Name

			cleanCis = append(cleanCis, clean)
		}

		response.DefaultOK()
		response.Message = "success get CIS"
		response.Data = map[string]interface{}{
			"cis_list": cleanCis,
			"pagination": map[string]int{
				"current_page": currentPageInt,
				"total_page":   totalPage,
			},
		}
		c.JSON(response.Code, response)
	}
}

func CisDetailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		cisId := c.Param("cis_id")
		protocol := "http://"
		if c.Request.TLS != nil {
			protocol = "https://"
		}

		var cis model.Cis

		db := database.Connection()
		result := db.Preload("CisType").
			Preload("CisStatus").
			Preload("CisDetail").
			Preload("Employee").
			Where("id = ?", cisId).
			Find(&cis)

		var count int64
		if result.Count(&count); count == 0 {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		type Cis struct {
			ID        int    `json:"id"`
			Type      string `json:"type"`
			Status    string `json:"status"`
			Name      string `json:"name"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
			File      string `json:"file"`
		}
		cleanCis := Cis{
			ID:        cis.ID,
			Type:      cis.CisType.Name,
			Status:    cis.CisStatus.Name,
			Name:      cis.Employee.Name,
			StartDate: cis.CisDetail.StartDate.String(),
			EndDate:   cis.CisDetail.EndDate.String(),
			File:      protocol + c.Request.Host + "/files/" + cis.CisDetail.File,
		}

		response.DefaultOK()
		response.Message = "get cis detail success"
		response.Data = map[string]interface{}{
			"cis": cleanCis,
		}
		c.JSON(response.Code, response)
	}
}

func EditCisHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response globalResponse.Response
		cisId := c.Param("cis_id")
		var editInput dto.EditInput
		if err := c.BindJSON(&editInput); err != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		validationErrors := inputValidator.RequestBodyValidator(editInput)
		if validationErrors != nil {
			response.DefaultBadRequest()
			response.Data = map[string][]string{"errors": validationErrors}
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		db := database.Connection()
		var count int64
		db.Where("id = ?", cisId).Find(&model.Cis{}).Count(&count)
		if count != 1 {
			response.DefaultNotFound()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		result := db.Model(&model.Cis{}).Where("id = ?", cisId).Update("cis_status_id", editInput.CisStatus)
		if result.Error != nil {
			response.DefaultInternalError()
			c.AbortWithStatusJSON(response.Code, response)
			return
		}

		response.DefaultOK()
		response.Message = "cis status updated successfully"
		c.JSON(response.Code, response)
	}

}
