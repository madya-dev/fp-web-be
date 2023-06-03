package model

import (
	"gorm.io/gorm"
	"hrd-be/pkg/database"
	"log"
	"time"
)

type SalaryCut struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	SalaryCut float64
}

type SalarySlip struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	Name         string
	Position     string
	Status       string
	StartPeriode time.Time `gorm:"type:date"`
	EndPeriode   time.Time `gorm:"type:date"`
	BasicSalary  float64
	Bonus        float64
	TotalA       float64
	PaidLeave    float64
	Permission   float64
	Insurance    float64
	TotalB       float64
	Total        float64
	GenerateDate time.Time `gorm:"type:date"`
}

type Project struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Name      string
	Client    string
	Budget    float64
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	Longtime  int
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
	ID        int       `gorm:"primaryKey;autoIncrement"`
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	File      []byte
}

type Employee struct {
	ID               int `gorm:"primaryKey;autoIncrement"`
	Name             string
	Age              int
	Salary           float64
	Position         string
	EmployeeStatusID int
	EmployeeStatus   EmployeeStatus
	Projects         []Project `gorm:"many2many:employee_projects"`
}

type Account struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	Username   string
	Email      string
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
	CisStatus   CisStatus
	CisDetailID int
	CisDetail   CisDetail
	EmployeeID  int
	Employee    Employee
}

func (s *SalarySlip) BeforeCreate(tx *gorm.DB) error {
	s.TotalA = s.BasicSalary + s.Bonus
	s.TotalB = s.PaidLeave + s.Permission + s.Permission
	s.Total = s.TotalA - s.TotalB
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
