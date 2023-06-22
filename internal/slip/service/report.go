package service

import (
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"hrd-be/model"
	"log"
	"path"
	"time"
)

func GenerateSlip(slip model.SalarySlip) (string, error) {
	p := message.NewPrinter(language.Indonesian)

	pdf := gofpdf.New("P", "mm", "A5", "")
	// header
	pdf.SetHeaderFunc(func() {
		pdf.ImageOptions("./files/slips/avatar.png", 45, 6, 60, 0, false, gofpdf.ImageOptions{}, 0, "")
		width, _ := pdf.GetPageSize()
		pdf.Line(0, 25, width, 25)
	})
	// footer
	pdf.SetFooterFunc(func() {
		pdf.SetFont("Arial", "", 8)
		_, height := pdf.GetPageSize()
		local, _ := time.LoadLocation("Asia/Makassar")
		pdf.Text(10, height-10, slip.GenerateDate.In(local).Format("Monday, 02 January 2006 15:04:05"))
	})

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Text(55.8, 34, "SALARY SLIP")

	pdf.SetFont("Arial", "", 12)
	pdf.Text(20, 44, "Name")
	pdf.Text(50, 44, ":")
	pdf.Text(54, 44, slip.Name)

	pdf.Text(20, 50, "Position")
	pdf.Text(50, 50, ":")
	pdf.Text(54, 50, slip.Position)

	pdf.Text(20, 56, "Status")
	pdf.Text(50, 56, ":")
	pdf.Text(54, 56, slip.Status)

	pdf.Text(20, 66, "Start Period")
	pdf.Text(50, 66, ":")
	pdf.Text(54, 66, slip.StartPeriode.Format("Monday, 02 January 2006"))

	pdf.Text(20, 72, "End Period")
	pdf.Text(50, 72, ":")
	pdf.Text(54, 72, slip.EndPeriode.Format("Monday, 02 January 2006"))

	pdf.Text(20, 82, "Basic Salary")
	pdf.Text(50, 82, ":")
	pdf.Text(54, 82, p.Sprintf("Rp%.2f", slip.BasicSalary))

	pdf.Text(20, 88, "Bonus")
	pdf.Text(50, 88, ":")
	pdf.Text(54, 88, p.Sprintf("Rp%.2f", slip.Bonus))

	pdf.Text(20, 94, "Total A")
	pdf.Text(50, 94, ":")
	pdf.Text(54, 94, p.Sprintf("Rp%.2f", slip.TotalA))

	pdf.Text(20, 104, "Paid Leave")
	pdf.Text(50, 104, ":")
	pdf.Text(54, 104, p.Sprintf("Rp%.2f", slip.PaidLeave))

	pdf.Text(20, 110, "Permission")
	pdf.Text(50, 110, ":")
	pdf.Text(54, 110, p.Sprintf("Rp%.2f", slip.Permission))

	pdf.Text(20, 116, "Insurance")
	pdf.Text(50, 116, ":")
	pdf.Text(54, 116, p.Sprintf("Rp%.2f", slip.Insurance))

	pdf.Text(20, 122, "Total B")
	pdf.Text(50, 122, ":")
	pdf.Text(54, 122, p.Sprintf("Rp%.2f", slip.TotalB))

	pdf.SetFont("Arial", "BU", 12)
	pdf.Text(20, 134, "Total")
	pdf.Text(54, 134, p.Sprintf("Rp%.2f", slip.Total))

	pdf.SetFont("Arial", "B", 12)
	pdf.Text(50, 134, ":")

	filename := uuid.New().String() + ".pdf"
	filePath := path.Join("files", "slips", filename)
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Fatalf("ERROR GenerateReport fatal error:, %v", err)
		return "", err
	}
	return filename, nil
}
