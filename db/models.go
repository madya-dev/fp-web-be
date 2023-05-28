package db

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID         uint
	Username   string
	Email      string
	Password   string
	Role       uint
	EmployeeID uint
	Employee   Employee
}

type Employee struct {
	ID             uint
	Name           string
	Age            uint
	Salary         float64
	Position       string
	EmployeeStatus EmployeeStatus
	Projects       []Project `gorm:"many2many:employee_projects"`
}

type EmployeeStatus struct {
	ID         uint
	Name       string
	EmployeeID uint
}

type Cis struct {
	ID         uint
	CisType    CisType
	CisStatus  CisStatus
	CisDetail  CisDetail
	EmployeeID uint
	Employee   Employee
}

type CisType struct {
	ID    uint
	Name  string
	CisID uint
}

type CisStatus struct {
	ID    uint
	Name  string
	CisID uint
}

type CisDetail struct {
	ID        uint
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	File      []byte
	CisID     uint
}

type Project struct {
	ID        uint
	Name      string
	Client    string
	Budget    float64
	StartDate time.Time `gorm:"type:date"`
	EndDate   time.Time `gorm:"type:date"`
	Longtime  int64
}

type SalaryCut struct {
	ID        uint
	Name      string
	SalaryCut float64
}

type SalarySlip struct {
	ID           uint
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
	GenerateDate time.Time `gorm:"type:date default:CURRENT_TIMESTAMP"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.Longtime = int64(p.EndDate.Sub(p.StartDate).Hours() / 24)
	return nil
}

func (s *SalarySlip) BeforeCreate(tx *gorm.DB) error {
	s.TotalA = s.BasicSalary + s.Bonus
	s.TotalB = s.PaidLeave + s.Permission + s.Insurance
	s.Total = s.TotalA - s.TotalB
	return nil
}

func InitialMigrate() {

}