package model

import (
	"hrd-be/pkg/database"
	"log"
	"time"

	"gorm.io/gorm"
)

type SalaryCut struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	SalaryCut float64
}

type SalarySlip struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `json:"name"`
	Position     string    `json:"position"`
	Status       string    `json:"status"`
	StartPeriode time.Time `gorm:"type:date" json:"start_periode"`
	EndPeriode   time.Time `gorm:"type:date" json:"end_periode"`
	BasicSalary  float64   `json:"basic_salary"`
	Bonus        float64   `json:"bonus"`
	TotalA       float64   `json:"total_a"`
	PaidLeave    float64   `json:"paid_leave"`
	Permission   float64   `json:"permission"`
	Insurance    float64   `json:"insurance"`
	TotalB       float64   `json:"total_b"`
	Total        float64   `json:"total"`
	GenerateDate time.Time `gorm:"type:date" json:"generate_date"`
}

type Project struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Name      string
	Client    string
	Budget    float64
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	Longtime  int
	Employees []Employee `gorm:"many2many:project_employees"`
}

type EmployeeStatus struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type CisType struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type CisStatus struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

type CisDetail struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	StartDate time.Time
	EndDate   time.Time
	File      string
}

type Employee struct {
	ID               int `gorm:"primaryKey;autoIncrement"`
	Name             string
	Age              int
	Salary           float64
	Position         string
	EmployeeStatusID int
	EmployeeStatus   EmployeeStatus
}

type Account struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	Username   string `gorm:"unique"`
	Email      string `gorm:"unique"`
	Password   string
	Role       int
	EmployeeID int
	Employee   Employee
}

type Cis struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	CisTypeID   int
	CisType     CisType
	CisStatusID int
	CisStatus   CisStatus `gorm:"default:null"`
	CisDetailID int
	CisDetail   CisDetail
	EmployeeID  int
	Employee    Employee
}

func (s *SalarySlip) BeforeCreate(tx *gorm.DB) error {
	s.TotalA = s.BasicSalary + s.Bonus
	s.TotalB = s.PaidLeave + s.Permission + s.Permission + s.Insurance
	s.Total = s.TotalA - s.TotalB
	s.GenerateDate = time.Now()
	return nil
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.Longtime = int(p.EndDate.Sub(p.StartDate).Hours() / 24)
	return nil
}

func InitialMigrate() {
	db := database.Connection()

	log.Println("INFO InitialMigrate: auto migrate start")
	err := db.AutoMigrate(&SalaryCut{}, &SalarySlip{}, &Project{}, &EmployeeStatus{},
		&CisType{}, &CisStatus{}, &CisDetail{}, &Employee{}, &Account{}, &Cis{})
	if err != nil {
		log.Fatalf("ERROR InitalMigrate fatal error: %v", err)
	}
	log.Println("INFO InitialMigrate: auto migrate success")
}
