package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/laschenkov67/docgen"
)

type Party struct {
	Name, INN, KPP, Address string
}

type Item struct {
	Name  string
	Qty   int
	Price float64
	Total float64
}

type Invoice struct {
	Number     string
	Date       time.Time
	Seller     Party
	Buyer      Party
	Items      []Item
	GrandTotal float64
}

func main() {
	eng := docgen.New()

	if err := eng.LoadTemplateFile("invoice.docx", "examples/invoice/templates/invoice.docx"); err != nil {
		log.Fatal(err)
	}
	if err := eng.LoadTemplateFile("invoice.pdf", "examples/invoice/templates/invoice.pdf.tmpl"); err != nil {
		log.Fatal(err)
	}

	inv := Invoice{
		Number: "СФ-000123",
		Date:   time.Now(),
		Seller: Party{Name: "ООО «АгроСофт»", INN: "7707083893", KPP: "770701001", Address: "г. Москва, ул. Тверская, 1"},
		Buyer:  Party{Name: "КФХ Иванов И.И.", INN: "500100732259", Address: "Московская обл., Дмитровский р-н"},
		Items: []Item{
			{"Семена пшеницы, кг", 100, 85.50, 8550.00},
			{"Удобрение NPK, кг", 50, 120.00, 6000.00},
		},
		GrandTotal: 14550.00,
	}

	docx, _ := os.Create("invoice.docx")
	defer docx.Close()
	if err := eng.Render(context.Background(), "invoice.docx", inv, docgen.FormatDOCX, docx); err != nil {
		log.Fatal(err)
	}

	pdf, _ := os.Create("invoice.pdf")
	defer pdf.Close()
	if err := eng.Render(context.Background(), "invoice.pdf", inv, docgen.FormatPDF, pdf); err != nil {
		log.Fatal(err)
	}

	log.Println("ok: invoice.docx + invoice.pdf")
}
